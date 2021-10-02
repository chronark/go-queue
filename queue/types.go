package queue

import (
	"time"
)

type Header struct {
	Id        string    `fauna:"id" json:"id"`
	Topic     string    `fauna:"topic" json:"topic"`
	CreatedAt time.Time `fauna:"createdAt" json:"createdAt"`
}

type Message struct {
	Header  Header `fauna:"header" json:"header"`
	Payload []byte `fauna:"payload" json:"payload"`
}

type SignedMessage struct {
	Message   Message `fauna:"message" json:"message"`
	Signature string  `fauna:"signature" json:"signature"`
}
