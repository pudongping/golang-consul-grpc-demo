package consul

import (
	"errors"
	"fmt"
	"regexp"
	"sync"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

const defaultPort = "8500"

var (
	errMissingAddr = errors.New("consul resolver: missing address")

	errAddrMisMatch = errors.New("consul resolver: invalid uri")

	regexConsul, _ = regexp.Compile("^([A-z0-9.]+)(:[0-9]{1,5})?/([A-z_-]+)/([A-z_-]+)$")
)

func Init() {
	fmt.Printf("calling consul init \n")
	resolver.Register(NewBuilder())
}

type consulBuilder struct {
}

type consulResolver struct {
	address              string  // consul 的服务地址，比如：127.0.0.1:8500
	wg                   sync.WaitGroup
	cc                   resolver.ClientConn
	name                 string
	tag                  string
	disableServiceConfig bool
	lastIndex            uint64
}

func NewBuilder() resolver.Builder {
	return &consulBuilder{}
}

func parseTarget(target string) (host, port, name, tag string, err error) {

	fmt.Printf("target uri: %v\n", target)
	if target == "" {
		return "", "", "", "", errMissingAddr
	}

	if !regexConsul.MatchString(target) {
		fmt.Println("url 匹配错误！")
		return "", "", "", "", errAddrMisMatch
	}

	groups := regexConsul.FindStringSubmatch(target)
	host = groups[1]
	port = groups[2]
	name = groups[3]
	tag = groups[4]
	if port == "" {
		port = defaultPort
	}

	return host, port, name, tag, nil
}

func (cb *consulBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	fmt.Printf("calling consul build \n")
	fmt.Printf("target: %+v\n", target)

	host, port, name, tag, err := parseTarget(fmt.Sprintf("%s%s", target.URL.Host, target.URL.Path))
	if err != nil {
		return nil, err
	}

	cr := &consulResolver{
		address:              fmt.Sprintf("%s%s", host, port),
		cc:                   cc,
		name:                 name,
		tag:                  tag,
		disableServiceConfig: opts.DisableServiceConfig,
		lastIndex:            0,
	}

	cr.wg.Add(1)
	go cr.watcher()
	return cr, nil
}

func (cb *consulBuilder) Scheme() string {
	return "consul"
}

func (cr *consulResolver) ResolveNow(opt resolver.ResolveNowOptions) {

}

func (cr *consulResolver) Close() {

}

// 从 consul 中发现服务
func (cr *consulResolver) watcher() {
	fmt.Printf("calling consul watcher \n")
	config := api.DefaultConfig()
	config.Address = cr.address
	client, err := api.NewClient(config)
	if err != nil {
		fmt.Printf("error create consul client: %v\n", err)
		return
	}

	for {
		// 只获取健康的 service
		services, metainfo, err := client.Health().Service(cr.name, cr.tag, true, &api.QueryOptions{
			WaitIndex: cr.lastIndex,  // 同步点，这个调用将一直阻塞直到有新的更新
		})

		if err != nil {
			fmt.Printf("error retrieving instances from Consul: %v", err)
		}

		cr.lastIndex = metainfo.LastIndex
		var newAddrs []resolver.Address
		for _, service := range services {
			// 健康的链接地址
			fmt.Println("service.Service.Address ==> ", service.Service.Address, "service.Service.Port ==> ", service.Service.Port)
			addr := fmt.Sprintf("%v:%v", service.Service.Address, service.Service.Port)
			newAddrs = append(newAddrs, resolver.Address{Addr: addr})
		}
		fmt.Printf("adding service addrs \n")
		// newAddrs = append(newAddrs, resolver.Address{
		// 	Addr: "127.0.0.1:50051",
		// })
		fmt.Printf("newAddrs: %+v\n", newAddrs)

		cr.cc.UpdateState(resolver.State{
			Addresses: newAddrs,
		})

		// cr.cc.NewAddress(newAddrs)
		cr.cc.NewServiceConfig(cr.name)
	}

}
