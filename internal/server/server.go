package server

import (
    "context"
    "fmt"
    "net"
    "sync"
    "time"

    "github.com/cloudwego/kitex/client"
    "github.com/cloudwego/kitex/server"
    "counter_service/internal/handler"
    "counter_service/kitex_gen/counter_service"
    "counter_service/kitex_gen/counter_service/counterservice"
    "counter_service/logger"
)

type Server interface {
    Start() error
    Shutdown()
}

type PrimaryServer struct {
    port       int
    backupAddr string
    backup     counterservice.Client
    isMaster   bool
    mu         sync.Mutex
    count      int32
    srv        server.Server
    handler    *handler.CounterServiceImpl
}

func NewPrimaryServer(port int, backupAddr string) *PrimaryServer {
    s := &PrimaryServer{
        port:       port,
        backupAddr: backupAddr,
        isMaster:   true,
    }
    s.handler = &handler.CounterServiceImpl{Server: s}
    return s
}

type BackupServer struct {
    port    int
    count   int32
    mu      sync.Mutex
    srv     server.Server
    handler *handler.CounterServiceImpl
}

func NewBackupServer(port int) *BackupServer {
    s := &BackupServer{
        port: port,
    }
    s.handler = &handler.CounterServiceImpl{Server: s}
    return s
}

// Implement the CounterServer interface for PrimaryServer
func (s *PrimaryServer) IncrementCounter(ctx context.Context, req *counter_service.Request) (*counter_service.Response, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    if s.isMaster {
        s.count += req.IncrementBy
        logger.Info("Incremented counter by %d. New count: %d", req.IncrementBy, s.count)
        return &counter_service.Response{Count: s.count, Message: "Processed by primary"}, nil
    }
    logger.Info("Forwarding increment request to backup server")
    return s.backup.IncrementCounter(ctx, req)
}

// Implement the CounterServer interface for BackupServer
func (s *BackupServer) IncrementCounter(ctx context.Context, req *counter_service.Request) (*counter_service.Response, error) {
    logger.Debug("Received IncrementCounter request: %v", req)
    s.mu.Lock()
    defer s.mu.Unlock()
    s.count += req.IncrementBy
    logger.Info("Incremented counter by %d. New count: %d", req.IncrementBy, s.count)
    resp := &counter_service.Response{Count: s.count, Message: "Processed by backup"}
    logger.Debug("Sending response: %v", resp)
    return resp, nil
}

func (s *PrimaryServer) Start() error {
    var err error
    for i := 0; i < 10; i++ {
        logger.Info("Attempting to create backup client (attempt %d)", i+1)
        s.backup, err = counterservice.NewClient("CounterService", 
            client.WithHostPorts(s.backupAddr),
            client.WithConnectTimeout(1*time.Second),
            client.WithRPCTimeout(3*time.Second))
        if err == nil {
            logger.Info("Successfully created backup client")
            break
        }
        logger.Warn("Failed to create backup client: %v. Retrying...", err)
        time.Sleep(3 * time.Second)
    }
    if err != nil {
        return fmt.Errorf("failed to create backup client after 10 attempts: %v", err)
    }

    go s.startHeartbeat()

    s.srv = server.NewServer(server.WithServiceAddr(&net.TCPAddr{IP: net.IPv4(0, 0, 0, 0), Port: s.port}))
    logger.Info("Registering service...")
    err = s.srv.RegisterService(counterservice.NewServiceInfo(), s.handler)
    if err != nil {
        return fmt.Errorf("failed to register service: %v", err)
    }

    logger.Info("Primary server starting on port %d", s.port)
    return s.srv.Run()
}

func (s *BackupServer) Start() error {
    s.srv = server.NewServer(server.WithServiceAddr(&net.TCPAddr{IP: net.IPv4(0, 0, 0, 0), Port: s.port}))
    logger.Info("Registering service...")
    err := s.srv.RegisterService(counterservice.NewServiceInfo(), s.handler)
    if err != nil {
        return fmt.Errorf("failed to register service: %v", err)
    }
    logger.Info("Backup server starting on port %d", s.port)
    return s.srv.Run()
}

func (s *PrimaryServer) Shutdown() {
    if s.srv != nil {
        logger.Info("Shutting down primary server...")
        s.srv.Stop()
    }
}

func (s *BackupServer) Shutdown() {
    if s.srv != nil {
        logger.Info("Shutting down backup server...")
        s.srv.Stop()
    }
}

func (s *PrimaryServer) startHeartbeat() {
    failCount := 0
    for {
        err := s.checkBackup()
        s.mu.Lock()
        if err != nil {
            failCount++
            logger.Warn("Backup server heartbeat failed (%d): %v", failCount, err)
            if failCount >= 3 && !s.isMaster {
                logger.Info("Backup server is down, primary taking over")
                s.isMaster = true
            }
        } else {
            if failCount > 0 {
                logger.Info("Backup server recovered after %d failures", failCount)
            }
            failCount = 0
            if s.isMaster {
                logger.Info("Backup server is up, primary stepping down")
                s.isMaster = false
            }
        }
        s.mu.Unlock()
        time.Sleep(15 * time.Second)
    }
}

func (s *PrimaryServer) checkBackup() error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    req := &counter_service.Request{IncrementBy: 0}
    for i := 0; i < 3; i++ {
        logger.Debug("Attempting to check backup (attempt %d)...", i+1)
        _, err := s.backup.IncrementCounter(ctx, req)
        if err == nil {
            logger.Debug("Backup check successful")
            return nil
        }
        logger.Warn("Attempt %d: Backup check failed: %v. Retrying...", i+1, err)
        time.Sleep(2 * time.Second)
    }
    return fmt.Errorf("backup check failed after 3 attempts")
}