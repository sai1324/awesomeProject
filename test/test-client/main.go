package main

import (
	"context"
	"fmt"
	"log"
	pd "test/grpc/test-server/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	//连接到本机 9090 端口 先采用不加密的方式
	conn, err := grpc.Dial("127.0.0.1:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect:%v", err)
	}
	//关闭连接
	defer conn.Close()
	//和服务端建立连接
	client := pd.NewSayHelloClient(conn)
	//执行rpc的调用
	res, err := client.SayHello(context.Background(), &pd.HelloRequest{RequestName: "shiheng"})
	if err != nil {
		log.Fatalf("did not connect1:%v", err)
	}
	fmt.Println(res.GetResponseMsg())
}
