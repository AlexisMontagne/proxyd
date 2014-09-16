package loadbalance

import (
	"net/http"
	"net/url"
)

type LoadBalancer interface {
	NextEndpoint(request *http.Request) (*url.URL, error)
}
