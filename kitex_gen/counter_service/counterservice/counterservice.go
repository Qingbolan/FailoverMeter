// Code generated by Kitex v0.10.3. DO NOT EDIT.

package counterservice

import (
	"context"
	counter_service "counter_service/kitex_gen/counter_service"
	"errors"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
)

var errInvalidMessageType = errors.New("invalid message type for service method handler")

var serviceMethods = map[string]kitex.MethodInfo{
	"IncrementCounter": kitex.NewMethodInfo(
		incrementCounterHandler,
		newCounterServiceIncrementCounterArgs,
		newCounterServiceIncrementCounterResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
}

var (
	counterServiceServiceInfo                = NewServiceInfo()
	counterServiceServiceInfoForClient       = NewServiceInfoForClient()
	counterServiceServiceInfoForStreamClient = NewServiceInfoForStreamClient()
)

// for server
func serviceInfo() *kitex.ServiceInfo {
	return counterServiceServiceInfo
}

// for stream client
func serviceInfoForStreamClient() *kitex.ServiceInfo {
	return counterServiceServiceInfoForStreamClient
}

// for client
func serviceInfoForClient() *kitex.ServiceInfo {
	return counterServiceServiceInfoForClient
}

// NewServiceInfo creates a new ServiceInfo containing all methods
func NewServiceInfo() *kitex.ServiceInfo {
	return newServiceInfo(false, true, true)
}

// NewServiceInfo creates a new ServiceInfo containing non-streaming methods
func NewServiceInfoForClient() *kitex.ServiceInfo {
	return newServiceInfo(false, false, true)
}
func NewServiceInfoForStreamClient() *kitex.ServiceInfo {
	return newServiceInfo(true, true, false)
}

func newServiceInfo(hasStreaming bool, keepStreamingMethods bool, keepNonStreamingMethods bool) *kitex.ServiceInfo {
	serviceName := "CounterService"
	handlerType := (*counter_service.CounterService)(nil)
	methods := map[string]kitex.MethodInfo{}
	for name, m := range serviceMethods {
		if m.IsStreaming() && !keepStreamingMethods {
			continue
		}
		if !m.IsStreaming() && !keepNonStreamingMethods {
			continue
		}
		methods[name] = m
	}
	extra := map[string]interface{}{
		"PackageName": "counter_service",
	}
	if hasStreaming {
		extra["streaming"] = hasStreaming
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.10.3",
		Extra:           extra,
	}
	return svcInfo
}

func incrementCounterHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*counter_service.CounterServiceIncrementCounterArgs)
	realResult := result.(*counter_service.CounterServiceIncrementCounterResult)
	success, err := handler.(counter_service.CounterService).IncrementCounter(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newCounterServiceIncrementCounterArgs() interface{} {
	return counter_service.NewCounterServiceIncrementCounterArgs()
}

func newCounterServiceIncrementCounterResult() interface{} {
	return counter_service.NewCounterServiceIncrementCounterResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) IncrementCounter(ctx context.Context, req *counter_service.Request) (r *counter_service.Response, err error) {
	var _args counter_service.CounterServiceIncrementCounterArgs
	_args.Req = req
	var _result counter_service.CounterServiceIncrementCounterResult
	if err = p.c.Call(ctx, "IncrementCounter", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
