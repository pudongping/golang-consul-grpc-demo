package main

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"

	"server-consumer/internal/consul"
	"server-consumer/proto"
)

const (
	target      = "consul://127.0.0.1:8500/say-hello-world/hello-world"
	defaultName = "Alex"
)

func main() {
	consul.Init()

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, target, grpc.WithBlock(), grpc.WithInsecure(), grpc.WithBalancerName("round_robin"))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := proto.NewSayClient(conn)

	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	for {
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		r, err := c.Hi(ctx, &proto.HelloRequest{
			Name: name,
		})
		if err != nil {
			log.Fatalf("could not say Hi: %v", err)
		}

		log.Printf("Success: Code ==> %d, Msg ==> %s, Data ==> %s", r.GetCode(), r.GetMsg(), r.GetData())
		time.Sleep(time.Second * 2)

	}

}
