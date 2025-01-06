package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	pb "product-service/gen/proto"
	"product-service/internal/products"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type productService struct {
	pb.UnimplementedProductServiceServer
	c *products.Conf
}

func (u productService) GetProductInfo(ctx context.Context, req *pb.ProductRequest) (*pb.ProductResponse, error) {
	productIds := req.GetProductIds()
	if productIds == nil {
		return nil, status.Error(codes.InvalidArgument, "ProductIds is nil")
	}

	var productInfos []products.ProductInfo
	productInfos, err := u.c.GetProductInfos(ctx, productIds)
	if err != nil {
		slog.Error("Error in fetching product Infos", slog.Any("productIds", productIds))
		return nil, status.Error(codes.InvalidArgument, "Invalid productIds")
	}

	fmt.Println(productInfos)
	var resp []*pb.Product
	for _, item := range productInfos {
		productResp := pb.Product{
			ProductId: item.Id,
			Stock:     int32(item.Stock),
			PriceId:   item.PriceId,
			Price:     int32(item.Price),
		}
		resp = append(resp, &productResp)
	}
	return &pb.ProductResponse{Products: resp}, nil
}

func RegistergRPCMethod(p *products.Conf) {
	listener, err := net.Listen("tcp", ":5001")

	if err != nil {
		panic(err)
	}

	//NewServer creates a gRPC server which has no service registered
	// creating an instance of the server
	s := grpc.NewServer()

	pb.RegisterProductServiceServer(s, &productService{c: p})

	//exposing gRPC service to be tested by postman
	reflection.Register(s)

	err = s.Serve(listener) // run the gRPC server
	if err != nil {
		panic(err)
	}
}
