package random

import (
	"github.com/coreos/go-etcd/etcd"
	"log"
	"math/rand"
	"net/url"
)

type RandomBalancer struct {
	Etcd       *etcd.Client
	Randomizer *rand.Rand
}

func NewBalancer(etcdURLs []string) (*RandomBalancer, error) {
	return &RandomBalancer{etcd.NewClient(etcdURLs), rand.New(rand.NewSource(10))}, nil
}

func (b *RandomBalancer) NextEndpoint(endURL *url.URL) (*url.URL, error) {
	resp, _ := b.Etcd.Get("/proxyd/endpoints", true, true)
	ips := []string{}
	for _, i := range resp.Node.Nodes {
		ips = append(ips, i.Value)
	}
	log.Println(ips)
	item := b.Randomizer.Int() % len(ips)
	return url.Parse("http://" + ips[item])
}
