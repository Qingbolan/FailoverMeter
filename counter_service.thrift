namespace go counter_service

struct Request {
    1: i32 IncrementBy
}

struct Response {
    1: i32 Count
    2: string Message
}

service CounterService {
    Response IncrementCounter(1: Request req)
}