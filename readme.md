# go-grpc
this simple using grpc with golang

## compile proto
```sh
    protoc --go_out=paths=source_relative:.  --go-grpc_out=paths=source_relative:. student/student.proto
```

## Running server

```sh
    go run server/main.go
```

## Running client
- new tab terminal and run:
```sh
    go run client/main.go
```
- or using postman