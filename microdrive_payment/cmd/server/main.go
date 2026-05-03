package main

import (
	"log"
	"net"

	pb "microdrive_payment/gen/go"
	"microdrive_payment/internal/service"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPaymentServiceServer(s, &service.PaymentServer{})

	log.Println("Payment Service gRPC сервер запущен на :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
