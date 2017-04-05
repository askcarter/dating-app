package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	pb "github.com/askcarter/dating-game/pb"
	"github.com/askcarter/io16/app/handlers"

	"github.com/braintree/manners"
)

const (
	address     = "localhost:8080"
	defaultName = "world"
)

var (
	port = flag.String("port", ":8081", "which port to client on")
	chat = flag.String("chat", ":18081", "which port to chat on")
)

type route struct {
	Name, Method, Pattern string
	Handler               http.HandlerFunc
}

const (
	apiPrefix = "/api/v1"
)

var routes = []route{
	route{
		Name:    "CardIndex",
		Method:  "GET",
		Pattern: apiPrefix + "/cards",
		Handler: nil,
	},
	route{
		Name:    "Review",
		Method:  "POST",
		Pattern: apiPrefix + "/review/{value}",
		// TODO: restrict path
		// Pattern: apiPrefix + "/review/{value:^(accept|forgot)$}",
		Handler: nil,
	},
	route{
		Name:    "Save",
		Method:  "POST",
		Pattern: apiPrefix + "/save",
		Handler: nil,
	},
}

func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, r := range routes {
		h := logger(r.Handler, r.Name)
		router.Methods(r.Method).Path(r.Pattern).Name(r.Name).Handler(h)
	}
	return router
}

func logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		log.Printf("%s\t%s\t%s\t%s", r.Method, r.RequestURI, name, time.Since(start))
	})
}

// chatServer is used to implement helloworld.chatServer.
type chatServer struct{}

func (s *chatServer) Chat(ctx context.Context, in *pb.ChatRequest) (*pb.ChatResponse, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		log.Fatal("no metadata")
	}
	from := md["from"][0]
	fmt.Printf("%v: %v\n", from, in.Message)
	return &pb.ChatResponse{}, nil
}

func main() {
	errChan := make(chan error, 10)

	// Setup Frontend Server
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./html")))

	httpServer := manners.NewServer()
	httpServer.Addr = *port
	httpServer.Handler = handlers.LoggingHandler(mux)

	go func() {
		errChan <- httpServer.ListenAndServe()
	}()
	fmt.Println("open browser to http://localhost" + *port)

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	fmt.Println("Connected to server server on " + address)

	defer func() {
		fmt.Println("Closing connection to server server on " + address)
		conn.Close()
	}()
	// c := pb.NewGreeterClient(conn)

	// Set up chat server.
	lis, err := net.Listen("tcp", *chat)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	/* TODO: should I gracefully shut down server with s.Stop() or s.GracefulStop()? */

	pb.RegisterClientChatServer(s, &chatServer{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	go func() {
		errChan <- s.Serve(lis)
	}()
	fmt.Println("Chat server started on http://localhost" + *chat)

	// // Contact the server and print out its response.
	// name := defaultName
	// if len(os.Args) > 1 {
	// 	name = os.Args[1]
	// }
	// req, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	// if err != nil {
	// 	log.Fatalf("could not greet: %v", err)
	// }
	// log.Printf("Greeting: %s", req.Message)

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
			httpServer.BlockingClose()
			os.Exit(0)
		}
	}
}
