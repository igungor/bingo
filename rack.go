package main

import (
	"strconv"

	"github.com/igungor/quackle"
	termbox "github.com/nsf/termbox-go"
)

type rack struct {
	player quackle.Player
	w, h   int
}

func newRack(pl quackle.Player) rack {
	return rack{
		player: pl,
		w:      7 * 2,
		h:      1,
	}
}

func (r *rack) draw(x, y int) {
	drawRect(x, y, r.w, r.h)
	tbprint("┤"+r.player.Name()+"├", x, y-1, fgcolor, bgcolor)

	// BUG: player1 score is not adding up. wtf is going on?
	// player score
	playerScore := strconv.Itoa(r.player.Score())
	tbprint(playerScore, x+r.w-len(playerScore), y-1, fgcolor, bgcolor)

	// i dont use tbprint helper function because letters have spacing between them
	var i int
	for _, r := range r.player.Rack().ToString() {
		termbox.SetCell(x+i*2, y, r, fgcolor|termbox.AttrUnderline|termbox.AttrBold, bgcolor)
		i++
	}
}

func (r *rack) highlight(x, y int) {
	tbprint(r.player.Name(), x+1, y-1, termbox.ColorWhite, termbox.ColorMagenta)
}
