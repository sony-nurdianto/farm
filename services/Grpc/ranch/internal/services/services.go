package services

import (
	"io"
	"log"

	"github.com/sony-nurdianto/farm/services/Grpc/ranch/internal/pbgen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RanchServiceServer struct {
	pbgen.UnimplementedRanchServiceServer
}

func (s RanchServiceServer) RegisterAnimal(stream pbgen.RanchService_RegisterAnimalServer) error {
	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			log.Println("RegisterAnimal Services Done")
			return nil
		default:
			msg, err := stream.Recv()

			log.Println(msg)

			if err == io.EOF {
				log.Println("register animal services response is done")
				return nil
			}

			if err != nil {
				return status.Error(codes.Internal, err.Error())
			}

			if err := stream.Send(&pbgen.RegisterAnimalResponse{}); err != nil {
				return status.Error(codes.Internal, err.Error())
			}
		}
	}
}
