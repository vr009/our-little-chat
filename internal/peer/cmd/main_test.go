package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	websocket2 "github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	models2 "our-little-chatik/internal/models"
	"sync"
	"testing"
	"time"
)

var (
	userID1     = uuid.New()
	userID2     = uuid.New()
	userID3     = uuid.New()
	chatID      = uuid.New()
	chatID1     = uuid.New()
	chatID2     = uuid.New()
	textFromA   = "Hi, B!"
	textFromB   = "Hi, A!"
	defaultText = "what's up"
)

func TestPeer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	wg := sync.WaitGroup{}
	wg.Add(2)

	t.Logf("user A: %v", userID1)
	t.Logf("user B: %v", userID2)

	clientA := websocket2.DefaultDialer
	clientB := websocket2.DefaultDialer

	host := os.Getenv("PEER_HOST")
	port := os.Getenv("PEER_PORT")

	u1 := url.URL{
		Scheme: "ws",
		Host:   host + ":" + port,
		Path:   "/ws/chat",
	}

	u2 := url.URL{
		Scheme: "ws",
		Host:   host + ":" + port,
		Path:   "/ws/chat",
	}

	url1 := u1.String() + fmt.Sprintf("?chat_id=%s&user_id=%s", chatID.String(), userID1.String())
	url2 := u2.String() + fmt.Sprintf("?chat_id=%s&user_id=%s", chatID.String(), userID2.String())

	time.Sleep(time.Second)
	log.Printf("connecting to %s", url1)
	connA, _, err := clientA.Dial(url1, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer connA.Close()
	log.Printf("Connected")
	_, bMsg, err := connA.ReadMessage()
	if err != nil {
		t.Error(err.Error())
	}
	if string(bMsg) != fmt.Sprintf("Welcome %s!", userID1.String()) {
		t.Errorf("wrong welcome msg")
	}

	log.Printf("connecting to %s", url2)
	connB, _, err := clientB.Dial(url2, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer connB.Close()
	_, bMsg, err = connB.ReadMessage()
	if err != nil {
		t.Error(err.Error())
	}
	if string(bMsg) != fmt.Sprintf("Welcome %s!", userID2.String()) {
		t.Errorf("wrong welcome msg")
	}
	log.Printf("Connected")

	msgFromA := models2.Message{
		ChatID:   chatID,
		SenderID: userID1,
		Payload:  textFromA,
	}
	msgFromA2 := models2.Message{
		ChatID:   chatID,
		SenderID: userID1,
		Payload:  defaultText,
	}

	msgFromB := models2.Message{
		ChatID:   chatID,
		SenderID: userID2,
		Payload:  textFromB,
	}

	go func() {
		err = connA.WriteMessage(1, []byte(msgFromA.Payload))
		if err != nil {
			t.Errorf("failed to send msg: %v", err)
		}
		t.Logf("sent %s", msgFromA.Payload)
		err = connA.WriteMessage(1, []byte(msgFromA2.Payload))
		if err != nil {
			t.Errorf("failed to send msg: %v", err)
		}
		t.Logf("sent %s", msgFromA2.Payload)

		msgRespFromB := &models2.Message{}
		for _ = range []int{1, 2, 3} {
			_, bData, err := connA.ReadMessage()
			t.Logf("%s", string(bData))
			if err != nil {
				t.Error(err)
			}
			err = json.Unmarshal(bData, &msgRespFromB)
			if err != nil {
				t.Error(err)
			}
			if msgRespFromB.Payload != textFromB &&
				msgRespFromB.Payload != msgFromA.Payload &&
				msgRespFromB.Payload != msgFromA2.Payload {
				t.Errorf("the received message is wrong")
			}
			t.Logf("client a finished")
		}
		wg.Done()
	}()

	go func() {
		err = connB.WriteMessage(1, []byte(msgFromB.Payload))
		if err != nil {
			t.Errorf("failed to send msg: %v", err)
		}
		t.Logf("sent %s", msgFromB.Payload)
		msgRespFromA := &models2.Message{}
		_, bData, err := connB.ReadMessage()
		if err != nil {
			t.Error(err)
		}
		err = json.Unmarshal(bData, &msgRespFromA)
		if err != nil {
			t.Error(err)
		}
		t.Logf("1 client b received %s", string(bData))

		if msgRespFromA.Payload != textFromA &&
			msgRespFromA.Payload != defaultText &&
			msgRespFromA.Payload != msgFromB.Payload {
			t.Errorf("the received message is wrong")
		}

		t.Logf("client b expects")
		_, bData, err = connB.ReadMessage()
		if err != nil {
			t.Error(err)
		}
		err = json.Unmarshal(bData, &msgRespFromA)
		if err != nil {
			t.Error(err)
		}
		t.Logf("2 client b received %s", string(bData))
		if msgRespFromA.Payload != textFromA && msgRespFromA.Payload != defaultText {
			t.Errorf("the received message is wrong")
		}
		t.Logf("client b finished")
		wg.Done()
	}()

	wg.Wait()
}

func TestDiff(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	t.Logf("user A: %v", userID1)
	t.Logf("user B: %v", userID2)

	clientA := websocket2.DefaultDialer
	clientB := websocket2.DefaultDialer

	host := os.Getenv("PEER_HOST")
	port := os.Getenv("PEER_PORT")

	u1 := url.URL{
		Scheme: "ws",
		Host:   host + ":" + port,
		Path:   "/ws/chat",
	}

	u2 := url.URL{
		Scheme: "ws",
		Host:   host + ":" + port,
		Path:   "/ws/chat",
	}

	u3 := url.URL{
		Scheme: "ws",
		Host:   host + ":" + port,
		Path:   "/ws/diff",
	}

	url1 := u1.String() + fmt.Sprintf("?chat_id=%s&user_id=%s", chatID1.String(), userID1.String())
	url2 := u2.String() + fmt.Sprintf("?chat_id=%s&user_id=%s", chatID2.String(), userID2.String())
	url3 := u3.String() + fmt.Sprintf("?user_id=%s", userID3.String())

	time.Sleep(time.Second)
	log.Printf("connecting to %s", url1)
	connA, _, err := clientA.Dial(url1, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Connected")
	_, bMsg, err := connA.ReadMessage()
	if err != nil {
		t.Error(err.Error())
	}
	if string(bMsg) != fmt.Sprintf("Welcome %s!", userID1.String()) {
		t.Errorf("wrong welcome msg: %s", string(bMsg))
	}

	log.Printf("connecting to %s", url2)
	connB, _, err := clientB.Dial(url2, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, bMsg, err = connB.ReadMessage()
	if err != nil {
		t.Error(err.Error())
	}
	if string(bMsg) != fmt.Sprintf("Welcome %s!", userID2.String()) {
		t.Errorf("wrong welcome msg: %s", string(bMsg))
	}
	log.Printf("Connected")

	connC, _, err := clientB.Dial(url3, nil)
	if err != nil {
		t.Fatal(err)
	}

	msgFromA := models2.Message{
		ChatID:   chatID1,
		SenderID: userID1,
		Payload:  textFromA,
	}

	msgFromB := models2.Message{
		ChatID:   chatID2,
		SenderID: userID2,
		Payload:  textFromB,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		chats := []models2.Chat{
			{
				ChatID: chatID1,
			},
			{
				ChatID: chatID2,
			},
		}
		bMsg, err := json.Marshal(chats)
		if err != nil {
			t.Errorf(err.Error())
		}

		err = connC.WriteMessage(1, bMsg)
		if err != nil {
			t.Error(err)
		}

		wg.Done()
		_, bData, err := connC.ReadMessage()
		if err != nil {
			t.Error(err)
		}
		t.Logf("got message 1")
		msg := models2.Message{}
		err = json.Unmarshal(bData, &msg)
		if err != nil {
			t.Error(err)
		}
		if msg.ChatID != chatID1 && msg.ChatID != chatID2 {
			t.Error("chat id is unknown")
		}

		_, bData, err = connC.ReadMessage()
		if err != nil {
			t.Error(err)
		}
		t.Logf("got message 2")
		msg2 := models2.Message{}
		err = json.Unmarshal(bData, &msg2)
		if err != nil {
			t.Error(err)
		}
		if msg2.ChatID != chatID1 && msg2.ChatID != chatID2 {
			t.Error("chat id is unknown")
		}
	}()

	wg.Wait()
	wg.Add(1)
	go func() {
		err = connA.WriteMessage(1, []byte(msgFromA.Payload))
		if err != nil {
			t.Errorf("failed to send msg: %v", err)
		}
		t.Logf("sent %s", msgFromA.Payload)

		err = connB.WriteMessage(1, []byte(msgFromB.Payload))
		if err != nil {
			t.Errorf("failed to send msg: %v", err)
		}
		wg.Done()
	}()

	wg.Wait()
}
