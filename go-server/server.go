package main

import (
	"context"
	"fmt"
	"time"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/common/proxy/proxy_factory"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/filter/filter_impl"
	_ "dubbo.apache.org/dubbo-go/v3/protocol/dubbo3"
	_ "dubbo.apache.org/dubbo-go/v3/registry/protocol"
	_ "dubbo.apache.org/dubbo-go/v3/registry/zookeeper"
	tripleConstant "github.com/dubbogo/triple/pkg/common/constant"
)

import (
	"dubbo3-demo/protobuf"
)

func setConfigByAPI() {
	providerConfig := config.NewProviderConfig(
		config.WithProviderAppConfig(config.NewDefaultApplicationConfig()),
		config.WithProviderProtocol("tripleProtocolKey", "tri", "20000"),
		config.WithProviderRegistry("registryKey", config.NewDefaultRegistryConfig("zookeeper")),

		config.WithProviderServices("greeterImpl", config.NewServiceConfigByAPI(
			config.WithServiceRegistry("registryKey"),
			config.WithServiceProtocol("tripleProtocolKey"),
			config.WithServiceInterface("org.apache.dubbo.UserProvider"),
		)),
	)
	config.SetProviderConfig(*providerConfig)
}

func init() {
	setConfigByAPI()
}

func main() {
	config.SetProviderService(NewGreeterProvider())
	config.Load()
	select {}
}

type GreeterProvider struct {
	*protobuf.GreeterProviderBase
}

func NewGreeterProvider() *GreeterProvider {
	return &GreeterProvider{
		GreeterProviderBase: &protobuf.GreeterProviderBase{},
	}
}

func (s *GreeterProvider) SayHelloStream(svr protobuf.Greeter_SayHelloStreamServer) error {
	for {
		c, err := svr.Recv()
		if err != nil {
			return err
		}
		fmt.Printf("server recv name = %s\n", c.Name)

		svr.Send(&protobuf.User{
			Name: "hello " + c.Name,
			Age:  19,
			Id:   "123456789",
		})
		fmt.Printf("server Send user: %+v\n", c.Name)
		time.Sleep(time.Second)
	}
}

func (s *GreeterProvider) SayHello(ctx context.Context, in *protobuf.HelloRequest) (*protobuf.User, error) {
	fmt.Printf("Dubbo3 GreeterProvider get user name = %s\n", in.Name)
	fmt.Println("get triple header tri-req-id = ", ctx.Value(tripleConstant.TripleCtxKey(tripleConstant.TripleRequestID)))
	fmt.Println("get triple header tri-service-version = ", ctx.Value(tripleConstant.TripleCtxKey(tripleConstant.TripleServiceVersion)))
	return &protobuf.User{Name: "Hello " + in.Name, Id: "12345", Age: 21}, nil
}

func (g *GreeterProvider) Reference() string {
	return "greeterImpl"
}
