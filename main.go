package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/pronkan/exec-service/execservice"
	"google.golang.org/grpc"
)

func main() {
	isServer := flag.Bool("server", false, "Run as server")
	port := flag.String("port", "8080", "Port to listen on (server) or connect to (client)")
	host := flag.String("host", "localhost", "Host to connect to (client)")
	binPath := flag.String("bin", "Rscript", "Path to the binary to execute (server)")
	execArgsStr := flag.String("exec-args", "", "Execution arguments (comma-separated)")
	tokenSecret := flag.String("token-secret", "", "Secret for token-based authentication (optional)")
	tokenString := flag.String("token", "", "Authentication token (client only)")
	flag.Parse()

	args := strings.Split(*execArgsStr, ",")

	if *isServer {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		execservice.RegisterExecServiceServer(grpcServer, &execservice.Server{BinPath: *binPath, TokenSecret: *tokenSecret})
		log.Printf("Server listening on port %s", *port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	} else {
		execservice.RunClient(*host, *port, args, *tokenString)
	}
}
