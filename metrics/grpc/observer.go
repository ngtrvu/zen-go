package grpc

import "time"

type GrpcObserver interface {
	ObserveRequest(start time.Time, serviceName, grpcType, method, statusCode string)
}
