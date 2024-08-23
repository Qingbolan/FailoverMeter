package main

import (
    "flag"
    "os"
    "os/signal"
    "syscall"
    "time"

    "counter_service/internal/server"
    "counter_service/logger"
)

func main() {
    role := flag.String("role", "primary", "Server role: primary or backup")
    port := flag.Int("port", 8080, "Server port")
    backupAddr := flag.String("backup", "localhost:8081", "Backup server address")
    logLevel := flag.String("log", "info", "Log level: debug, info, warn, error")
    flag.Parse()

    // 设置日志级别
    switch *logLevel {
    case "debug":
        logger.SetLogLevel(logger.DEBUG)
    case "info":
        logger.SetLogLevel(logger.INFO)
    case "warn":
        logger.SetLogLevel(logger.WARN)
    case "error":
        logger.SetLogLevel(logger.ERROR)
    default:
        logger.SetLogLevel(logger.INFO)
    }

    logger.Info("Starting server with role: %s, port: %d, backup: %s", *role, *port, *backupAddr)

    var s server.Server
    if *role == "primary" {
        logger.Info("Waiting for backup server to start...")
        time.Sleep(10 * time.Second)
        s = server.NewPrimaryServer(*port, *backupAddr)
    } else {
        s = server.NewBackupServer(*port)
    }

    go func() {
        if err := s.Start(); err != nil {
            logger.Fatal("%s server error: %v", *role, err)
        }
    }()

    // 等待中断信号以优雅地关闭服务器
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    logger.Info("Shutting down server...")
    s.Shutdown()
}