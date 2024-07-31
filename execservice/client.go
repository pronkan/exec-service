package execservice

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
)

func RunClient(host, port string, args []string, tokenString string) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := NewExecServiceClient(conn)

	ctx := context.Background()
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
