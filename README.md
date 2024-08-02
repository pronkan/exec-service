# exec-service

This is a simple service that executes a command on the server and returns the output to the client.

## Usage

### Server
```bash
go run main.go -server -bin=echo -port=8080 -token-secret=$(openssl rand -base64 32)
```
Server started:
```bash
2024/08/02 03:14:51 gRPC Server listening on port 8080
2024/08/02 03:14:51 HTTP Token endpoint available at http://localhost:8080/token
```

### Client
```bash
export JWT=$(curl http://localhost:8080/token) # get jwt token for 1 hour
go run main.go -host=localhost -port=8080 -token=${JWT} -command="Hello, World!"
```