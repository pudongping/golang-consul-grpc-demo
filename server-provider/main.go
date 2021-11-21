package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"server-provider/internal/consul"
	"server-provider/proto"
)

const port = ":50051"
const consulAddress = "127.0.0.1:8500"

type helloWorldServer struct {
}

func (h *helloWorldServer) Hi(ctx context.Context, in *proto.HelloRequest) (*proto.HelloResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &proto.HelloResponse{
		Code: 0,
		Msg:  "Success",
		Data: "Hello " + in.Name,
	}, nil
}

func RegisterToConsul() {
	consul.RegisterService(consulAddress, &consul.ConsulService{
		Ip:   "127.0.0.1",
		Port: 50051,
		Tag:  []string{"hello-world"},
		Name: "say-hello-world",
	})
}

type HealthImpl struct {
}

// consul 服务端会自己发送请求，来进行健康检查
func (health *HealthImpl) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	fmt.Println("health checking")
	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
}

func (health *HealthImpl) Watch(req *grpc_health_v1.HealthCheckRequest, w grpc_health_v1.Health_WatchServer) error {
	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Printf("grpc server is starting at 0.0.0.0:%s\n", port)

	s := grpc.NewServer()
	proto.RegisterSayServer(s, &helloWorldServer{})
	grpc_health_v1.RegisterHealthServer(s, &HealthImpl{})

	RegisterToConsul()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
