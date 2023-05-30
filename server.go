package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	pb "github.com/quangtm-hblab/react-grpc-go-no-envoy/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	pb.UnimplementedCalculateServer
}

func (s *Server) Sum(ctx context.Context, in *pb.SumRequest) (*pb.SumResponse, error) {
	log.Printf("Receiverd : %v %v", in.GetNum1(), in.GetNum2())
	return &pb.SumResponse{Result: in.GetNum1() + in.GetNum2()}, nil
}

type grpcMultiplexer struct {
	*grpcweb.WrappedGrpcServer
}

// Handler is used to route requests to either grpc or to regular http
func (m *grpcMultiplexer) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.IsGrpcWebRequest(r) {
			m.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Generate a TLS grpc API
	apiserver, err := generateTLSApi("cert/server.crt", "cert/server.key")
	if err != nil {
		log.Fatal(err)
	}
	// Start listening on a TCP Port
	lis, err := net.Listen("tcp", "localhost:50056")
	if err != nil {
		log.Fatal(err)
	}
	// Register the API server as a Calculate Server
	// The register function is a generated piece by protoc.
	pb.RegisterCalculateServer(apiserver, &Server{})

	// Start serving in a goroutine to not block
	go func() {
		log.Fatal(apiserver.Serve(lis))
	}()

	// Wrap the GRPC Server in grpc-web and also host the UI
	grpcWebServer := grpcweb.WrapServer(apiserver)

	multiplex := grpcMultiplexer{
		grpcWebServer,
	}

	//new http server
	r := http.NewServeMux()
	// Load the static webpage with a http fileserver
	webapp := http.FileServer(http.Dir("ui/react-client/build"))

	r.Handle("/", multiplex.Handler((webapp)))

	// Create a HTTP server and bind the router to it, and set wanted address
	srv := &http.Server{
		Handler:      r,
		Addr:         "localhost:50022",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	// Serve the webapp over TLS
	log.Fatal(srv.ListenAndServeTLS("cert/server.crt", "cert/server.key"))
}

func generateTLSApi(pemPath string, keyPath string) (*grpc.Server, error) {
	cred, err := credentials.NewServerTLSFromFile(pemPath, keyPath)
	if err != nil {
		return nil, err
	}
	s := grpc.NewServer(grpc.Creds(cred))
	return s, nil
}
