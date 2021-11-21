package consul

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/consul/api"
)

type ConsulService struct {
	Ip   string
	Port int
	Tag  []string
	Name string
}

func RegisterService(consulAddress string, svc *ConsulService) {
	// 注册 consul
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulAddress
	client, err := api.NewClient(consulConfig)
	if err != nil {
		log.Printf("consul NewClient error %v", err)
		return
	}

	agent := client.Agent()
	interval := time.Duration(10) * time.Second
	deregister := time.Duration(1) * time.Second

	reg := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%v-%v-%v", svc.Name, svc.Ip, svc.Port), // 服务节点的名称，比如：say-hello-world-127.0.0.1-50051
		Name:    svc.Name,                                            // 服务名称，比如：say-hello-world
		Tags:    svc.Tag,                                             // tag，可以为空，比如：hello-world
		Port:    svc.Port,                                            // 服务端口，比如：50051
		Address: svc.Ip,                                              // 服务 IP，比如：127.0.0.1
		Check: &api.AgentServiceCheck{ // 健康检查
			Interval:                       interval.String(),                                   // 健康检查间隔
			GRPC:                           fmt.Sprintf("%v:%v/%v", svc.Ip, svc.Port, svc.Name), // grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
			DeregisterCriticalServiceAfter: deregister.String(),                                 // 注销时间，check 失败后 1s 删除本服务，相当于过期时间
		},
	}

	fmt.Printf("registing to %v\n", consulAddress)
	if err := agent.ServiceRegister(reg); err != nil {
		fmt.Printf("Consul Service Register error \n%v", err)
		return
	}

}
