package main

import (
	"flag"
	"github.com/AlexisMontagne/proxyd/discovery"
	"github.com/AlexisMontagne/proxyd/proxy"
	"log"
	"strconv"
)

var (
	Port = flag.Int("port", 1080, "listen on this port")
	Ip   = flag.String("ip", "127.0.0.1", "IP to reach this endpoint")
	Ttl  = flag.Int("ttl", 5, "TTL to the endpoint")
)

func main() {
	flag.Parse()
	log.Println("Listen 0.0.0.0:" + strconv.Itoa(*Port))
	discovery, _ := discovery.NewDiscovery(*Ip, *Ttl, *Port, "http://127.0.0.1:4001")
	go discovery.KeepDeclared()
	proxy, _ := proxy.NewProxyServer(*Port, "0.0.0.0", false, nil)
	proxy.ListenAndServe()
}
