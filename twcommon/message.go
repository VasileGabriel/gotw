package twcommon

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"image"
	"log"
	"strconv"
	"time"

	"github.com/BurntSushi/xgbutil/xgraphics"
	turboJpeg "github.com/pixiv/go-libjpeg/jpeg"
)

const (
	TypeLogin  = 0
	TypeLogout = 1
	TypeText   = 2
	TypeImage  = 4
)

type Message struct {
	MessageType int
	Content     []byte
}

func (m Message) String() string {
	s := "mesaj de tip " + strconv.Itoa(m.MessageType)
	if len(m.Content) > 10 {
		s += " si ceva continut: " + string(m.Content[:10])
	}
	return s
}
func NewImageMessage(img *xgraphics.Image) (*Message, error) {
	buf := new(bytes.Buffer)
	inceput := time.Now()
	err := turboJpeg.Encode(buf, img, &turboJpeg.EncoderOptions{Quality: 100})
	duration := time.Since(inceput)
	log.Println("encodarea a durat ", duration.Seconds())

	if err != nil {
		return nil, err
	}
	return &Message{TypeImage, buf.Bytes()}, nil

}

func (m Message) GetJpegImage() (image.Image, error) {
	buffer := bytes.NewBuffer(m.Content)
	r := bufio.NewReader(buffer)
	return turboJpeg.Decode(r, &turboJpeg.DecoderOptions{})
}

func (m Message) WriteToEncoder(e *gob.Encoder) error {
	return e.Encode(m)
}

func (m Message) Dispatch(clients *Clients) {
	for index, client := range *clients {
		go func() {
			err := m.WriteToEncoder(client.Encoder)
			if err != nil {
				log.Println(err)
				delete(*clients, index)
			}
		}()
	}
}
