package twcommon

import (
	"encoding/gob"
	"fmt"
	"net"
)

type Clients map[int]Client

type Client struct {
	Id      string
	Conn    net.Conn
	Encoder *gob.Encoder
	Decoder *gob.Decoder
}

func (c Clients) String() string {
	if len(c) == 0 {
		return "No clients!!"
	}
	var s = fmt.Sprintln(" clients: ", len(c))
	i := 1
	for _, client := range c {
		s += fmt.Sprintln("Client no. ", i, " with id ", client.Id)
		i++
	}
	return s
}
