package main

import (
	"github.com/google/uuid"
	websocket2 "github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"our-little-chatik/internal/peer/internal/models"
	"path/filepath"
	"sync"
	"testing"
)

var (
	uuidA       = uuid.New()
	uuidB       = uuid.New()
	uuidChat    = uuid.New()
	textFromA   = "Hi, B!"
	textFromB   = "Hi, A!"
	defaultText = "what's up"
	hellotext   = ""
)

func TestPeer(t *testing.T) {
	path, _ := os.Getwd()
	os.Setenv("CONFIG", filepath.Join(path))
	go main()

	wg := sync.WaitGroup{}
	wg.Add(2)

	t.Logf("user A: %v", uuidA)
	t.Logf("user B: %v", uuidB)

	clientA := websocket2.DefaultDialer
	clientB := websocket2.DefaultDialer

	u := url.URL{Scheme: "ws", Host: ":8080", Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	connA, _, err := clientA.Dial(u.String(), nil)
	if err != nil {
		t.Fatal(err)
	}
	connB, _, err := clientB.Dial(u.String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	HelloMsgFromA := models.Message{
		ChatID:       uuidChat,
		ReceiverID:   uuidB,
		SenderID:     uuidA,
		Payload:      hellotext,
		SessionStart: true,
	}
	HelloMsgFromB := models.Message{
		ChatID:       uuidChat,
		ReceiverID:   uuidA,
		SenderID:     uuidB,
		Payload:      hellotext,
		SessionStart: true,
	}

	msgFromA := models.Message{
		ChatID:     uuidChat,
		ReceiverID: uuidB,
		SenderID:   uuidA,
		Payload:    textFromA,
	}
	msgFromA2 := models.Message{
		ChatID:     uuidChat,
		ReceiverID: uuidB,
		SenderID:   uuidA,
		Payload:    defaultText,
	}

	msgFromB := models.Message{
		ChatID:     uuidChat,
		ReceiverID: uuidA,
		SenderID:   uuidB,
		Payload:    textFromB,
	}

	go func() {
		err = connA.WriteJSON(HelloMsgFromA)
		if err != nil {
			t.Errorf("failed to send msg: %v", err)
		}
		err = connA.WriteJSON(msgFromA)
		if err != nil {
			t.Errorf("failed to send msg: %v", err)
		}
		err = connA.WriteJSON(msgFromA2)
		if err != nil {
			t.Errorf("failed to send msg: %v", err)
		}

		msgRespFromB := &models.Message{}
		err = connA.ReadJSON(msgRespFromB)
		if err != nil {
			t.Error(err)
		}
		if msgRespFromB.Payload != textFromB {
			t.Errorf("the received message is wrong")
		}
		wg.Done()
	}()

	go func() {
		err = connB.WriteJSON(HelloMsgFromB)
		if err != nil {
			t.Errorf("failed to send msg: %v", err)
		}
		err = connB.WriteJSON(msgFromB)
		if err != nil {
			t.Errorf("failed to send msg: %v", err)
		}
		msgRespFromA := &models.Message{}
		err = connB.ReadJSON(msgRespFromA)
		if err != nil {
			t.Error(err)
		}
		if msgRespFromA.Payload != textFromA {
			t.Errorf("the received message is wrong")
		}

		err = connB.ReadJSON(msgRespFromA)
		if err != nil {
			t.Error(err)
		}
		if msgRespFromA.Payload != defaultText {
			t.Errorf("the received message is wrong")
		}

		wg.Done()
	}()

	wg.Wait()
}
