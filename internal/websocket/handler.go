package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/justine-george/nexus-decentralized-messaging/internal/directory"
	"github.com/justine-george/nexus-decentralized-messaging/internal/peer"
	"github.com/justine-george/nexus-decentralized-messaging/internal/webrtc"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleConnection(dirService *directory.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}
		defer conn.Close()

		p := peer.New(conn)
		rtcConn := webrtc.NewConnection(p)

		// Register peer with directory service
		dirService.RegisterPeer(r.Context(), &pb.RegisterRequest{Id: p.ID, Address: r.RemoteAddr})

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}

			switch messageType {
			case websocket.TextMessage:
				// Handle text message (e.g., chat or signaling)
				handleTextMessage(p, message)
			case websocket.BinaryMessage:
				// Handle binary message (e.g., file transfer)
				handleBinaryMessage(p, message)
			}
		}

		// Unregister peer when connection closes
		dirService.UnregisterPeer(r.Context(), &pb.UnregisterRequest{Id: p.ID})
	}
}

func handleTextMessage(p *peer.Peer, message []byte) {
	// Implement text message handling (e.g., JSON parsing, chat routing)
}

func handleBinaryMessage(p *peer.Peer, message []byte) {
	// Implement binary message handling (e.g., file transfer)
}
