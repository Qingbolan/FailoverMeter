package handler

import (
    "context"
    "counter_service/internal/interfaces"
    "counter_service/kitex_gen/counter_service"
)

type CounterServiceImpl struct {
    Server interfaces.CounterServer
}

func (s *CounterServiceImpl) IncrementCounter(ctx context.Context, req *counter_service.Request) (resp *counter_service.Response, err error) {
    return s.Server.IncrementCounter(ctx, req)
}