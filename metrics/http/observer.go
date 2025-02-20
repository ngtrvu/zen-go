package http

import "time"

type HttpObserver interface {
	ObserveRequest(start time.Time, statusCode int, method string, endpoint string)
}

type OutboundHttpObserver interface {
	ObserveRequest(start time.Time, success bool, statusCode int, method string, host string)
}
