package main

import (
	"log"
	"net"
	"os"

	"grpcChatServer/chatserver"

	"google.golang.org/grpc"
)

func main() {

	//assign port
	Port := os.Getenv("PORT")
	if Port == "" {
		Port = "5000" // default Port ser to 5000 if the PORT is not set
	}

	//init listener
	listener, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("Could not listen @ %v :: %v", Port, err)
	}
	log.Println("Listening @ : " + Port)

	//gRPC server instance
	grpcserver := grpc.NewServer()

	//register ChatService
	cs := chatserver.ChatServer{}
	chatserver.RegisterServicesServer(grpcserver, &cs)

	//grpc listen and serve
	err = grpcserver.Serve(listener)
	if err != nil {
		log.Fatalf("Failed to start the gRPC Server :: %v", err)
	}
}
