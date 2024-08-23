package logger

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "runtime"
    "time"
)

type LogLevel int

const (
    DEBUG LogLevel = iota
    INFO
    WARN
    ERROR
    FATAL
)

var (
    logger *log.Logger
    level  LogLevel
)

const (
    colorReset  = "\033[0m"
    colorRed    = "\033[31m"
    colorGreen  = "\033[32m"
    colorYellow = "\033[33m"
    colorBlue   = "\033[34m"
    colorPurple = "\033[35m"
    colorCyan   = "\033[36m"
    colorWhite  = "\033[37m"
)

func init() {
    logger = log.New(os.Stdout, "", 0)
    level = INFO // 默认日志级别
}

func SetLogLevel(l LogLevel) {
    level = l
}

func logf(l LogLevel, format string, v ...interface{}) {
    if l < level {
        return
    }

    now := time.Now().Format("2006-01-02 15:04:05.000")
    _, file, line, _ := runtime.Caller(2)
    file = filepath.Base(file)

    var levelStr, colorCode string
    switch l {
    case DEBUG:
        levelStr = "DEBUG"
        colorCode = colorCyan
    case INFO:
        levelStr = "INFO "
        colorCode = colorGreen
    case WARN:
        levelStr = "WARN "
        colorCode = colorYellow
    case ERROR:
        levelStr = "ERROR"
        colorCode = colorRed
    case FATAL:
        levelStr = "FATAL"
        colorCode = colorPurple
    }

    msg := fmt.Sprintf(format, v...)
    logMsg := fmt.Sprintf("%s%s [%s] %s:%d - %s%s", colorCode, now, levelStr, file, line, msg, colorReset)
    logger.Println(logMsg)

    if l == FATAL {
        os.Exit(1)
    }
}

func Debug(format string, v ...interface{}) {
    logf(DEBUG, format, v...)
}

func Info(format string, v ...interface{}) {
    logf(INFO, format, v...)
}

func Warn(format string, v ...interface{}) {
    logf(WARN, format, v...)
}

func Error(format string, v ...interface{}) {
    logf(ERROR, format, v...)
}

func Fatal(format string, v ...interface{}) {
    logf(FATAL, format, v...)
}