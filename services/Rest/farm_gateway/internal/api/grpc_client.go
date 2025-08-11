package api

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClientConn = grpc.ClientConn

type GrpcConn interface {
	Close() error
	Connection() *GrpcClientConn
}

type grpcClient struct {
	conn *grpc.ClientConn
}

func NewGrpcClient(target string) (GrpcConn, error) {
	insCred := insecure.NewCredentials()
	trsCred := grpc.WithTransportCredentials(insCred)

	conn, err := grpc.NewClient(target, trsCred)
	if err != nil {
		return nil, err
	}

	return grpcClient{conn}, nil
}

func (c grpcClient) Close() error {
	return c.conn.Close()
}

func (c grpcClient) Connection() *GrpcClientConn {
	return c.conn
}
