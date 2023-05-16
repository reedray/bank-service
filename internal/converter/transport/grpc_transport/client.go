package grpc_transport

import (
	"context"
	"github.com/reedray/bank-service/api/pb/converter/gen_converter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type Client struct {
	Conn   *grpc.ClientConn
	Client gen_converter.ConvertServiceClient
}

func NewClient(addr string) *Client {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Create a client instance
	client := gen_converter.NewConvertServiceClient(conn)
	return &Client{
		Conn:   conn,
		Client: client,
	}
}

func (c *Client) Convert(ctx context.Context, in *gen_converter.Money, opts ...grpc.CallOption) (*gen_converter.Money, error) {
	converted, err := c.Client.Convert(ctx, in)
	if err != nil {
		return nil, err
	}
	return converted, nil
}
