package discovery

import (
	"errors"
	"github.com/coreos/go-etcd/etcd"
	"strconv"
	"time"
)

const EndpointsKey = "/proxyd/endpoints"

type Endpoint struct {
	Ip   string
	Port int
}

type Discovery struct {
	Endpoint *Endpoint
	Etcd     *etcd.Client
	Ttl      int
	Key      string
}

func NewDiscovery(ip string, ttl, port int, etcdUrl string) (*Discovery, error) {
	discovery := &Discovery{
		&Endpoint{ip, port},
		etcd.NewClient([]string{etcdUrl}),
		ttl,
		"",
	}

	discovery.Etcd.CreateDir(EndpointsKey, 0)

	return discovery, nil
}

func (e *Endpoint) FormatUrl() string {
	return e.Ip + ":" + strconv.Itoa(e.Port)
}

func (d *Discovery) Declare() (string, error) {
	var err error = nil
	if d.Key != "" {
		_, err = d.Etcd.Update(d.Key, d.Endpoint.FormatUrl(), uint64(d.Ttl))
	} else if d.Key == "" || err != nil {
		resp, _ := d.Etcd.CreateInOrder(EndpointsKey, d.Endpoint.FormatUrl(), uint64(d.Ttl))
		d.Key = resp.Node.Key
	} else {
		err = errors.New("Can't either create or update the key")
	}
	return d.Key, err
}

func (d *Discovery) KeepDeclared() {
	for {
		d.Declare()
		time.Sleep(time.Duration(d.Ttl-1) * time.Second)
	}
}
