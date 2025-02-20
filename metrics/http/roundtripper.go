package http

import (
	"net/http"
	"time"
)

type NetHttpMetricsRoundTripper struct {
	Next     http.RoundTripper
	Observer OutboundHttpObserver
}

func NewMetricsRoundTripper(next http.RoundTripper, observer OutboundHttpObserver) *NetHttpMetricsRoundTripper {
	return &NetHttpMetricsRoundTripper{
		Next:     next,
		Observer: observer,
	}
}

func (m *NetHttpMetricsRoundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
	defer func(time time.Time) {
		// observe request
		if res != nil && req != nil {
			m.Observer.ObserveRequest(time, err == nil, res.StatusCode, req.Method, req.URL.Host)
		}
	}(time.Now())

	// execute request
	res, err = m.Next.RoundTrip(req)

	return res, err
}
