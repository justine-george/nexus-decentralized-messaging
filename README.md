# Nexus Decentralized Messaging

Nexus is a decentralized messaging service that ensures secure, peer-to-peer communication with end-to-end encryption. Built using Go and WebRTC, Nexus offers a scalable infrastructure for seamless and private messaging.

## Features

- **WebRTC Integration**: Real-time communication using WebRTC.
- **gRPC Directory Service**: Manages peer registration and discovery.
- **JSON Messaging**: Messages are encoded in JSON format for simplicity and interoperability.

## Getting Started

### Prerequisites

- Go 1.22.6 or later
- Protocol Buffers compiler (`protoc`)
- Web browser with WebRTC support

### Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/justine-george/nexus-decentralized-messaging.git
    cd nexus-decentralized-messaging
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

3. Generate the gRPC code from the proto files:
    ```sh
    ./scripts/generate_proto.sh
    ```

### Running the Application

1. Start the directory service:
    ```sh
    go run cmd/server/main.go
    ```

2. Open `web/templates/index.html` in your web browser to start the chat interface.

## Project Structure

- **cmd/server/main.go**: Entry point for the server application.
- **internal/directory/service.go**: Implementation of the directory service.
- **pkg/message/message.go**: Message struct and JSON serialization/deserialization.
- **proto/directory.proto**: Protocol Buffers definitions for the directory service.
- **web/templates/index.html**: Frontend HTML for the chat application.
- **scripts/generate_proto.sh**: Script to generate gRPC code from proto files.

## Usage

### Registering a Peer

Peers can register with the directory service to be discoverable by other peers.

```go
dirService.RegisterPeer(context.Background(), &pb.RegisterRequest{
    Id:      "peer1",
    Address: "192.168.1.2",
})
```

### Sending a Message

Messages are created and serialized to JSON format for transmission.

```go
msg := message.New("chat", "peer1", "peer2", "Hello, World!")
jsonData, err := msg.ToJSON()
if err != nil {
    log.Fatal(err)
}
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License.