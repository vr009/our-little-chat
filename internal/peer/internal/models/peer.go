package models

import "github.com/google/uuid"

type Peer struct {
	PeerID     uuid.UUID
	MsgsToRecv chan []Message
	MsgToSend  chan *Message
	Connected  bool
}

func GetPeerFromMessage(msg *Message) *Peer {
	return &Peer{
		PeerID:     msg.SenderID,
		MsgToSend:  make(chan *Message, 100),
		MsgsToRecv: make(chan []Message, 100),
		Connected:  false,
	}
}
