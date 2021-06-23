package main

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/cluster/cluster_impl"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/loadbalance"
	_ "dubbo.apache.org/dubbo-go/v3/common/proxy/proxy_factory"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/filter/filter_impl"
	_ "dubbo.apache.org/dubbo-go/v3/protocol/dubbo3"
	_ "dubbo.apache.org/dubbo-go/v3/registry/protocol"
	_ "dubbo.apache.org/dubbo-go/v3/registry/zookeeper"
)

import (
	"dubbo3-demo/protobuf"
)

var greeterProvider = new(protobuf.GreeterClientImpl)

func setConfigByAPI() {
	consumerConfig := config.NewConsumerConfig(
		config.WithConsumerAppConfig(config.NewDefaultApplicationConfig()),
		config.WithConsumerRegistryConfig("registryKey", config.NewDefaultRegistryConfig("zookeeper")),
		config.WithConsumerReferenceConfig("greeterImpl", config.NewReferenceConfigByAPI(
			config.WithReferenceRegistry("registryKey"),
			config.WithReferenceProtocol("tri"),
			config.WithReferenceInterface("org.apache.dubbo.UserProvider"),
		)),
	)
	config.SetConsumerConfig(*consumerConfig)
}

func init() {
	config.SetConsumerService(greeterProvider)
	setConfigByAPI()
}

func main() {
	config.Load()
	//sayHello()
	sayHelloStream()
	//select {}
}

func sayHello() {
	ctx := context.Background()
	req := protobuf.HelloRequest{
		Name: "dubbogo-v3",
	}
	user := protobuf.User{}
	err := greeterProvider.SayHello(ctx, &req, &user)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Receive user = %+v\n", user)
}

func sayHelloStream() {
	ctx := context.Background()
	stream, err := greeterProvider.SayHelloStream(ctx)
	if err != nil {
		panic(err)
	}
	name := "dubbogo-v3-"
	i := 0

	req := protobuf.HelloRequest{
		Name: name,
	}
	for {
		i++
		req.Name = name + strconv.Itoa(i)
		err = stream.Send(&req)
		if err != nil {
			panic(err)
		}
		fmt.Printf("client:Send %+v\n", req)
		user, err := stream.Recv()
		if err != nil {
			panic(err)
		}
		fmt.Printf("client:Recv %+v\n", user)
		time.Sleep(time.Second)
	}
}
