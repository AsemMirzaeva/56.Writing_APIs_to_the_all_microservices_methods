package main

import (
	"context"
	"io"
	"log"
	"net"
	"time"

	productpb "grpc-interceptors/proto/product"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type server struct {
	productpb.UnimplementedProductServiceServer
}

func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	h, err := handler(ctx, req)

	end := time.Now()

	log.Printf("Method: %s, Duration: %s, Error: %v", info.FullMethod, end.Sub(start), err)
	if err != nil {
		st, _ := status.FromError(err)
		log.Printf("Error code: %s, Error message: %s", st.Code(), st.Message())
	}

	return h, err
}

func (s *server) GetProduct(ctx context.Context, req *productpb.ProductRequest) (*productpb.ProductResponse, error) {
	res := &productpb.ProductResponse{
		Id:    req.GetId(),
		Name:  "Example Product",
		Price: 99.99,
	}
	return res, nil
}

func (s *server) BidirectionalStreaming(stream productpb.ProductService_BidirectionalStreamingServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("Received message: %s", req.Id)
		
		if err := stream.Send(&productpb.ProductResponse{
			Id:    req.Id,
			Name:  "Received",
			Price: 0.0,
		}); err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)
	productpb.RegisterProductServiceServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
