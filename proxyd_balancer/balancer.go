package main

import (
	"flag"
	"github.com/AlexisMontagne/proxyd/loadbalance/roundrobin"
	"github.com/AlexisMontagne/proxyd/proxy"
	"log"
	"strconv"
)

var (
	Port = flag.Int("port", 1081, "listen on this port")
	Ip   = flag.String("ip", "127.0.0.1", "IP to reach this endpoint")
	Ttl  = flag.Int("ttl", 5, "TTL to the endpoint")
)

func main() {
	flag.Parse()
	log.Println("Listen 0.0.0.0:" + strconv.Itoa(*Port))
	balancer, _ := roundrobin.NewBalancer([]string{"http://127.0.0.1:4001"})
	proxy, _ := proxy.NewProxyServer(*Port, "0.0.0.0", true, balancer)
	proxy.ListenAndServe()
}
