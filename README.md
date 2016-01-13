# bingo

bingo is a scrabble game for terminal. it uses [quackle](https://github.com/quackle/quackle) as the
engine.

![](http://i.imgur.com/RvFeEyz.png)

## installation

requires recent version of `Go`,  `Qt4 QtCore` headers and `SWIG 3.0.6` or higher.

`go get github.com/igungor/bingo`

## how to play

human vs. cpu. cpu always wins!

- input box accepts quackle-format
  * `H2 NABER` means start from `H2` and place the word `NABER` from top to bottom.
  * `3B HELO` means start from `3B` and place the word `HELO` from left to right.
  * words must be typed in all uppercase, except jokers.
  * jokers must be typed in lowercase. `n` is the joker letter for move `4F CAMEKAn`.
  * type `-` to pass your turn.
  * use `.` if a letter of the word you type is already on the board.
- `ctrl-s` toggles multipliers and scores.
- `ctrl-t` proposes a highscore move for you.
- `ctrl-f` shuffles your rack.
- `ctrl-l` toggles legend.
- `ctrl-c` or `esc` quit the game.

## there are bugs

not just bugs, there are many shortcomings as well. it was just a proof-of-concept thingy but it's
a playable game right now. i play and lose every damn game.

## credits

[@gokceneraslan](https://github.com/gokceneraslan)

## license

same as the `quackle` license. see
[LICENSE](https://github.com/quackle/quackle/blob/master/LICENSE)
