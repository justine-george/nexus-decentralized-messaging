package message

import "encoding/json"

type Message struct {
	Type    string `json:"type"`
	From    string `json:"from"`
	To      string `json:"to,omitempty"`
	Content string `json:"content"`
}

func New(msgType, from, to, content string) *Message {
	return &Message{
		Type:    msgType,
		From:    from,
		To:      to,
		Content: content,
	}
}

func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func FromJSON(data []byte) (*Message, error) {
	var m Message
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
