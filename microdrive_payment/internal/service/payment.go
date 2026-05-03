package service

import (
	"context"
	"fmt"
	"math/rand"

	pb "microdrive_payment/gen/go"
)

type PaymentServer struct {
	pb.UnimplementedPaymentServiceServer
}

func (s *PaymentServer) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	fmt.Printf("Получен запрос на оплату: Заказ %s, Сумма: %.2f\n", req.OrderId, req.Amount)

	success := true
	if req.Amount > 10000 {
		success = false
	}

	transactionID := fmt.Sprintf("tr_%d", rand.Intn(1000000))

	return &pb.PaymentResponse{
		TransactionId: transactionID,
		Success:       success,
		Message:       fmt.Sprintf("Платеж %s", map[bool]string{true: "успешен", false: "отклонен"}[success]),
	}, nil
}
