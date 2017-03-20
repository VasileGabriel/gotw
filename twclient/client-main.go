package main

import (
	// "bufio"
	// "bytes"
	// "image/png"
	// "io"
	"encoding/gob"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"net"
	"os"

	"github.com/VasileGabriel/gotw/twcommon"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

var img image.Image
var decoder *gob.Decoder
var msg twcommon.Message

var (
	blue0    = color.RGBA{0x00, 0x00, 0x1f, 0xff}
	blue1    = color.RGBA{0x00, 0x00, 0x3f, 0xff}
	darkGray = color.RGBA{0x3f, 0x3f, 0x3f, 0xff}
	red      = color.RGBA{0x7f, 0x00, 0x00, 0x7f}
	yellow   = color.RGBA{0x3f, 0x3f, 0x00, 0x3f}

	cos30 = math.Cos(math.Pi / 6)
	sin30 = math.Sin(math.Pi / 6)
)

func main() {

	service := "127.0.0.1:1200"

	conn, err := net.Dial("tcp", service)
	checkError(err)

	// encoder := gob.NewEncoder(conn)
	decoder = gob.NewDecoder(conn)

	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(nil)
		if err != nil {
			log.Fatal(err)
		}
		defer w.Release()

		winSize := image.Point{1000, 1000}
		b, err := s.NewBuffer(winSize)
		if err != nil {
			log.Fatal(err)
		}
		defer b.Release()

		t, err := s.NewTexture(winSize)
		if err != nil {
			log.Fatal(err)
		}
		defer t.Release()
		t.Upload(image.Point{}, b, b.Bounds())

		var sz size.Event
		for {
			e := w.NextEvent()

			// This print message is to help programmers learn what events this
			// example program generates. A real program shouldn't print such
			// messages; they're not important to end users.
			format := "got %#v\n"
			if _, ok := e.(fmt.Stringer); ok {
				format = "got %v\n"
			}
			fmt.Printf(format, e)

			switch e := e.(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}
			case paint.Event:
				const inset = 10
				for _, r := range imageutil.Border(sz.Bounds(), inset) {
					w.Fill(r, blue0, screen.Src)
				}
				w.Fill(sz.Bounds().Inset(inset), blue1, screen.Src)
				w.Upload(image.Point{120, 0}, b, b.Bounds())
				// var i uint8
				// i = 0
				go func() {
					for {

						decoder.Decode(&msg)
						img, err = msg.GetJpegImage()
						checkError(err)

						draw.Draw(b.RGBA(), b.RGBA().Bounds(), img, image.ZP, draw.Src)

						w.Upload(image.Point{0, 0}, b, b.Bounds())

						// blue := color.RGBA{0, 0, i, 255}
						// draw.Draw(b.RGBA(), b.RGBA().Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)
						// i++
						// w.Upload(image.Point{120, 0}, b, b.Bounds())
						// time.Sleep(time.Millisecond * 1000 / 24)

					}
				}()
				// w.Publish()

			case size.Event:
				sz = e

			case error:
				log.Print(e)
			}
		}
	})

	// for {
	// 	decoder.Decode(&msg)
	// 	img, err = msg.GetJpegImage()
	// 	checkError(err)
	// 	// toimg, _ := os.Create("new.jpg")
	// 	// defer toimg.Close()

	// 	// jpeg.Encode(toimg, img, &jpeg.Options{jpeg.DefaultQuality})
	// }

	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func initiateDrawSurface(s screen.Screen) (screen.Window, screen.Buffer) {
	w, err := s.NewWindow(nil)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Release()

	winSize := image.Point{1000, 1000}
	b, err := s.NewBuffer(winSize)
	if err != nil {
		log.Fatal(err)
	}
	defer b.Release()

	t, err := s.NewTexture(winSize)
	if err != nil {
		log.Fatal(err)
	}
	defer t.Release()
	t.Upload(image.Point{}, b, b.Bounds())
	return w, b
}

func drawOnSurface(decoder *gob.Decoder, w screen.Window, b screen.Buffer) {
	var err error
	for {
		decoder.Decode(&msg)
		img, err = msg.GetJpegImage()
		checkError(err)
		draw.Draw(b.RGBA(), b.RGBA().Bounds(), img, image.ZP, draw.Src)

		w.Upload(image.Point{0, 0}, b, b.Bounds())
		// toimg, _ := os.Create("new.jpg")
		// defer toimg.Close()

		// jpeg.Encode(toimg, img, &jpeg.Options{jpeg.DefaultQuality})
	}
}
