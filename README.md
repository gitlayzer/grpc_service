## gRPC学习文档

![image](https://user-images.githubusercontent.com/77761224/218545624-4d8ccfa9-a4fb-4205-bf53-b365557b61fa.png)

### 1：什么是gRPC

```shell
gRPC是Google开源的一个RPC框架和库，支持多语言之间的通信，底层通信采用的是HTTP2协议，gRPC在设计上使用了ProtoBuf这种接口描述语言，这种IDL语言可以定义各种服务，Google还提供了一种工具，protoc来编译这种IDL语言，生成各种各样的语言来操作服务。
```

### 2：gRPC的特点

|                             特点                             |
| :----------------------------------------------------------: |
|      1：定义服务简单，可以很快的搭建出一个RPC调度服务器      |
| 2：gRPC是与语言无关，平台无关。你定义好一个protobuf协议，就可以用protoc生成不同语言的协议框架 |
|   3：使用HTTP2协议，支持双向流，客户端和服务端可以双向通信   |

### 3：RPC与RESTful的区别

|                             区别                             |
| :----------------------------------------------------------: |
| 1：在客户端和服务端通信还有一种基于http的协议，RESTful 架构模式，RESTful是一种基于资源的操作，它的名词是（资源地址），然后添加了一些动作对这些资源进行操作，而RPC是基于函数，它是动词 |
| 2：RPC一般基于TCP协议，当然gRPC是基于HTTP2协议，但是也是比HTTP协议更加有效率和更多特性，RESTful一般都基于HTTP协议 |
| 3：传输方面，自定义的TCP协议或者使用HTTP2协议，报文体积更小，所以传输数据效率更高，RESTful一般都基于HTTP协议，报文体积大 |
| 4：gRPC用的是Protobuf的IDL语言，会编码为二进制协议的数据，而RESTful一般用Json的数据格式，json格式的编码更耗时 |

### 4：gRPC通信流程

![image](https://user-images.githubusercontent.com/77761224/218545684-e8ac5748-723b-4e7f-9b94-97e8efb9ad0c.png)

```shell
1：客户端（gRPC Stub）调用 A 方法，发起RPC调用，对方请求信息使用Protobuf进行对象系列化压缩（IDL）
2：服务器（gRPC Server）接收到请求后，解码请求体，进行业务逻辑处理并返回，对应相应结果使用Protobuf进行对象序列化压缩（IDL）
3：客户端接收到服务器的响应，解码请求体，回调被调用的 A 方法，唤醒正在等待响应（阻塞）的客户端调用并返回响应的结果
```

### 5：gRPC环境安装

```shell
1：安装protoc，protoc是protobuf编译器，可以将源文件xxx.proto编译成开源语言文件，例如：xxx.pb.go
	1：下载地址：https://github.com/protocolbuffers/protobuf/releases/download/v3.20.0/protoc-3.20.0-win64.zip
	2：解压后将其设置为环境变量（也就是把protoc的目录内bin的目录的绝对路径复制添加到Windows的环境变量中）

2：验证
PS E:\code\goland\demo> protoc --version
libprotoc 3.20.0

3：安装protoc-gen-go-grpc和protoc-gen-go，因为protc没有内置go生成器，想实现protoc -> go的转换的话所以需要安装一下
	1：go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	2：go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	
	安装完成之后它会在我们的GOPATH下面，验证
	
	PS E:\code\goland\demo> protoc-gen-go-grpc --version
	protoc-gen-go-grpc 1.2.0
	PS E:\code\goland\demo> protoc-gen-go --version     
	protoc-gen-go.exe v1.28.1
```

### 6：gRPC使用

```shell
创建项目，创建目录pb，proto，然后再proto下创建proto文件message.proto
```

```protobuf
//  定义protoc的版本
syntax = "proto3";
// ../pb定义了生成的Go文件的位置，pb是go文件的包名
option go_package = "../pb;pb";
// 定义响应数据
message MessageResponse {
    string responseSomething = 1;
}
// 定义请求数据
message MessageRequest {
    string requestSomething = 1;
}

// 定义服务
service MessageService {
    // 定义一个rpc方法，请求数据为MessageRequest，响应数据为MessageResponse
    rpc Send(MessageRequest) returns (MessageResponse) {}
}
```

```shell
然后去生成Go文件，使用方法如下。

PS E:\code\goland\demo> cd .\proto\
PS E:\code\goland\demo\proto> protoc --go_out=. message.proto
PS E:\code\goland\demo\proto> protoc --go-grpc_out=. message.proto
# 然后会在上级目录的pb内生成一个message.pb.go文件，内容很多这里就不列出了。
这里基本上生成的两个文件就是把我们定义的proto的定义去生成对应的结构体
# 然后在pb同级目录创建serviceImpl目录，创建MessageSenderServerImpl.go，内容如下
```

```go
package serviceImpl

import (
	"context"
	"demo/pb"
	"fmt"
)

// 这里是做了一个结构体的嵌套
type MessageSenderServerImpl struct {
	*pb.UnimplementedMessageServiceServer
}


// 这里的方法是根据Proto文件内写的Service去做的具体实现
func (MessageSenderServerImpl) Send(ctx context.Context, request *pb.MessageRequest) (*pb.MessageResponse, error) {
	fmt.Println("Received Message: ", request.GetRequestSomething())
	resp := &pb.MessageResponse{}
	resp.ResponseSomething = "Roger That!"
	return resp, nil
}
```

```shell
随后我们就可以去启动Server服务了，在目录下创建server目录，然后创建main.go的文件内容如下
```

```go
package main

import (
	"demo/pb"
	"demo/serviceImpl"
	"google.golang.org/grpc"
	"net"
)

func main() {
	// new一个grpc服务
	grpcServer := grpc.NewServer()

	// 注册服务
	pb.RegisterMessageServiceServer(grpcServer, &serviceImpl.MessageSenderServerImpl{})
	listener, err := net.Listen("tcp", ":8002") // 监听端口
	if err != nil {
		panic("TCP Listen Error: " + err.Error())
	}

	// 启动RPC服务
	err = grpcServer.Serve(listener)
	if err != nil {
		panic("gRPC Server Error: " + err.Error())
	}
}
```

```shell
这里我们的服务端就完事儿了，后面我们再来实现一个客户端，在同目录创建client目录，然后创建main.go文件，内容如下
```

```go
package main

import (
	"context"
	"demo/pb"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 建立GRPC连接
	conn, err := grpc.Dial("localhost:8002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("gRPC Connection Error: " + err.Error())
	}
	// 关闭连接
	defer conn.Close()

	// 创建客户端
	client := pb.NewMessageServiceClient(conn)
	resp, err := client.Send(context.Background(), &pb.MessageRequest{
		RequestSomething: "Hello GRPC!",
	})
	if err != nil {
		panic("gRPC Client Error: " + err.Error())
	}
	fmt.Println("Response: ", resp.GetResponseSomething())
}

```

```shell
# 然后我们就可以测试服务了，首先启动服务端

PS E:\code\goland\demo\server> go run .\main.go
# 这里阻塞住是正常的

# 然后去操作客户端。
PS E:\code\goland\demo\client> go run .\main.go
Response:  Roger That!

# 我们可以看到，Client收到了我们上面服务的方法实现的返回数据

PS E:\code\goland\demo\server> go run .\main.go
Received Message:  Hello GRPC!

# 服务端也收到了Client发送过来的数据

# 这里我们需要记住一个问题，gRPC客户端必须拿到服务端的pb文件才可以，否则是不可调用的，所以一般在客户端的代码里会去引入那个pb内的Go文件

# 总结：
# 1：GRPC用于微服务内部调用
# 2：网关从http转到grpc的功能
```


