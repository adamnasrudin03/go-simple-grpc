package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/mail"
	"os"
	"sync"

	"github.com/adamnasrudin03/go-simpel-grpc/config"
	pb "github.com/adamnasrudin03/go-simpel-grpc/student"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type dataStudentServer struct {
	pb.UnimplementedStudentServiceServer
	mu       sync.Mutex
	students []*pb.Student
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (s *dataStudentServer) GetStudentByEmail(ctx context.Context, in *pb.Student) (*pb.Student, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Printf("[GetStudentByEmail] incoming request : %v", in.Email)

	if !IsValidEmail(in.Email) {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid email format")
	}

	for _, v := range s.students {
		if v.Email == in.Email {
			return v, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "Data student not found")
}
func (s *dataStudentServer) loadData() {
	data, err := os.ReadFile("data/students.json")
	if err != nil {
		log.Fatalf("failed to read data: %v", err)
	}
	if err := json.Unmarshal(data, &s.students); err != nil {
		log.Fatalf("failed to load data: %v", err)
	}

}

func newServer() *dataStudentServer {
	s := dataStudentServer{}
	s.loadData()
	return &s
}

func isValidApiKey(keys []string) bool {
	if len(keys) < 1 {
		return false
	}
	for _, v := range keys {
		if v == config.ApiKey {
			return true
		}
	}

	return false
}

// serverInterceptor intercepts gRPC calls and logs the request and response.
func serverInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("[grpc-handler] incoming request: %v \n", info.FullMethod)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Printf("[grpc-handler] Invalid request: %v \n", info.FullMethod)
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	if !isValidApiKey(md["api-key"]) {
		log.Printf("[grpc-handler] Invalid request: %v \n", info.FullMethod)
		return nil, status.Errorf(codes.Unauthenticated, "Unauthorized")
	}

	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("[grpc-handler] Invalid request: %v err: %v \n", info.FullMethod, err)
		return nil, err
	}

	log.Printf("[grpc-handler] Success request: %v \n", info.FullMethod)
	return resp, nil
}

func main() {
	add := fmt.Sprintf("%s:%d", config.DefaultHost, config.DefaultPort)
	lis, err := net.Listen("tcp", add)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(serverInterceptor))
	pb.RegisterStudentServiceServer(grpcServer, newServer())

	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
