package main

import (
	"fmt"
	"os"
	"time"

	"github.com/igungor/quackle"
	termbox "github.com/nsf/termbox-go"
)

const (
	lexicon  = "turkish"
	alphabet = "turkish"
)

// XXX: hacky stuff.
var datadir = fmt.Sprintf("%v/src/github.com/igungor/bingo/data", os.Getenv("GOPATH"))

const (
	fgcolor = termbox.ColorDefault
	bgcolor = termbox.ColorDefault
)

var dm quackle.DataManager
var flexAbc quackle.FlexibleAlphabetParameters

type game struct {
	qg      quackle.Game
	board   board
	rack1   rack
	rack2   rack
	legend  legend
	editbox editbox

	isOver     bool
	showLegend bool
}

func (g *game) draw() {
	// update board and racks
	g.board.qb = g.pos().Board()
	g.rack1.qr = g.player(0).Rack()
	g.rack2.qr = g.player(1).Rack()

	termbox.Clear(fgcolor, bgcolor)
	defer termbox.Flush()

	sw, sh := termbox.Size()

	// board
	boardx := (sw - g.board.w*2 + 2 + 1) / 2
	boardy := (sh - g.board.h + 1 + 1) / 2
	g.board.draw(boardx, boardy)

	// legend
	if g.showLegend {
		legendx := (sw+g.board.w)/2 + 1
		legendy := (sh-g.board.h)/2 + 1 + 1 + g.board.h
		g.legend.draw(legendx, legendy)
	}

	// racks
	rack1x := boardx
	rack1y := boardy + g.board.h + 2
	g.rack1.draw(rack1x, rack1y)
	rack2x := boardx + g.rack1.w + 2
	rack2y := boardy + g.board.h + 2
	g.rack2.draw(rack2x, rack2y)
	if g.curPlayer().Id() == 0 {
		g.rack1.highlight(rack1x, rack1y)
	} else {
		g.rack2.highlight(rack2x, rack2y)
	}

	// editbox
	boxx := (sw-g.editbox.w)/2 + 1
	boxy := (sh+g.board.h)/2 + g.rack1.h + 5
	g.editbox.draw(boxx, boxy)
}

func (g *game) loop() {
mainloop:
	for {
		if g.pos().GameOver() {
			g.isOver = true
			g.over()
			break mainloop
		}

		if g.curPlayer().Xtype() == 0 {
			g.qg.HaveComputerPlay()
			g.draw()
		}
		g.qg.AdvanceToNoncomputerPlayer()

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEnter:
				// new human move
				g.doHumanMove()
			case termbox.KeyCtrlS:
				g.board.showScore = !g.board.showScore
			case termbox.KeyCtrlL:
				g.showLegend = !g.showLegend
			case termbox.KeyCtrlT:
				g.showHint()
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				g.editbox.DeleteRuneBackward()
			case termbox.KeySpace:
				g.editbox.InsertRune(' ')
			case termbox.KeyEsc, termbox.KeyCtrlC:
				break mainloop
			default:
				if ev.Ch != 0 {
					g.editbox.InsertRune(ev.Ch)
				}
			}
		case termbox.EventResize:
			g.draw()
		case termbox.EventMouse:
			// TODO(ig): handle mouse clicks
		case termbox.EventError:
			panic(ev.Err)
		}
		g.draw()
	}
}

// FIXME(ig): print proper errors
func (g *game) doHumanMove() {
	place, word, err := g.editbox.getPlaceWord()
	if err != nil {
		return
	}
	var move quackle.Move
	if place == "-" {
		move = quackle.MoveCreatePassMove()
	} else {
		move = quackle.MoveCreatePlaceMove(place, flexAbc.Encode(word))
	}
	if g.pos().ValidateMove(move) != 0 {
		return
	}
	g.qg.CommitMove(move)
	g.editbox.clear()
}

func (g *game) showHint() {
	g.editbox.clear()
	// accept top scoreing advice from the beast
	g.pos().Kibitz(1)
	if g.pos().Moves().Size() == 0 {
		return
	}
	move := g.pos().Moves().SwigGetMoveVector().Get(0)
	for _, r := range move.ToString() {
		g.editbox.InsertRune(r)
	}
	termbox.Flush()
}

// over draws game-over screen.
func (g *game) over() {
	termbox.Clear(fgcolor, bgcolor)
	sw, sh := termbox.Size()
	tbprint("GAME OVER", sw/2-4, sh/2, fgcolor, bgcolor)
	termbox.Flush()
	time.Sleep(1 * time.Second)
}

// pos returns current game position
func (g *game) pos() quackle.GamePosition {
	return g.qg.CurrentPosition().(quackle.GamePosition)
}

// player returns current player
func (g *game) curPlayer() quackle.Player {
	return g.pos().CurrentPlayer().(quackle.Player)
}

func (g *game) player(id int) quackle.Player {
	found := make([]bool, 1)
	return g.pos().Players().PlayerForId(id, found)
}

// newGame initializes a new game and constructs game object.
func newGame() *game {
	dm = quackle.NewDataManager()
	dm.SetComputerPlayers(quackle.ComputerPlayerCollectionFullCollection().SwigGetPlayerList())
	dm.SetBackupLexicon(lexicon)
	dm.SetAppDataDirectory(datadir)

	// set up alphabet
	abc := quackle.AlphabetParametersFindAlphabetFile(alphabet)
	qabc := quackle.UtilStdStringToQString(abc)
	flexAbc = quackle.NewFlexibleAlphabetParameters()
	flexAbc.Load(qabc)
	dm.SetAlphabetParameters(flexAbc)

	// set up board
	bp := quackle.NewBoardParameters()
	for y := 0; y < boardsize; y++ {
		for x := 0; x < boardsize; x++ {
			bp.SetLetterMultiplier(x, y, quackle.QuackleBoardParametersLetterMultiplier(boardLetterMult[x][y]))
			bp.SetWordMultiplier(x, y, quackle.QuackleBoardParametersWordMultiplier(boardWordMult[x][y]))
		}
	}
	dm.SetBoardParameters(bp)

	// find lexicon
	dawg := quackle.LexiconParametersFindDictionaryFile(lexicon + ".dawg")
	gaddag := quackle.LexiconParametersFindDictionaryFile(lexicon + ".gaddag")
	dm.LexiconParameters().LoadDawg(dawg)
	dm.LexiconParameters().LoadGaddag(gaddag)
	dm.StrategyParameters().Initialize(lexicon)

	dm.SeedRandomNumbers(uint(time.Now().UnixNano()))

	newCompPlayer := func(name string, id int) quackle.Player {
		found := make([]bool, 1)
		player := dm.ComputerPlayers().PlayerForName("Speedy Player", found)
		if !found[0] {
			panic("player could not be found")
		}
		comp := player.ComputerPlayer()

		p := quackle.NewPlayer(name, int(quackle.PlayerComputerPlayerType), id)
		p.SetComputerPlayer(comp)
		return p
	}

	// set up players and game
	g := quackle.NewGame()
	player1 := quackle.NewPlayer("iby", int(quackle.PlayerHumanPlayerType), 0)
	player2 := newCompPlayer("Computer", 1)
	players := quackle.NewPlayerList()
	players.Add(player1)
	players.Add(player2)
	g.SetPlayers(players)
	g.AssociateKnownComputerPlayers()
	g.AddPosition()

	b := board{
		qb: g.CurrentPosition().(quackle.GamePosition).Board(),
		w:  boardsize,
		h:  boardsize,
	}

	return &game{
		qg:         g,
		board:      b,
		rack1:      newRack(player1.Name(), player1.Rack()),
		rack2:      newRack(player2.Name(), player2.Rack()),
		editbox:    newEditbox(),
		showLegend: true,
	}
}
