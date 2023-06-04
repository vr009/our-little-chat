package models

import (
	"fmt"

	"github.com/google/uuid"
)

type Chat struct {
	ChatID uuid.UUID
	Peers  map[uuid.UUID]*Peer //TODO add sync.Map
}

func (c *Chat) SubscribePeer(newPeer *Peer) error {
	_, ok := c.Peers[newPeer.PeerID]
	if !ok {
		c.Peers[newPeer.PeerID] = newPeer
	}
	c.Peers[newPeer.PeerID].Connected = true
	return nil
}

func (c *Chat) UnsubscribePeer(peer *Peer) error {
	if peer == nil {
		return nil
	}
	peer, ok := c.Peers[peer.PeerID]
	if !ok {
		return fmt.Errorf("peer not found")
	}
	delete(c.Peers, peer.PeerID)
	peer.Connected = false
	return nil
}

func GetChatFromInitialMessage(msg *Message) *Chat {
	return &Chat{
		ChatID: msg.ChatID,
		Peers:  make(map[uuid.UUID]*Peer, 10),
	}
}
