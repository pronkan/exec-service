## JWT Auth

Can be used to generate and JWT tokens from secrets.  

### Usage

```bash
go run main.go -key=secret
```

### Example

```bash
go run main.go -key=$(openssl rand -base64 32)
```
