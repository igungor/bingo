package main

// TODO(ig): show score for each letters in rack
// TODO(ig): help screen (shortcuts etc.)
// TODO(ig): tile exchange support

import (
	"flag"
	"log"
	"math/rand"
	"time"

	termbox "github.com/nsf/termbox-go"
)

// flags
var (
	p1 = flag.String("p1", "iby", "player1 name. cpu is chosen if empty string is given")
	p2 = flag.String("p2", "", "player2 name. cpu is chosen if empty string is given")
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
	flag.Parse()
	if err := termbox.Init(); err != nil {
		return err
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	opts := &gameOpts{}
	if *p1 == "" {
		opts.p1type = computer
		opts.p1name = "8086"
	} else {
		opts.p1type = human
		opts.p1name = *p1
	}
	if *p2 == "" {
		opts.p2type = computer
		opts.p2name = "hal"
	} else {
		opts.p2type = human
		opts.p2name = *p2
	}
	game := newGame(opts)
	game.draw()
	game.loop()
	return nil
}
