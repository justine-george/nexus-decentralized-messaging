package webrtc

import (
	"log"

	"github.com/justine-george/nexus-decentralized-messaging/internal/peer"
	"github.com/pion/webrtc/v3"
)

type Connection struct {
	peer *peer.Peer
	pc   *webrtc.PeerConnection
}

func NewConnection(p *peer.Peer) *Connection {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	pc, err := webrtc.NewPeerConnection(config)
	if err != nil {
		log.Printf("Failed to create peer connection: %v", err)
		return nil
	}

	conn := &Connection{
		peer: p,
		pc:   pc,
	}

	conn.setupDataChannel()
	return conn
}

func (c *Connection) setupDataChannel() {
	dc, err := c.pc.CreateDataChannel("chat", nil)
	if err != nil {
		log.Printf("Failed to create data channel: %v", err)
		return
	}

	dc.OnOpen(func() {
		log.Println("Data channel opened")
	})

	dc.OnMessage(func(msg webrtc.DataChannelMessage) {
		// Handle incoming WebRTC messages
		c.peer.HandleMessage(msg.Data)
	})
}

// Add methods for handling ICE candidates, session descriptions, etc.
