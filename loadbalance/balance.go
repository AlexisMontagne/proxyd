package loadbalance

import "net/url"

type LoadBalancer interface {
	NextEndpoint(endURL *url.URL) (*url.URL, error)
}
