package translate

import (
	"fmt"
	"net/http"
	"time"
)

var client *http.Client

func GetHttpClient() *http.Client {
	if client != nil {
		return client
	}
	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
	}
	defaultTransport := *defaultTransportPointer
	defaultTransport.MaxIdleConns = 20
	defaultTransport.MaxIdleConnsPerHost = 20

	client = &http.Client{}

	client.Transport = &defaultTransport
	client.Timeout = 3 * time.Second
	return client
}
