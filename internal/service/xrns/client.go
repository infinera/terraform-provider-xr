package xrns

import (
	"context"
	"log"

	pb "terraform-provider-xrcm/internal/service/xrns/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type XrnsClient struct {
	Endpoint string
}

func (c *XrnsClient) GetDeviceByName(deviceName string) (*pb.Device, error) {
	log.Printf("xrnsClient => %+v\n", c)
	conn, err := grpc.Dial(c.Endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v\n", err)
	}
	defer conn.Close()
	log.Printf("Connection establsihed - %+v\n", conn)

	client := pb.NewNamingServiceClient(conn)
	res, err := client.GetDeviceByName(context.Background(), &pb.GetDeviceByNameRequest{
		Name: deviceName,
	})
	if err != nil {
		log.Printf("Could not get device (%v) by name: %v\n", deviceName, err)
		return nil, err
	}

	return res, nil
}
