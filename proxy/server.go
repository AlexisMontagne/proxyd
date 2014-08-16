package proxy

import (
	"github.com/AlexisMontagne/proxyd/loadbalance"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type ProxyServer struct {
	Port         int
	Ip           string
	UseResolver  bool
	LoadBalancer loadbalance.LoadBalancer
}

func NewProxyServer(Port int, Ip string, forwardProxy bool, loadBalancer loadbalance.LoadBalancer) (*ProxyServer, error) {
	server := &ProxyServer{
		Port,
		Ip,
		forwardProxy,
		loadBalancer,
	}

	return server, nil
}

func (s *ProxyServer) ProxyResolver(endURL *url.URL) (*url.URL, error) {
	return s.LoadBalancer.NextEndpoint(endURL)
}

func (server *ProxyServer) ListenAndServe() {
	http.ListenAndServe(":"+strconv.Itoa(server.Port), http.HandlerFunc(server.ProxyHandler))
}

func (s *ProxyServer) ProxyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Host)
	log.Println(r.RequestURI)
	log.Println(r.RemoteAddr)

	var defaultTransport http.RoundTripper = nil

	if s.UseResolver {
		defaultTransport = &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) { return s.ProxyResolver(r.URL) },
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
		}
	}

	client := &http.Client{
		Transport: defaultTransport,
	}
	r.URL.Scheme = strings.Map(unicode.ToLower, r.URL.Scheme)
	r.RequestURI = ""

	resp, err := client.Do(r)

	defer resp.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}

	w.WriteHeader(resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}
	w.Write(body)
}
