package main

import (
	"fmt"
	"log"
	"net"

	"github.com/stebin13/product-srv/pkg/config"
	"github.com/stebin13/product-srv/pkg/db"
	"github.com/stebin13/product-srv/pkg/pb"
	"github.com/stebin13/product-srv/pkg/services"
	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("failed at config", err)
	}
	h := db.InitDb(&c)
	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Product Svc on", c.Port)

	srv := services.Server{
		H: h,
	}
	grpcserver := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcserver, &srv)

	if err := grpcserver.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
