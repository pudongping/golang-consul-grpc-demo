# golang-consul-grpc-demo

基于 golang 用 consul 做 grpc 的服务注册与发现示例代码。

## 版本依赖

go 1.16

## 运行

- 下载源码

```shell
git clone https://github.com/pudongping/golang-consul-grpc-demo.git
```

- 启动 consul

```shell
# 我这里使用的环境是 Mac，如果你的环境和我不一样，则需要下载和你系统相符的 consul
# 我使用的 consul 版本是 1.3.0
# 这里原本是准备将 consul 提交上去的，但是 consul 太大了，无法提交，因此就需要你自己下载了
# 下载方式很简单，直接执行以下命令，然后解压即可
wget wget https://releases.hashicorp.com/consul/1.3.0/consul_1.3.0_darwin_amd64.zip && unzip consul_1.3.0_darwin_amd64.zip
./consul agent -dev
```

- 启动服务端

```shell
cd server-provider
go mod tidy
go run main.go
```

- 启动客户端

```shell
cd server-consumer
go mod tidy
go run main.go
```

## 验证

服务端控制台会输出类似于如下内容

```shell
grpc server is starting at 0.0.0.0::50051
registing to 127.0.0.1:8500
health checking
2021/11/21 18:51:56 Received: Alex
2021/11/21 18:51:58 Received: Alex
2021/11/21 18:52:00 Received: Alex
2021/11/21 18:52:02 Received: Alex
health checking
2021/11/21 18:52:04 Received: Alex
2021/11/21 18:52:06 Received: Alex
2021/11/21 18:52:08 Received: Alex
2021/11/21 18:52:10 Received: Alex
2021/11/21 18:52:12 Received: Alex
health checking
2021/11/21 18:52:14 Received: Alex
2021/11/21 18:52:16 Received: Alex
2021/11/21 18:52:18 Received: Alex
```

客户端控制台会输出类似于如下内容

```shell
calling consul init 
calling consul build 
target: {Scheme:consul Authority:127.0.0.1:8500 Endpoint:say-hello-world/hello-world URL:{Scheme:consul Opaque: User: Host:127.0.0.1:8500 Path:/say-hello-world/hello-world RawPath: ForceQuery:false RawQuery: Fragment: RawFragment:}}
target uri: 127.0.0.1:8500/say-hello-world/hello-world
calling consul watcher 
adding service addrs 
newAddrs: [{Addr:127.0.0.1:50051 ServerName: Attributes:<nil> BalancerAttributes:<nil> Type:0 Metadata:<nil>}]
2021/11/21 18:51:56 Success: Code ==> 0, Msg ==> Success, Data ==> Hello Alex
2021/11/21 18:51:58 Success: Code ==> 0, Msg ==> Success, Data ==> Hello Alex
2021/11/21 18:52:00 Success: Code ==> 0, Msg ==> Success, Data ==> Hello Alex
2021/11/21 18:52:02 Success: Code ==> 0, Msg ==> Success, Data ==> Hello Alex
2021/11/21 18:52:04 Success: Code ==> 0, Msg ==> Success, Data ==> Hello Alex
2021/11/21 18:52:06 Success: Code ==> 0, Msg ==> Success, Data ==> Hello Alex
2021/11/21 18:52:08 Success: Code ==> 0, Msg ==> Success, Data ==> Hello Alex
2021/11/21 18:52:10 Success: Code ==> 0, Msg ==> Success, Data ==> Hello Alex
2021/11/21 18:52:12 Success: Code ==> 0, Msg ==> Success, Data ==> Hello Alex
2021/11/21 18:52:14 Success: Code ==> 0, Msg ==> Success, Data ==> Hello Alex
2021/11/21 18:52:16 Success: Code ==> 0, Msg ==> Success, Data ==> Hello Alex
2021/11/21 18:52:18 Success: Code ==> 0, Msg ==> Success, Data ==> Hello Alex
```

表示已经注册成功了。希望对你有所帮助。
