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
	w, h          int
}

func newEditbox() editbox {
	return editbox{
		w: editboxWidth,
		h: 1,
	}
}

// Draws the editbox in the given location, 'h' is not used at the moment
func (eb *editbox) draw(x, y int) {
	drawRect(x, y, eb.w, eb.h)

	t := eb.text
	lx := 0
	for {
		if len(t) == 0 {
			break
		}
		r, size := utf8.DecodeRune(t)
		// topdown words
		if lx == 0 && unicode.IsLetter(r) {
			termbox.SetCell(x, y, '↓', fgcolor|termbox.AttrBold, bgcolor)
		}
		// leftright words
		if lx == 0 && unicode.IsNumber(r) {
			termbox.SetCell(x, y, '→', fgcolor|termbox.AttrBold, bgcolor)
		}

		// pass
		if lx == 0 && r == '-' {
			termbox.SetCell(x, y, '⚐', fgcolor|termbox.AttrBold, bgcolor)
		}

		termbox.SetCell(x+lx+2, y, r, fgcolor, bgcolor)
		lx += 1
		t = t[size:]
	}
	termbox.SetCursor(x+lx+2, y)
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
	eb.MoveCursorTo(0)
	eb.text = eb.text[:eb.curByteOffset]
}

func (eb *editbox) MoveCursorTo(boffset int) {
	eb.curByteOffset = boffset
}

func (eb *editbox) RuneUnderCursor() (rune, int) {
	return utf8.DecodeRune(eb.text[eb.curByteOffset:])
}

func (eb *editbox) RuneBeforeCursor() (rune, int) {
	return utf8.DecodeLastRune(eb.text[:eb.curByteOffset])
}

func (eb *editbox) MoveCursorOneRuneBackward() {
	if eb.curByteOffset == 0 {
		return
	}
	_, size := eb.RuneBeforeCursor()
	eb.MoveCursorTo(eb.curByteOffset - size)
}

func (eb *editbox) MoveCursorOneRuneForward() {
	if eb.curByteOffset == len(eb.text) {
		return
	}
	_, size := eb.RuneUnderCursor()
	eb.MoveCursorTo(eb.curByteOffset + size)
}

func (eb *editbox) DeleteRuneBackward() {
	if eb.curByteOffset == 0 {
		return
	}

	eb.MoveCursorOneRuneBackward()
	_, size := eb.RuneUnderCursor()
	eb.text = byteSliceRemove(eb.text, eb.curByteOffset, eb.curByteOffset+size)
}

func (eb *editbox) InsertRune(r rune) {
	// ___ _______ -> total of 11 runes
	if utf8.RuneCount(eb.text) >= 11 {
		return
	}
	r = unicode.TurkishCase.ToUpper(r)
	var buf [utf8.UTFMax]byte
	n := utf8.EncodeRune(buf[:], r)
	eb.text = byteSliceInsert(eb.text, eb.curByteOffset, buf[:n])
	eb.MoveCursorOneRuneForward()
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
