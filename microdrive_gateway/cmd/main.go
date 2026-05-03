package main

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"microdrive_gateway/handler"
	pb "microdrive_gateway/pkg/proto"
)

func main() {

	authConn, err := grpc.NewClient("auth-service:44044", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	authClient := pb.NewAuthClient(authConn)

	storageConn, err := grpc.NewClient("storage-service:50055", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	storageClient := pb.NewImageServiceClient(storageConn)

	paymentConn, err := grpc.NewClient("payment-service:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	paymentClient := pb.NewPaymentServiceClient(paymentConn)

	h := handler.NewHandler(authClient, storageClient, paymentClient)

	srv := h.InitRoutes()

	log.Println("Gateway running on :8080")
	srv.Run(":8080")
}
