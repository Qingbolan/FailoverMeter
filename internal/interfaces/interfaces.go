package interfaces

import (
    "context"
    "counter_service/kitex_gen/counter_service"
)

type CounterServer interface {
    IncrementCounter(ctx context.Context, req *counter_service.Request) (*counter_service.Response, error)
}