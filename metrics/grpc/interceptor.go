package grpc

import (
	"context"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type grpcType string

const (
	Unary        grpcType = "unary"
	ClientStream grpcType = "client_stream"
	ServerStream grpcType = "server_stream"
	BidiStream   grpcType = "bidi_stream"
)

func NewMetricsUnaryInterceptor(observer GrpcObserver) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (result interface{}, err error) {
		defer func(time time.Time) {
			serviceName, methodName := parseMethodName(info.FullMethod)
			// observe request
			observer.ObserveRequest(time, serviceName, string(Unary), methodName, status.Code(err).String())
		}(time.Now())

		// exec handler
		result, err = handler(ctx, req)

		return result, err
	}
}

func parseMethodName(fullMethodName string) (serviceName string, methodName string) {
	fullMethodName = strings.TrimPrefix(fullMethodName, "/") // remove leading slash
	arrStr := strings.Split(fullMethodName, "/")
	if len(arrStr) > 1 {
		return arrStr[0], arrStr[len(arrStr)-1]
	} else {
		return "unknown", fullMethodName
	}
}
