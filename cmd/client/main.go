package main

import (
    "counter_service/internal/client"
    "counter_service/logger"
)

func main() {
    logger.SetLogLevel(logger.INFO)
    c := client.NewClient(1, "localhost:8080")
    c.Run()
}