package client

import (
    "bufio"
    "context"
    "fmt"
    "os"
    "strconv"
    "strings"
    "time"

    "github.com/cloudwego/kitex/client"
    "counter_service/kitex_gen/counter_service"
    "counter_service/kitex_gen/counter_service/counterservice"
    "counter_service/logger"
)

type Client struct {
    id     int
    client counterservice.Client
}

func NewClient(id int, serverAddr string) *Client {
    var c counterservice.Client
    var err error
    for i := 0; i < 5; i++ {
        c, err = counterservice.NewClient("CounterService", 
            client.WithHostPorts(serverAddr),
            client.WithConnectTimeout(1*time.Second),
            client.WithRPCTimeout(3*time.Second))
        if err == nil {
            logger.Info("Successfully connected to server")
            break
        }
        logger.Warn("Failed to create client (attempt %d): %v. Retrying...", i+1, err)
        // time.Sleep(2 * time.Second)
    }
    if err != nil {
        logger.Fatal("Failed to create client after 5 attempts: %v", err)
    }
    return &Client{id: id, client: c}
}

func (c *Client) IncrementCounter(incrementBy int32) (*counter_service.Response, error) {
    req := &counter_service.Request{IncrementBy: incrementBy}
    logger.Debug("Sending increment request: %v", req)
    
    var resp *counter_service.Response
    var err error
    for i := 0; i < 3; i++ {
        ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
        resp, err = c.client.IncrementCounter(ctx, req)
        cancel()
        if err == nil {
            logger.Info("Received response: %v", resp)
            return resp, nil
        }
        logger.Warn("Failed to increment counter (attempt %d): %v. Retrying...", i+1, err)
        time.Sleep(1 * time.Second)
    }
    logger.Error("Failed to increment counter after 3 attempts: %v", err)
    return nil, err
}

func (c *Client) Run() {
    reader := bufio.NewReader(os.Stdin)

    for {
        fmt.Print("Enter increment value (or 'q' to quit): ")
        input, _ := reader.ReadString('\n')
        input = strings.TrimSpace(input)

        if input == "q" {
            break
        }

        increment, err := strconv.Atoi(input)
        if err != nil {
            logger.Warn("Invalid input. Please enter a number.")
            continue
        }

        resp, err := c.IncrementCounter(int32(increment))
        if err != nil {
            logger.Error("Error: %v", err)
            continue
        }

        fmt.Printf("Response: Count = %d, Message = %s\n", resp.Count, resp.Message)
    }
}