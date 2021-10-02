package queue

import (
	"encoding/json"
	"time"
)

type Header struct {
	Id         string    `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Message struct {
	Header Header `json:"header"`
	Body   []byte `json:"body"`
}

type SignedMessage struct {
	Message   Message `json:"message"`
	Signature string  `json:"signature"`
}

func (m SignedMessage) Serialize() ([]byte, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}

func Deserialize(rawMessage string) (message SignedMessage, err error) {
	err = json.Unmarshal([]byte(rawMessage), &message)
	if err != nil {
		return SignedMessage{}, err
	}
	return
}
