

# FailoverMeter: Resilient Distributed Counter Service

FailoverMeter is a distributed counter service built with a primary-backup architecture, implemented in Go using the Kitex RPC framework.

## Features

- Primary-backup architecture for high availability
- Automatic failover and recovery
- Real-time counter updates
- Concurrent request handling
- Customizable logging levels

## Project Structure

```bash
distcount/
├── cmd/
│   ├── server/
│   │   └── main.go
│   └── client/
│       └── main.go
├── internal/
│   ├── server/
│   │   └── server.go
│   ├── client/
│   │   └── client.go
│   ├── handler/
│   │   └── handler.go
│   └── interfaces/
│       └── interfaces.go
├── logger/
│   └── logger.go
├── kitex_gen/
│   └── ...
└── counter_service.thrift
```

## Prerequisites

- Go 1.16 or higher
- Kitex framework
- Thrift compiler

## Installation

1. Clone the repository:
```

   git clone https://github.com/yourusername/distcount.git
   cd distcount

```

2. Install dependencies:
```

   go mod tidy

```

3. Generate Kitex code:
```

   kitex -module "distcount" counter_service.thrift

```

## Compilation

1. Build the server:
```

   go build -o server cmd/server/main.go

```

2. Build the client:
```

   go build -o client cmd/client/main.go

```

## Deployment

1. Start the backup server:
```

   ./server -role backup -port 8081 -log debug

```

2. Start the primary server:
```

   ./server -role primary -port 8080 -backup localhost:8081 -log debug

```

3. Run the client:
```

   ./client

```

## Usage

- After starting the client, enter a number to increment the counter.
- Enter 'q' to quit the client.

## Configuration

Configure the server using command-line arguments:

- `-role`: Server role, either "primary" or "backup"
- `-port`: Port for the server to listen on
- `-backup`: Address of the backup server (for primary server only)
- `-log`: Log level, can be "debug", "info", "warn", or "error"

## API

The service exposes a single RPC method:

```thrift
service CounterService {
    Response IncrementCounter(1: Request req)
}

struct Request {
    1: i32 IncrementBy
}

struct Response {
    1: i32 Count
    2: string Message
}
```

## Development

To contribute to DistCount:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## Troubleshooting

If you encounter connection issues:

1. Ensure no firewalls are blocking the connection.
2. Check the network settings of both server and client.
3. Make sure the backup server is started before the primary server.

For compilation issues:

1. Verify that Go and Kitex are correctly installed.
2. Check your GOPATH and GOROOT settings.
3. Ensure all dependencies are correctly installed.

## Performance Considerations

- The system uses goroutines for concurrent request handling, but be mindful of resource usage with a high number of concurrent clients.
- Adjust the heartbeat interval based on your network conditions and requirements.
- For production environments, consider implementing rate limiting and additional security measures.

## Future Enhancements

- Implement data persistence
- Add support for multiple backup servers
- Introduce a web-based dashboard for monitoring
- Implement client-side load balancing

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## Acknowledgments

- Kitex framework developers
- The Go community for their excellent documentation and support

## Contact

Your Name - silan-hu@comp.nus.edu.sg

Project Link: [https://github.com/Qingbolan/FailoverMeter](https://github.com/Qingbolan/FailoverMeter)
