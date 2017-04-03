package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/askcarter/dating-game/pb"
)

// Contestant is a player of the game.
type Contestant struct {
	Username string
	Email    string
	Bio      string
}

// server is used to implement helloworld.GreeterServer.
type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

// datingGameServer is used to implement helloworld.datingGameServer.
type datingGameServer struct{}

func (s *datingGameServer) DebugListUsers(ctx context.Context, in *pb.DebugListUsersRequest) (*pb.DebugListUsersReply, error) {
	// Create a User for every planet (and pluto?).  Each planet matches with the planet adjacent to it.
	// var users = []*pb.User{
	// 	{ID: 1, DisplayName: "Mercury", Matches: []int64{2}},
	// 	{ID: 2, DisplayName: "Venus", Matches: []int64{1, 3}},
	// 	{ID: 3, DisplayName: "Earth", Matches: []int64{2, 4}},
	// 	{ID: 4, DisplayName: "Mars", Matches: []int64{3, 5}},
	// 	{ID: 5, DisplayName: "Jupiter", Matches: []int64{4, 6}},
	// 	{ID: 6, DisplayName: "Saturn", Matches: []int64{5, 7}},
	// 	{ID: 7, DisplayName: "Uranus", Matches: []int64{6, 8}},
	// 	{ID: 8, DisplayName: "Neptune", Matches: []int64{7, 9}},
	// 	{ID: 9, DisplayName: "Pluto", Matches: []int64{8}},
	// }

	// ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)

	client, err := spanner.NewClient(ctx, "projects/askcarter-talks/instances/test-instance/databases/test-dating-game")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	iter := client.Single().Read(ctx, "users", spanner.AllKeys(),
		[]string{"id", "username", "matches"})
	defer iter.Stop()

	var users = []*pb.User{}
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var user pb.User
		var matches []spanner.NullInt64
		if err := row.Columns(&user.ID, &user.DisplayName, &matches); err != nil {
			return nil, err
		}

		for _, v := range matches {
			user.Matches = append(user.Matches, v.Int64)
		}

		users = append(users, &user)
	}

	return &pb.DebugListUsersReply{Users: users}, nil
}

func main() {
	errChan := make(chan error, 10)

	port := flag.String("port", ":8080", "Selects action to take.")

	flag.Parse()

	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Creates a new gRPC server
	s := grpc.NewServer()
	pb.RegisterDatingGameServer(s, &datingGameServer{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	go func() {
		errChan <- s.Serve(lis)
	}()

	mockDatabase()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}
		case s := <-signalChan:
			log.Println(fmt.Sprintf("Captured %v. Exiting...", s))
			// httpServer.BlockingClose()

			os.Exit(0)
		}
	}
}

// UserID represents the user.
type UserID int64

// User is a player in the game.
type User struct {
	ID          UserID
	DisplayName string
	Matches     []int64
}

func mockDatabase() error {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)

	client, err := spanner.NewClient(ctx, "projects/askcarter-talks/instances/test-instance/databases/test-dating-game")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	userColumns := []string{"id", "username", "matches"}
	// Create a User for every planet (and pluto?).  Each planet matches with the planet adjacent to it.
	m := []*spanner.Mutation{
		spanner.InsertOrUpdate("users", userColumns, []interface{}{1, "Mercury", []int64{2}}),
		spanner.InsertOrUpdate("users", userColumns, []interface{}{2, "Venus", []int64{1, 3}}),
		spanner.InsertOrUpdate("users", userColumns, []interface{}{3, "Earth", []int64{2, 4}}),
		spanner.InsertOrUpdate("users", userColumns, []interface{}{4, "Mars", []int64{3, 5}}),
		spanner.InsertOrUpdate("users", userColumns, []interface{}{5, "Jupiter", []int64{4, 6}}),
		spanner.InsertOrUpdate("users", userColumns, []interface{}{6, "Saturn", []int64{5, 7}}),
		spanner.InsertOrUpdate("users", userColumns, []interface{}{7, "Uranus", []int64{6, 8}}),
		spanner.InsertOrUpdate("users", userColumns, []interface{}{8, "Neptune", []int64{7, 9}}),
		spanner.InsertOrUpdate("users", userColumns, []interface{}{9, "Pluto", []int64{8}}),
	}

	_, err = client.Apply(ctx, m)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
