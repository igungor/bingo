package main

import (
	"github.com/igungor/quackle"
	termbox "github.com/nsf/termbox-go"
)

type rack struct {
	qr    quackle.Rack
	title string
	w, h  int
}

func newRack(title string, r quackle.Rack) rack {
	return rack{
		qr:    r,
		title: title,
		w:     7 * 2,
		h:     1,
	}
}

func (r *rack) draw(x, y int) {
	drawRect(x, y, r.w, r.h)
	tbprint("┤"+r.title+"├", x+2, y-1, fgcolor, bgcolor)
	var i int
	for _, r := range r.qr.ToString() {
		termbox.SetCell(x+i*2, y, r, fgcolor|termbox.AttrUnderline|termbox.AttrBold, bgcolor)
		i++
	}
}

func (r *rack) highlight(x, y int) {
	tbprint(r.title, x+2+1, y-1, termbox.ColorWhite, termbox.ColorMagenta)
}

func (r *rack) shuffle() {
	r.qr.Shuffle()
}
