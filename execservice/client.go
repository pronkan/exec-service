package execservice

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func RunClient(host, port string, args []string, tokenString string, timeout time.Duration) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := NewExecServiceClient(conn)

	// Create context with timeout (if provided)
	ctx := context.Background()
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	// Add token to metadata
	if tokenString != "" {
		md := metadata.Pairs("authorization", tokenString)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	stream, err := client.Execute(ctx, &ExecuteRequest{Args: args})
	if err != nil {
		log.Fatalf("could not execute: %v", err)
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break // End of stream
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		if response.ExitCode != 0 {
			log.Fatalf("command failed with exit code %d", response.ExitCode)
		}
		fmt.Println(response.Output)
	}
}
