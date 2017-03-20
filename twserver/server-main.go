package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/VasileGabriel/gotw/twcommon"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	// "github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/nu7hatch/gouuid"
)

func main() {

	clients := twcommon.Clients{}

	service := "0.0.0.0:1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	messagesChannel := make(chan twcommon.Message)
	go queueMessages(messagesChannel)
	go sendMessages(&clients, messagesChannel)

	listen(listener, messagesChannel, clients)
}

func listen(listener *net.TCPListener,
	messagesChannel chan twcommon.Message,
	clients twcommon.Clients) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		encoder := gob.NewEncoder(conn)
		uuid, _ := uuid.NewV4()
		client := twcommon.Client{Id: uuid.String(), Conn: conn, Encoder: encoder, Decoder: nil}
		clients[len(clients)] = client
		log.Println(clients)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func queueMessages(messagesChannel chan twcommon.Message) {
	//TODO stop recording if no client connected
	var ximg *xgraphics.Image
	var err error
	var msg *twcommon.Message

	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	drawable := xproto.Drawable(X.RootWin())

	for {
		inceput := time.Now()
		ximg, err = xgraphics.NewDrawable(X, drawable)
		duration := time.Since(inceput)
		log.Println("luarea imaginii a durat ", duration.Seconds())

		if err != nil {
			log.Fatal(err)
		}
		msg, err = twcommon.NewImageMessage(ximg)
		checkError(err)
		messagesChannel <- *msg
		// time.Sleep(40 * time.Millisecond)
	}
}

func sendMessages(clients *twcommon.Clients, messagesChannel chan twcommon.Message) {
	var msg twcommon.Message
	for {
		msg = <-messagesChannel
		msg.Dispatch(clients)
	}
}
