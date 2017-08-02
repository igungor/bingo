package main

// TODO(ig): show score for each letters in rack
// TODO(ig): help screen (shortcuts etc.)
// TODO(ig): tile exchange support

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
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
	flag.Usage = usage
	flag.Parse()
	if err := termbox.Init(); err != nil {
		log.Fatal(err)
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
}

func usage() {
	fmt.Fprintf(os.Stderr, "bingo is a crossword game.")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "Rules:")
	fmt.Fprintf(os.Stderr, "  * input box accepts quackle-format")
	fmt.Fprintf(os.Stderr, "    - `H2 NABER` means start from `H2` and place the word `NABER` from top to bottom.")
	fmt.Fprintf(os.Stderr, "    - `3B HELO` means start from `3B` and place the word `HELO` from left to right.")
	fmt.Fprintf(os.Stderr, "    - words must be typed in all uppercase, except jokers.")
	fmt.Fprintf(os.Stderr, "    - jokers must be typed in lowercase. `n` is the joker letter for move `4F CAMEKAn`.")
	fmt.Fprintf(os.Stderr, "    - type `-` to pass your turn.")
	fmt.Fprintf(os.Stderr, "    - use `.` if a letter of the word you type is already on the board.")
	fmt.Fprintf(os.Stderr, "  * `ctrl-t` toggles multipliers and scores.")
	fmt.Fprintf(os.Stderr, "  * `ctrl-f` fill the input box with a highscore move for you.")
	fmt.Fprintf(os.Stderr, "  * `ctrl-s` shuffles your rack.")
	fmt.Fprintf(os.Stderr, "  * `ctrl-l` toggles legend.")
	fmt.Fprintf(os.Stderr, "  * `ctrl-c` or `esc` quit the game.")
}
