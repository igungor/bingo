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

	"github.com/igungor/quackle"
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
		opts.p1type = quackle.PlayerComputerPlayerType
		opts.p1name = "8086"
	} else {
		opts.p1type = quackle.PlayerHumanPlayerType
		opts.p1name = *p1
	}
	if *p2 == "" {
		opts.p2type = quackle.PlayerComputerPlayerType
		opts.p2name = "hal"
	} else {
		opts.p2type = quackle.PlayerHumanPlayerType
		opts.p2name = *p2
	}
	game := newGame(opts)
	game.draw()
	game.loop()
}

func usage() {
	fmt.Fprintln(os.Stderr, "bingo is a crossword game.")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Rules:")
	fmt.Fprintln(os.Stderr, "  * input box accepts quackle-format")
	fmt.Fprintln(os.Stderr, "    - `H2 NABER` means start from `H2` and place the word `NABER` from top to bottom.")
	fmt.Fprintln(os.Stderr, "    - `3B HELO` means start from `3B` and place the word `HELO` from left to right.")
	fmt.Fprintln(os.Stderr, "    - words must be typed in all uppercase, except jokers.")
	fmt.Fprintln(os.Stderr, "    - jokers must be typed in lowercase. `n` is the joker letter for move `4F CAMEKAn`.")
	fmt.Fprintln(os.Stderr, "    - type `-` to pass your turn.")
	fmt.Fprintln(os.Stderr, "    - use `.` if a letter of the word you type is already on the board.")
	fmt.Fprintln(os.Stderr, "  * `ctrl-t` toggles multipliers and scores.")
	fmt.Fprintln(os.Stderr, "  * `ctrl-f` fill the input box with a highscore move for you.")
	fmt.Fprintln(os.Stderr, "  * `ctrl-s` shuffles your rack.")
	fmt.Fprintln(os.Stderr, "  * `ctrl-l` toggles legend.")
	fmt.Fprintln(os.Stderr, "  * `ctrl-c` or `esc` quit the game.")
}
