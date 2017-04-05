package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

type clientChatConn struct {
	client pb.ClientChatClient
	// I need to store this to close the connection or else the server
	// will keep trying to connect to it.
	conn *grpc.ClientConn
}

// datingGameServer is used to implement helloworld.datingGameServer.
type datingGameServer struct {
	// This is temporary until the IP list is moved to the
	// the database.  Otherwise this wouldn't be able to scale.
	//
	ConnMap map[string]clientChatConn
}

func (s *datingGameServer) ListUsers(ctx context.Context, in *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
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

	// ctx, _ := context.WithTimeout(ctx, 1*time.Minute)

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

	return &pb.ListUsersResponse{Users: users}, nil
}

// SendChat takes a message from a client and sends it to another client.  The receiver is determined
// by the metadata associated with the request.
// TODO: Return different errors codes, based on results.
func (s *datingGameServer) SendChat(ctx context.Context, in *pb.ChatRequest) (*pb.ChatResponse, error) {
	// // streaming api
	// headers, ok := metadata.FromContext(stream.Context())
	// token := headers["authorization"]
	// for {
	//     request := new(Request)
	//     err := stream.RecvMsg(request)
	//     // do work
	//     err := stream.SendMsg(response)
	// }

	// Retrieve metadata sent with this message.
	md, ok := metadata.FromContext(ctx)
	if !ok {
		log.Fatal("no metadata")
	}
	to := md["to"][0]
	from := md["from"][0]
	fmt.Printf(" (%v -> %v): %#v\n", from, to, in.Message)

	// Store message in spanner for later retrieval by clients.
	// TODO: Add a read field.  I think I can store the message
	// after calling the gRPC chat function.  If the request errors
	// out, then I'll say that message is unread.  Or, maybe I'll
	// I have the frontend send a read receipt?  Might be simpler.
	err := s.storeChatInDB(to, from, in)
	if err != nil {
		log.Fatal("error storing chat in DB:", err)
	}

	// Send message to client chat server.

	// TODO: if the client isn't connected, don't send.
	c, found := s.ConnMap[to]
	if !found {
		// TODO: return error message.
		return nil, fmt.Errorf("No connection found in map for '%v'", to)
	}
	header := metadata.New(map[string]string{"from": from})
	// this is the critical step that includes your headers
	newCtx := metadata.NewContext(ctx, header)
	return c.client.Chat(newCtx, &pb.ChatRequest{Message: in.Message})
}

// ChatHistory returns the history between two clients. The pb.HistoryRequest metadata tells the server which
// key to request. The key is in the format of '<lower id UserName>#<higher id UserName>'. For example, if there were
// two users Mercury and Venus with IDs 1 and 2, respectively then the metadata sent over would be "key:Mercury:Venus".
func (s *datingGameServer) ChatHistory(ctx context.Context, in *pb.HistoryRequest) (*pb.HistoryResponse, error) {
	client, err := spanner.NewClient(ctx, "projects/askcarter-talks/instances/test-instance/databases/test-dating-game")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	md, ok := metadata.FromContext(ctx)
	if !ok {
		log.Fatal("no metadata")
	}
	mark := md["mark"][0]

	stmt := spanner.NewStatement("SELECT timestamp_created, sender, message FROM test_chat WHERE mark = @key")
	stmt.Params["key"] = mark
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()

	var chats = []*pb.ChatRequest{}
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var chat pb.ChatRequest
		if err := row.Columns(&chat.TimestampCreated, &chat.Sender, &chat.Message); err != nil {
			return nil, err
		}

		chats = append(chats, &chat)
	}

	return &pb.HistoryResponse{ChatHistory: chats}, nil
}

// This operates as a hacky GUUID for chats. This doesn't scale, as there no mutex around it.
var numChats = 0

func getSizeOfTable(ctx context.Context, db string) (int, error) {
	client, err := spanner.NewClient(ctx, "projects/askcarter-talks/instances/test-instance/databases/test-dating-game")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	iter := client.Single().Read(ctx, db, spanner.AllKeys(),
		[]string{"timestamp_created", "sender", "message"})
	defer iter.Stop()

	numRows := 1
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		numRows++
	}
	return numRows, nil
}

// Store message in spanner for later retrieval by clients.  Messages between two clients are stored by concatenating
// the name of the two clients, as such: The key is in the format of '<lower id UserName>#<higher id UserName>'.
// For example, if there were two users Mercury and Venus with IDs 1 and 2, respectively then the metadata sent
// over would be "key:Mercury:Venus".
// TODO: Add a 'Read' field.
func (s *datingGameServer) storeChatInDB(to, from string, msg *pb.ChatRequest) error {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)

	client, err := spanner.NewClient(ctx, "projects/askcarter-talks/instances/test-instance/databases/test-dating-game")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Create a User for every planet (and pluto?).  Each planet matches with the planet adjacent to it.
	ids := map[string]int{
		"Mercury": 1,
		"Venus":   2,
		"Earth":   3,
		"Mars":    4,
		"Jupiter": 5,
		"Saturn":  6,
		"Uranus":  7,
		"Neptune": 8,
		"Pluto":   9,
	}

	toID := ids[to]
	fromID := ids[from]

	mark := ""
	if toID < fromID {
		mark = to + "#" + from
	} else {
		mark = from + "#" + to
	}
	chatColumns := []string{"id", "mark", "sender", "message", "timestamp_created"}

	if numChats == 0 {
		numChats, _ = getSizeOfTable(ctx, "test_chat")
	}
	m := []*spanner.Mutation{
		spanner.InsertOrUpdate("test_chat", chatColumns, []interface{}{numChats, mark, from, msg.Message, time.Now().Unix()}),
	}

	_, err = client.Apply(ctx, m)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// Matches returns the matches for a specific client.  This function uses metadata
// to determine whose mathces to send back.
// TODO: Consider just passing the requester in the pb.MatchesRequest struct.
func (s *datingGameServer) Matches(ctx context.Context, in *pb.MatchesRequest) (*pb.MatchesResponse, error) {
	client, err := spanner.NewClient(ctx, "projects/askcarter-talks/instances/test-instance/databases/test-dating-game")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	md, ok := metadata.FromContext(ctx)
	if !ok {
		log.Fatal("no metadata")
	}
	id := md["id"][0]
	// id, _ := strconv.Atoi(md["id"][0])

	row, err := client.Single().ReadRow(ctx, "users",
		spanner.Key{id}, []string{"username", "matches"})
	if err != nil {
		return nil, err
	}

	var user pb.User
	var matches []spanner.NullInt64
	if err = row.Columns(&user.DisplayName, &matches); err != nil {
		return nil, err
	}
	for _, v := range matches {
		user.Matches = append(user.Matches, v.Int64)
	}

	if len(matches) == 0 {
		return &pb.MatchesResponse{}, nil
	}

	str := "SELECT id, username, matches FROM users WHERE id IN ("
	str += strings.Trim(strings.Join(strings.Fields(fmt.Sprint(user.Matches)), ","), "[]")
	str += ")"

	stmt := spanner.NewStatement(str)
	// stmt.Params["keys"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(user.Matches)), ","), "[]")
	stmt.Params["keys"] = user.Matches
	iter := client.Single().Query(ctx, stmt)
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

	return &pb.MatchesResponse{Users: users}, nil
}

// Connect registers a client with the server.
// TODO: rename to Register.
func (s *datingGameServer) Connect(ctx context.Context, in *pb.ConnectRequest) (*pb.ConnectResponse, error) {
	// Setup connection to chat server
	conn, err := grpc.Dial(in.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	s.ConnMap[in.UserName] = clientChatConn{client: pb.NewClientChatClient(conn), conn: conn}
	fmt.Println("Registed " + in.UserName + ".")

	return &pb.ConnectResponse{}, nil
}

// Disconnect unregisters a client with the server.
// TODO: rename to Unregister.
func (s *datingGameServer) Disconnect(ctx context.Context, in *pb.DisconnectRequest) (*pb.DisconnectResponse, error) {
	s.ConnMap[in.UserName].conn.Close()
	delete(s.ConnMap, in.UserName)
	fmt.Println("Unregisted " + in.UserName + ".")
	return &pb.DisconnectResponse{}, nil
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
	pb.RegisterDatingGameServer(s, &datingGameServer{
		ConnMap: make(map[string]clientChatConn),
	})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	/* TODO: should I gracefully shut down server with s.Stop() or s.GracefulStop()? */

	go func() {
		errChan <- s.Serve(lis)
	}()
	fmt.Println("Server started on http://localhost" + *port)

	// TODO: remove once users can signup.
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
