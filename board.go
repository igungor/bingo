package main

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/igungor/quackle"
	termbox "github.com/igungor/termbox-go"
)

const boardsize = 15

type board struct {
	qb         quackle.Board
	w, h       int
	x, y       int
	showScore  bool
	curPosVert bool
}

func (b *board) draw() {
	// columns on the top
	for dx := 0; dx < b.w; dx++ {
		termbox.SetCell(b.x+dx*2, b.y-2, rune('A'+dx), fgcolor, bgcolor)
		termbox.SetCell(b.x+dx*2+1, b.y-2, ' ', fgcolor, bgcolor)
	}

	// rows on the left
	for dy := 0; dy < b.h; dy++ {
		if dy+1 < 10 {
			tbprint(strconv.Itoa(dy+1), b.x-2, b.y+dy, fgcolor, bgcolor)
		} else {
			tbprint(strconv.Itoa(dy+1), b.x-3, b.y+dy, fgcolor, bgcolor)
		}
	}

	// borders
	drawRect(b.x, b.y, b.w*2, b.h)

	// multipliers and letters
	for row := 0; row < b.h; row++ {
		for col := 0; col < b.w; col++ {
			// mark letters
			bl := b.qb.Letter(row, col)
			if bl != byte(quackle.QUACKLE_NULL_MARK) {
				letter := flexAbc.UserVisible(bl)

				r, _ := utf8.DecodeRuneInString(letter)
				// BUG(ig): score lookups for joker letters are broken and crashes the app.
				// FIXME(ig): fix this ugly hack!
				var score int
				if unicode.IsUpper(r) {
					score = flexAbc.Score(bl)
				}
				termbox.SetCell(b.x+col*2, b.y+row, r, fgcolor, bgcolor)
				if b.showScore {
					termbox.SetCell(b.x+col*2+1, b.y+row, getScoreRune(score), fgcolor, bgcolor)
				}
				continue
			}

			// mark multipliers
			letterMult := dm.BoardParameters().LetterMultiplier(row, col)
			wordMult := dm.BoardParameters().WordMultiplier(row, col)
			multChar := "★"
			var ch string
			var altch string
			color := fgcolor
			switch {
			case letterMult == 2:
				color = termbox.ColorBlue
				altch = "h²"
				ch = multChar
			case letterMult == 3:
				color = termbox.ColorMagenta
				altch = "h³"
				ch = multChar
			case wordMult == 2:
				color = termbox.ColorGreen
				altch = "k²"
				ch = multChar
			case wordMult == 3:
				color = termbox.ColorBlack
				altch = "k³"
				ch = multChar
			default:
				ch = " "
				altch = " "
			}
			if b.showScore {
				tbprint(altch, b.x+col*2, b.y+row, color, bgcolor)
			} else {
				tbprint(ch, b.x+col*2, b.y+row, color, bgcolor)
			}
		}
	}
}

// in reports whether the given (x,y) coordinates overlaps with the board.
func (b *board) in(x, y int) bool {
	if x >= b.x && x < b.x+b.w*2 &&
		y >= b.y && y < b.y+b.h {
		return true
	}
	return false
}

func (b *board) highlightPos(x, y int) {
	if !b.in(x, y) {
		return
	}

	b.curPosVert = !b.curPosVert

	fg := fgcolor | termbox.AttrBold
	bg := bgcolor | termbox.AttrBold

	defch := '·'
	var xch, ych rune
	if b.curPosVert {
		xch = defch
		ych = '⇣'
	} else {
		xch = '⇢'
		ych = defch
	}

	// draw position indicators
	for i := b.x; i < x; i += 2 {
		termbox.SetCell(i, y, xch, fg, bg)
	}
	for i := b.y; i <= y; i++ {
		if (x-b.x)%2 == 1 {
			termbox.SetCell(x-1, i, ych, fg, bg)
		} else {
			termbox.SetCell(x, i, ych, fg, bg)
		}
	}
	termbox.Flush()
}

func (b *board) pos(x, y int) string {
	if !b.in(x, y) {
		return ""
	}

	xpos := (x - b.x) / 2
	ypos := y - b.y
	var place string
	if b.curPosVert {
		place = fmt.Sprintf("%v%v", string(rune('A'+xpos)), ypos+1)
	} else {
		place = fmt.Sprintf("%v%v", ypos+1, string(rune('A'+xpos)))
	}
	return place
}

type legend struct {
	x, y int
}

func (l *legend) draw() {
	tbprint("H²", l.x+0, l.y, termbox.ColorWhite, termbox.ColorBlue)
	tbprint("H³", l.x+2, l.y, termbox.ColorWhite, termbox.ColorMagenta)
	tbprint("K²", l.x+4, l.y, termbox.ColorWhite, termbox.ColorGreen)
	tbprint("K³", l.x+6, l.y, termbox.ColorWhite, termbox.ColorBlack)
}

var score2rune = []rune{'₀', '₁', '₂', '₃', '₄', '₅', '₆', '₇', '₈', '₉', '⏨'}

func getScoreRune(score int) (r rune) {
	return score2rune[score]
}

// kelimelik board
var (
	boardLetterMult = [boardsize][boardsize]int{
		{1, 1, 1, 1, 1, 2, 1, 1, 1, 2, 1, 1, 1, 1, 1},
		{1, 3, 1, 1, 1, 1, 2, 1, 2, 1, 1, 1, 1, 3, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 3, 1, 1, 1, 1, 1, 3, 1, 1, 1, 1},
		{2, 1, 1, 1, 1, 2, 1, 1, 1, 2, 1, 1, 1, 1, 2},
		{1, 2, 1, 1, 1, 1, 2, 1, 2, 1, 1, 1, 1, 2, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 2, 1, 1, 1, 1, 2, 1, 2, 1, 1, 1, 1, 2, 1},
		{2, 1, 1, 1, 1, 2, 1, 1, 1, 2, 1, 1, 1, 1, 2},
		{1, 1, 1, 1, 3, 1, 1, 1, 1, 1, 3, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 3, 1, 1, 1, 1, 2, 1, 2, 1, 1, 1, 1, 3, 1},
		{1, 1, 1, 1, 1, 2, 1, 1, 1, 2, 1, 1, 1, 1, 1},
	}
	boardWordMult = [boardsize][boardsize]int{
		{1, 1, 3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{3, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 3},
		{1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 2, 1, 1, 1, 1, 2, 1, 1, 1, 1, 2, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1},
		{3, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 3},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 1},
	}
)
