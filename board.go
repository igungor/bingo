package main

import (
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/igungor/quackle"
	termbox "github.com/nsf/termbox-go"
)

const boardsize = 15

type board struct {
	qb        quackle.Board
	w, h      int
	showScore bool
}

func (b *board) draw(x, y int) {
	// columns on the top
	for dx := 0; dx < b.w; dx++ {
		termbox.SetCell(x+dx*2, y-2, rune('A'+dx), fgcolor, bgcolor)
		termbox.SetCell(x+dx*2+1, y-2, ' ', fgcolor, bgcolor)
	}

	// rows on the left
	for dy := 0; dy < b.h; dy++ {
		if dy+1 < 10 {
			tbprint(strconv.Itoa(dy+1), x-2, y+dy, fgcolor, bgcolor)
		} else {
			tbprint(strconv.Itoa(dy+1), x-3, y+dy, fgcolor, bgcolor)
		}
	}

	// borders
	drawRect(x, y, b.w*2, b.h)

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
				termbox.SetCell(x+col*2, y+row, r, fgcolor, bgcolor)
				if b.showScore {
					termbox.SetCell(x+col*2+1, y+row, getScoreRune(score), fgcolor, bgcolor)
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
				tbprint(altch, x+col*2, y+row, color, bgcolor)
			} else {
				tbprint(ch, x+col*2, y+row, color, bgcolor)
			}
		}
	}
}

type legend struct {
}

func (l *legend) draw(x, y int) {
	tbprint("H²", x+0, y, termbox.ColorWhite, termbox.ColorBlue)
	tbprint("H³", x+2, y, termbox.ColorWhite, termbox.ColorMagenta)
	tbprint("K²", x+4, y, termbox.ColorWhite, termbox.ColorGreen)
	tbprint("K³", x+6, y, termbox.ColorWhite, termbox.ColorBlack)
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
