package random

import (
	"github.com/coreos/go-etcd/etcd"
	"math/rand"
	"net/http"
	"net/url"
)

type RandomBalancer struct {
	Etcd       *etcd.Client
	Randomizer *rand.Rand
}

func NewBalancer(etcdURLs []string) (*RandomBalancer, error) {
	return &RandomBalancer{
		etcd.NewClient(etcdURLs),
		rand.New(rand.NewSource(10)),
	}, nil
}

func (b *RandomBalancer) NextEndpoint(request *http.Request) (*url.URL, error) {
	resp, _ := b.Etcd.Get("/proxyd/endpoints", true, true)
	ips := []string{}

	if len(resp.Node.Nodes) == 0 {
		panic("No endpoint available")
	}

	for _, i := range resp.Node.Nodes {
		ips = append(ips, i.Value)
	}
	item := b.Randomizer.Int() % len(ips)
	return url.Parse("http://" + ips[item])
}
