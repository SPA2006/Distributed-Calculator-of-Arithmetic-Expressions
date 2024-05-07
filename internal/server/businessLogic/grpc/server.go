package grpc_server

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "github.com/SPA2006/Distributed-Calculator-of-Arithmetic-Expressions/proto/agent/agent"
)

func Compute(ctx context.Context, exp string, addr string) (string, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("could not connect to grpc server (agent):", err)
		return "0", status.Errorf(status.Code(err), "could not connect to grpc server (agent): %v", err)
	}
	defer conn.Close()

	state := conn.GetState().String()
	log.Println("connection state:", state)
	grpcClient := pb.NewAgentClient(conn)

	answer, err := grpcClient.CalculateExpression(ctx, &pb.ExpressionToCalculateRequest{})

	if err != nil {
		log.Println("internal gRPC error, couldn't calculate")
		return "0", status.Error(status.Code(err), " Internal gRPC error, couldn't calculate")
	}

	log.Println("Successful gRPC connection. Operation is done")
	return answer.Result, nil
}
