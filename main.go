package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Peer struct {
	ID   string
	Name string
	Conn *websocket.Conn
}

type Message struct {
	Type    string `json:"type"`
	From    string `json:"from"`
	FromID  string `json:"fromId"`
	To      string `json:"to,omitempty"`
	Content string `json:"content"`
}

var (
	peers = make(map[string]*Peer)
	mutex sync.RWMutex
)

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/", handleHome)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	peer := &Peer{
		ID:   fmt.Sprintf("peer-%d", len(peers)+1),
		Name: fmt.Sprintf("Anonymous-%d", len(peers)+1),
		Conn: conn,
	}

	mutex.Lock()
	peers[peer.ID] = peer
	mutex.Unlock()

	log.Printf("New peer connected: %s (%s)\n", peer.Name, peer.ID)

	// Send the peer its ID
	if err := conn.WriteJSON(Message{Type: "id_assigned", Content: peer.ID}); err != nil {
		log.Printf("Error sending ID to peer: %v", err)
		return
	}

	// Send peer list to the new peer
	sendPeerList(peer)

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}

		msg.FromID = peer.ID
		handleMessage(peer, msg)
	}

	mutex.Lock()
	delete(peers, peer.ID)
	mutex.Unlock()
	log.Printf("Peer disconnected: %s (%s)\n", peer.Name, peer.ID)
	broadcastPeerList()
}

func sendPeerList(peer *Peer) {
	peerList := make(map[string]string)
	mutex.RLock()
	for id, p := range peers {
		if id != peer.ID {
			peerList[id] = p.Name
		}
	}
	mutex.RUnlock()

	msg := Message{
		Type:    "peer_list",
		Content: formatPeerList(peerList),
	}
	if err := peer.Conn.WriteJSON(msg); err != nil {
		log.Printf("Error sending peer list to %s: %v", peer.ID, err)
	}
}

func broadcastPeerList() {
	peerList := make(map[string]string)
	mutex.RLock()
	for id, p := range peers {
		peerList[id] = p.Name
	}
	mutex.RUnlock()

	msg := Message{
		Type:    "peer_list",
		Content: formatPeerList(peerList),
	}

	mutex.RLock()
	for _, peer := range peers {
		if err := peer.Conn.WriteJSON(msg); err != nil {
			log.Printf("Error broadcasting peer list to %s: %v", peer.ID, err)
		}
	}
	mutex.RUnlock()
}

func formatPeerList(peerList map[string]string) string {
	jsonBytes, _ := json.Marshal(peerList)
	return string(jsonBytes)
}

func handleMessage(sender *Peer, msg Message) {
	switch msg.Type {
	case "chat":
		mutex.RLock()
		targetPeer, exists := peers[msg.To]
		senderName := sender.Name
		mutex.RUnlock()

		if exists {
			msg.From = senderName
			msg.FromID = sender.ID
			if err := targetPeer.Conn.WriteJSON(msg); err != nil {
				log.Printf("Error sending message to %s: %v", targetPeer.ID, err)
			}
			// Send the message back to the sender as well
			if err := sender.Conn.WriteJSON(msg); err != nil {
				log.Printf("Error sending message back to sender %s: %v", sender.ID, err)
			}
		} else {
			log.Printf("Peer %s not found\n", msg.To)
		}
	case "set_name":
		oldName := sender.Name
		mutex.Lock()
		sender.Name = msg.Content
		mutex.Unlock()
		log.Printf("Peer %s changed name from %s to %s\n", sender.ID, oldName, sender.Name)
		// Notify all peers about the name change
		nameUpdateMsg := Message{
			Type:    "name_updated",
			FromID:  sender.ID,
			Content: sender.Name,
		}
		mutex.RLock()
		for _, peer := range peers {
			if err := peer.Conn.WriteJSON(nameUpdateMsg); err != nil {
				log.Printf("Error sending name update to %s: %v", peer.ID, err)
			}
		}
		mutex.RUnlock()
		broadcastPeerList()
	case "get_id":
		if err := sender.Conn.WriteJSON(Message{
			Type:    "id_assigned",
			Content: sender.ID,
		}); err != nil {
			log.Printf("Error sending ID to peer %s: %v", sender.ID, err)
		}
	default:
		log.Printf("Unknown message type: %s\n", msg.Type)
	}
}
