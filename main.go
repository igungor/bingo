package main

// TODO(ig): show score for each letters in rack
// TODO(ig): help screen (shortcuts etc.)
// TODO(ig): highlight error causes on move validity

import (
	"log"
	"math/rand"
	"time"

	termbox "github.com/igungor/termbox-go"
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
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	game := newGame(nil)
	game.draw()
	game.loop()
	return nil
}
