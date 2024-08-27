package peer

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/justine-george/nexus-decentralized-messaging/pkg/message"
	"github.com/pion/webrtc/v3"
)

type Peer struct {
	ID          string
	Name        string
	WSConn      *websocket.Conn
	RTCPeerConn *webrtc.PeerConnection
	DataChannel *webrtc.DataChannel
	messageChan chan *message.Message
	mu          sync.Mutex
}

func New(id string, wsConn *websocket.Conn) *Peer {
	return &Peer{
		ID:          id,
		Name:        "Anonymous",
		WSConn:      wsConn,
		messageChan: make(chan *message.Message, 100),
	}
}

func (p *Peer) SetRTCPeerConnection(rtcPeerConn *webrtc.PeerConnection) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.RTCPeerConn = rtcPeerConn
}

func (p *Peer) SetDataChannel(dataChannel *webrtc.DataChannel) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.DataChannel = dataChannel
}

func (p *Peer) SendMessage(msg *message.Message) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.DataChannel != nil && p.DataChannel.ReadyState() == webrtc.DataChannelStateOpen {
		bytes, err := msg.ToJSON()
		if err != nil {
			return err
		}
		return p.DataChannel.Send(bytes)
	}

	return p.WSConn.WriteJSON(msg)
}

func (p *Peer) HandleMessage(data []byte) {
	msg, err := message.FromJSON(data)
	if err != nil {
		log.Printf("Error parsing message from peer %s: %v", p.ID, err)
		return
	}

	p.messageChan <- msg
}

func (p *Peer) Start() {
	go p.readPump()
	go p.writePump()
}

func (p *Peer) readPump() {
	defer func() {
		p.WSConn.Close()
	}()

	for {
		_, message, err := p.WSConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message from peer %s: %v", p.ID, err)
			}
			break
		}
		p.HandleMessage(message)
	}
}

func (p *Peer) writePump() {
	for msg := range p.messageChan {
		err := p.SendMessage(msg)
		if err != nil {
			log.Printf("Error sending message to peer %s: %v", p.ID, err)
			return
		}
	}
}

func (p *Peer) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.RTCPeerConn != nil {
		p.RTCPeerConn.Close()
	}
	if p.DataChannel != nil {
		p.DataChannel.Close()
	}
	close(p.messageChan)
	p.WSConn.Close()
}
