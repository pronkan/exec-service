package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/pronkan/exec-service/execservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	isServer := flag.Bool("server", false, "Run as server")
	port := flag.String("port", "8080", "Port to listen on (server) or connect to (client)")
	host := flag.String("host", "localhost", "Host to connect to (client)")
	binPath := flag.String("bin", "echo", "Path to the binary to execute (server)")
	execArgsStr := flag.String("exec-args", "", "Execution arguments (comma-separated)")
	tokenSecret := flag.String("token-secret", "", "Secret for token-based authentication (optional)")
	tokenString := flag.String("token", "", "Authentication token (client only)")
	clientTimeout := flag.Duration("timeout", 0, "Client timeout (e.g., 5s, 1m, 2h)")
	flag.Parse()

	args := strings.Split(*execArgsStr, ",")

	if *isServer {
		grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		httpListener, err := net.Listen("tcp", "localhost:8080") // Listen only on localhost
		if err != nil {
			log.Fatalf("failed to listen on localhost:8080: %v", err)
		}

		grpcServer := grpc.NewServer()
		execservice.RegisterExecServiceServer(grpcServer, &execservice.Server{BinPath: *binPath, TokenSecret: *tokenSecret})
		reflection.Register(grpcServer) // Enable reflection for debugging/testing

		// Start the HTTP server for token generation
		expirationTime := time.Hour // Set the token expiration time (1 hour in this example)
		http.HandleFunc("/token", execservice.GenerateTokenHandler(*tokenSecret, expirationTime))
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%s", *port),
			Handler: http.DefaultServeMux,
		}

		log.Printf("gRPC Server listening on port %s", *port)
		log.Printf("HTTP Token endpoint available at http://localhost:%s/token", "8080") // Updated log message

		// Start gRPC server and HTTP server in separate goroutines
		go func() {
			if err := srv.Serve(httpListener); err != nil && err != http.ErrServerClosed {
				log.Fatalf("failed to serve HTTP: %v", err)
			}
		}()

		go func() {
			if err := grpcServer.Serve(grpcListener); err != nil {
				log.Fatalf("failed to serve gRPC: %v", err)
			}
		}()

		// Keep the main goroutine running
		select {}
	} else {
		execservice.RunClient(*host, *port, args, *tokenString, *clientTimeout)
	}
}
