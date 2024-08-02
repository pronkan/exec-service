package execservice

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Server represents the gRPC server for the exec service.
type Server struct {
	UnimplementedExecServiceServer
	BinPath     string
	TokenSecret string
}

// Execute handles the Execute RPC call, executing the command and streaming output.
func (s *Server) Execute(req *ExecuteRequest, stream ExecService_ExecuteServer) error {
	// Authentication (Optional)
	if s.TokenSecret != "" {
		md, ok := metadata.FromIncomingContext(stream.Context())
		if !ok {
			return status.Errorf(codes.Unauthenticated, "missing metadata")
		}
		values := md["authorization"]
		if len(values) == 0 {
			return status.Errorf(codes.Unauthenticated, "missing authorization token")
		}
		tokenString := values[0]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.TokenSecret), nil
		})
		if err != nil || !token.Valid {
			return status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}
	}

	log.Println("Executing command:", s.BinPath, strings.Join(req.Args, " "))

	// Execute command in a new goroutine to allow parallel executions
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		cmd := exec.Command(s.BinPath, req.Args...)
		stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()

		cmd.Start()
		log.Println("Command started with PID:", cmd.Process.Pid)

		// Stream output and wait for the command to finish.
		waitCh := make(chan error, 1)
		go func() { waitCh <- cmd.Wait() }() // Wait for command in a goroutine
		streamOutput(stdout, stream, "stdout")
		streamOutput(stderr, stream, "stderr")
		err := <-waitCh

		if err != nil {
			log.Println("Command finished with error:", err)
			stream.Send(&ExecuteResponse{ExitCode: -1})
		} else {
			log.Println("Command finished successfully")
			stream.Send(&ExecuteResponse{ExitCode: int32(cmd.ProcessState.ExitCode())})
		}
	}()
	wg.Wait()

	return nil
}

// streamOutput streams output from the command to the client.
func streamOutput(pipe io.ReadCloser, stream ExecService_ExecuteServer, pipeName string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		log.Println(pipeName, ":", scanner.Text())
		stream.Send(&ExecuteResponse{Output: scanner.Text()})
	}
}

// GenerateTokenHandler creates an HTTP handler function that generates a JWT token.
func GenerateTokenHandler(secretKey string, expirationTime time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create claims
		claims := jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
		}

		// Create and sign the token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(secretKey))
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, tokenString)
	}
}
