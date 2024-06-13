package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/adamnasrudin03/go-simpel-grpc/config"
	pb "github.com/adamnasrudin03/go-simpel-grpc/student"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func getDataStudentByEmail(ctx context.Context, client pb.StudentServiceClient, email string) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	student, err := client.GetStudentByEmail(ctx, &pb.Student{
		Email: email,
	})
	if err != nil {
		log.Printf("Failed to get student: (%v)", err)
		return
	}
	log.Printf("Received: %v", student)
}

// main is the entry point of the client application.
//
// It does not take any arguments.
// It does not return anything.
func main() {
	add := fmt.Sprintf("%s:%d", config.DefaultHost, config.DefaultPort)

	// Create a connection to the server.
	conn, err := grpc.Dial(add, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a client for the StudentService.
	client := pb.NewStudentServiceClient(conn)

	// Add metadata to the context.
	md := metadata.Pairs("api-key", config.ApiKey)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Get data student by email
	fmt.Println("To exit press CTRL+C")
	for {
		email := ""
		fmt.Print("Enter email: ")
		fmt.Scanln(&email)

		getDataStudentByEmail(ctx, client, email)
		time.Sleep(1 * time.Second)
	}
}
