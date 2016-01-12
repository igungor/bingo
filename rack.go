package main

import (
	"strconv"

	"github.com/igungor/quackle"
	termbox "github.com/nsf/termbox-go"
)

type rack struct {
	player quackle.Player
	x, y   int
	w, h   int
}

func newRack(pl quackle.Player) rack {
	return rack{
		player: pl,
		w:      7 * 2,
		h:      1,
	}
}

func (r *rack) draw() {
	drawRect(r.x, r.y, r.w, r.h)
	tbprint("┤"+r.player.Name()+"├", r.x, r.y-1, fgcolor, bgcolor)

	// BUG: player1 score is not adding up. wtf is going on?
	// player score
	playerScore := strconv.Itoa(r.player.Score())
	tbprint(playerScore, r.x+r.w-len(playerScore), r.y-1, fgcolor, bgcolor)

	// tbprint helper function is not used because letters have spacing between them
	var i int
	for _, ch := range r.player.Rack().ToString() {
		termbox.SetCell(r.x+i*2, r.y, ch, fgcolor|termbox.AttrUnderline|termbox.AttrBold, bgcolor)
		i++
	}
}

func (r *rack) highlight() {
	tbprint(r.player.Name(), r.x+1, r.y-1, termbox.ColorWhite, termbox.ColorMagenta)
}
