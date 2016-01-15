package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/igungor/quackle"
	termbox "github.com/igungor/termbox-go"
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
	qg quackle.Game

	// widgets
	board   board
	rack1   rack
	rack2   rack
	legend  legend
	editbox editbox

	// screen width, height
	w, h int

	// toggle switch legend
	showLegend bool

	// gameplay error
	err error
}

func (g *game) draw() {
	// update board and racks
	g.board.qb = g.pos().Board()
	g.rack1.player = g.player(0)
	g.rack2.player = g.player(1)
	g.w, g.h = termbox.Size()

	termbox.Clear(fgcolor, bgcolor)
	defer termbox.Flush()

	// board
	g.board.x = (g.w - g.board.w*2 + 2 + 1) / 2
	g.board.y = (g.h - g.board.h - g.rack1.h - g.editbox.h - 2) / 2
	g.board.draw()

	// racks
	g.rack1.x = g.board.x
	g.rack1.y = g.board.y + g.board.h + 2
	g.rack1.draw()
	g.rack2.x = g.board.x + g.rack1.w + 2
	g.rack2.y = g.board.y + g.board.h + 2
	g.rack2.draw()
	if g.curPlayer().Id() == 0 {
		g.rack1.highlight()
	} else {
		g.rack2.highlight()
	}

	// editbox
	g.editbox.x = g.board.x + g.editbox.w/2 + 1
	g.editbox.y = g.rack1.y - g.rack1.h + 4
	g.editbox.draw()

	// legend
	if g.showLegend {
		g.legend.x = g.board.x + g.board.w*2 - 8
		g.legend.y = g.board.y + g.board.h
		g.legend.draw()
	}
}

func (g *game) loop() {
mainloop:
	for {
		if g.pos().GameOver() {
			g.over()
			break mainloop
		}

		// go on cold heartless cpu.
		if g.curPlayer().Xtype() == int(quackle.PlayerComputerPlayerType) {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
			g.qg.HaveComputerPlay()
			g.draw()
			continue
		}
		// it is our turn, human.
		g.qg.AdvanceToNoncomputerPlayer()

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEnter:
				// new human move
				g.doHumanMove()
			case termbox.KeyCtrlS:
				g.curPlayer().Rack().Shuffle()
			case termbox.KeyCtrlL:
				g.showLegend = !g.showLegend
			case termbox.KeyCtrlT:
				g.board.showScore = !g.board.showScore
			case termbox.KeyCtrlF:
				g.showHint()
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				g.editbox.deleteRuneBackward()
			case termbox.KeySpace:
				g.editbox.insertRune(' ')
			case termbox.KeyEsc, termbox.KeyCtrlC:
				break mainloop
			default:
				if ev.Ch != 0 {
					g.editbox.insertRune(ev.Ch)
				}
			}
		case termbox.EventResize:
			g.draw()
		case termbox.EventMouse:
			if ev.Key == termbox.MouseLeft || ev.Key == termbox.MouseDown {
				g.draw()
				g.board.highlightPos(ev.MouseX, ev.MouseY)
				pos := g.board.pos(ev.MouseX, ev.MouseY)
				if pos == "" {
					break
				}
				g.editbox.clear()
				for _, r := range pos {
					g.editbox.insertRune(r)
				}
				g.editbox.insertRune(' ')
				termbox.Flush()
				continue
			}
			if ev.Key == termbox.MouseRelease {
				g.draw()
				continue
			}
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
		g.editbox.err = true
		return
	}
	var move quackle.Move
	if place == "-" {
		move = quackle.MoveCreatePassMove()
	} else {
		move = quackle.MoveCreatePlaceMove(place, flexAbc.Encode(word))
	}

	// score the move
	pos := g.pos()

	// known move.
	if pos.Moves().Contains(move) {
		pos.ScoreMove(move)
		g.qg.SetCandidate(move)
		g.qg.CommitMove(move)
		g.editbox.clear()
		return
	}

	// validate unknown move
	validityFlags := pos.ValidateMove(move)
	if validityFlags == int(quackle.GamePositionValidMove) {
		pos.ScoreMove(move)
		pos.AddAndSetMoveMade(move)
		g.qg.CommitMove(move)
		g.editbox.clear()
		return
	}

	// could not validate the move. reasons?
	//
	// very first move doesn't cover the center star
	if validityFlags&int(quackle.GamePositionInvalidOpeningPlace) > 0 {
		g.editbox.warn = true
		g.setErr("CENTER STAR")
	}
	// word doesn't connect to other plays on board
	if validityFlags&int(quackle.GamePositionInvalidPlace) > 0 {
		g.editbox.warn = true
		g.setErr("COME CLOSE")
	}
	// there are missing tiles in the rack
	if validityFlags&int(quackle.GamePositionInvalidTiles) > 0 {
		g.editbox.warn = true
		g.setErr("IMAGINARY")
	}
	// there is no such word mate
	if validityFlags&int(quackle.GamePositionUnacceptableWord) > 0 {
		g.editbox.warn = true
		g.setErr("IS IT A WORD?")
	}
	// invalid action
	if validityFlags&int(quackle.GamePositionInvalidAction) > 0 {
		g.editbox.warn = true
		g.setErr("SAD FACE")
	}
}

func (g *game) showHint() {
	g.editbox.clear()
	// accept top scoring advice from the beast
	move := g.pos().StaticBestMove()
	for _, r := range move.ToString() {
		g.editbox.insertRune(r)
	}
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

func (g *game) setErr(err string) {
	g.err = fmt.Errorf(err)
}

// newGame initializes a new game and constructs game object.
func newGame(opts *gameOpts) *game {
	if opts == nil {
		opts = &gameOpts{
			p1type:   human,
			p2type:   computer,
			p1name:   "iby",
			p2name:   "hal",
			alphabet: alphabet,
		}
	} else {
		switch {
		case opts.alphabet == "":
			opts.alphabet = alphabet
		}
	}

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
	player1 := quackle.NewPlayer(opts.p1name, 1, 0)
	player2 := newCompPlayer("cpu", 1)
	players := quackle.NewPlayerList()
	players.Add(player1)
	players.Add(player2)
	g.SetPlayers(players)
	g.AssociateKnownComputerPlayers()
	g.AddPosition()

	b := board{
		qb:         g.CurrentPosition().(quackle.GamePosition).Board(),
		w:          boardsize,
		h:          boardsize,
		curPosVert: true,
	}

	return &game{
		qg:      g,
		board:   b,
		rack1:   newRack(player1),
		rack2:   newRack(player2),
		editbox: newEditbox(),
	}
}

type playerType int

const (
	human playerType = iota
	computer
)

type gameOpts struct {
	p1type, p2type playerType
	p1name, p2name string
	alphabet       string
}
