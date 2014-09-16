package roundrobin

import (
	"github.com/coreos/go-etcd/etcd"
	"log"
	"net/http"
	"net/url"
)

const Timeout uint64 = 360

type RoundRobinBalancer struct {
	Etcd *etcd.Client
}

func NewBalancer(etcdURLs []string) (*RoundRobinBalancer, error) {
	return &RoundRobinBalancer{etcd.NewClient(etcdURLs)}, nil
}

func (b *RoundRobinBalancer) NextProxyChoose(key string) (*etcd.Node, string) {
	previous, err := b.Etcd.Get(key, false, false)
	previousKey := ""
	var previousNode *etcd.Node = nil

	if err == nil {
		previousNode = previous.Node
		previousKey = previousNode.Value
	}

	log.Println(previousNode)

	resp, err := b.Etcd.Get("/proxyd/endpoints", true, true)

	if err != nil {
		panic("Endpoint key doesn't exist")
	}
	if len(resp.Node.Nodes) == 0 {
		panic("No endpoint available")
	}

	proxyIp := ""
	if previousKey == "" {
		proxyIp = resp.Node.Nodes[0].Value
	} else {
		id := -1
		for idx, node := range resp.Node.Nodes {
			if node.Value == previousKey {
				id = idx
				break
			}
		}
		proxyIp = resp.Node.Nodes[(id+1)%len(resp.Node.Nodes)].Value
	}
	return previousNode, proxyIp
}

func (b *RoundRobinBalancer) NextEndpoint(request *http.Request) (*url.URL, error) {
	key := "proxyd/roundrobin/targets/" + url.QueryEscape(request.URL.Host)
	log.Println(key)
	currentNode, nextIp := b.NextProxyChoose(key)
	var err error = nil

	if currentNode != nil {
		_, err = b.Etcd.CompareAndSwap(
			key, nextIp, Timeout, currentNode.Value, currentNode.ModifiedIndex)
	} else {
		b.Etcd.Set(key, nextIp, Timeout)
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return url.Parse("http://" + nextIp)
}
