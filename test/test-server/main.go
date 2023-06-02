package main

import (
	"context"
	"net"
	pd "test/grpc/test-server/proto"

	"google.golang.org/grpc"
) //调用自己的包的格式不是很熟

type server struct {
	pd.UnimplementedSayHelloServer
}

func (s *server) SayHello(ctx context.Context, req *pd.HelloRequest) (*pd.HelloResponse, error) {
	return &pd.HelloResponse{ResponseMsg: "hello" + req.RequestName}, nil
}
func main() {
	//开启端口进行监听
	listen, _ := net.Listen("tcp", ":9090")
	//创建grpc服务
	grpcServer := grpc.NewServer()

	//注册服务
	pd.RegisterSayHelloServer(grpcServer, &server{})
	//启动服务
	grpcServer.Serve(listen)
}
