package main

import (
	"log"
	"math/rand"
	"time"

	termbox "github.com/nsf/termbox-go"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	if err := realMain(); err != nil {
		log.Fatal(err)
	}
}

func realMain() error {
	if err := termbox.Init(); err != nil {
		return err
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	game := newGame()
	game.draw()
	game.loop()
	return nil
}
