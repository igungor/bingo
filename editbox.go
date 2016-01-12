package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	termbox "github.com/nsf/termbox-go"
)

const preferredHorizontalThreshold = 3
const editboxWidth = 14

type editbox struct {
	text          []byte
	curByteOffset int // cursor offset in bytes
	x, y          int
	w, h          int
	warn, err     bool
}

func newEditbox() editbox {
	return editbox{
		w: editboxWidth,
		h: 1,
	}
}

// Draws the editbox at (x,y) position.
func (eb *editbox) draw() {
	drawRect(eb.x, eb.y, eb.w, eb.h)

	t := eb.text
	lx := 0
	for {
		if len(t) == 0 {
			break
		}
		r, size := utf8.DecodeRune(t)
		// topdown words
		if lx == 0 && unicode.IsLetter(r) {
			termbox.SetCell(eb.x, eb.y, '↓', fgcolor|termbox.AttrBold, bgcolor)
		}
		// leftright words
		if lx == 0 && unicode.IsNumber(r) {
			termbox.SetCell(eb.x, eb.y, '→', fgcolor|termbox.AttrBold, bgcolor)
		}
		// pass
		if lx == 0 && r == '-' {
			termbox.SetCell(eb.x, eb.y, '⚐', fgcolor|termbox.AttrBold, bgcolor)
		}
		termbox.SetCell(eb.x+lx+2, eb.y, r, fgcolor, bgcolor)
		lx += 1
		t = t[size:]
	}
	termbox.SetCursor(eb.x+lx+2, eb.y)

	// error/warning indicator
	if eb.warn {
		termbox.SetCell(eb.x, eb.y, '❗', fgcolor|termbox.AttrBold, bgcolor)
	} else if eb.err {
		termbox.SetCell(eb.x, eb.y, '⊗', fgcolor|termbox.AttrBold, bgcolor)
	}
}

func (eb *editbox) getPlaceWord() (string, string, error) {
	s := strings.Fields(string(eb.text))
	if len(s) == 1 && s[0] == "-" {
		return s[0], "", nil
	}
	if len(s) != 2 {
		return "", "", fmt.Errorf("`B2 SELAM` formatinda olmali")
	}
	return s[0], s[1], nil
}

func (eb *editbox) clear() {
	eb.moveCursorTo(0)
	eb.text = eb.text[:eb.curByteOffset]
	eb.warn = false
	eb.err = false
}

func (eb *editbox) moveCursorTo(boffset int) {
	eb.curByteOffset = boffset
}

func (eb *editbox) runeUnderCursor() (rune, int) {
	return utf8.DecodeRune(eb.text[eb.curByteOffset:])
}

func (eb *editbox) runeBeforeCursor() (rune, int) {
	return utf8.DecodeLastRune(eb.text[:eb.curByteOffset])
}

func (eb *editbox) moveCursorOneRuneBackward() {
	if eb.curByteOffset == 0 {
		return
	}
	_, size := eb.runeBeforeCursor()
	eb.moveCursorTo(eb.curByteOffset - size)
}

func (eb *editbox) moveCursorOneRuneForward() {
	if eb.curByteOffset == len(eb.text) {
		return
	}
	_, size := eb.runeUnderCursor()
	eb.moveCursorTo(eb.curByteOffset + size)
}

func (eb *editbox) deleteRuneBackward() {
	if eb.curByteOffset == 0 {
		return
	}

	eb.moveCursorOneRuneBackward()
	_, size := eb.runeUnderCursor()
	eb.text = byteSliceRemove(eb.text, eb.curByteOffset, eb.curByteOffset+size)
}

func (eb *editbox) insertRune(r rune) {
	// ___ _______ -> total of 11 runes
	if utf8.RuneCount(eb.text) >= 11 {
		return
	}
	var buf [utf8.UTFMax]byte
	n := utf8.EncodeRune(buf[:], r)
	eb.text = byteSliceInsert(eb.text, eb.curByteOffset, buf[:n])
	eb.moveCursorOneRuneForward()
}

func byteSliceGrow(s []byte, desired_cap int) []byte {
	if cap(s) < desired_cap {
		ns := make([]byte, len(s), desired_cap)
		copy(ns, s)
		return ns
	}
	return s
}

func byteSliceRemove(text []byte, from, to int) []byte {
	size := to - from
	copy(text[from:], text[to:])
	text = text[:len(text)-size]
	return text
}

func byteSliceInsert(text []byte, offset int, what []byte) []byte {
	n := len(text) + len(what)
	text = byteSliceGrow(text, n)
	text = text[:n]
	copy(text[offset+len(what):], text[offset:])
	copy(text[offset:], what)
	return text
}
