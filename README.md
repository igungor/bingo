# bingo

bingo is a crossword game for terminal. it uses [quackle](https://github.com/quackle/quackle) as
the engine.

![](http://i.imgur.com/RvFeEyz.png)

## installation

requires recent version of `Go`,  `Qt4 QtCore` headers and `SWIG 3.0.6` or higher.

`go get -u github.com/igungor/bingo`

## how to play

human vs. cpu. cpu always wins!

- input box accepts quackle-format
  * `H2 NABER` means start from `H2` and place the word `NABER` from top to bottom.
  * `3B HELO` means start from `3B` and place the word `HELO` from left to right.
  * words must be typed in all uppercase, except jokers.
  * jokers must be typed in lowercase. `n` is the joker letter for move `4F CAMEKAn`.
  * type `-` to pass your turn.
  * use `.` if a letter of the word you type is already on the board.

## keymapping

|key                             |description                             |
|--------------------------------|----------------------------------------|
|<kbd>ctrl-t</kbd>               |toggles multipliers and scores          |
|<kbd>ctrl-f</kbd>               |fill the input box with a highscore move|
|<kbd>ctrl-s</kbd>               |shuffles the rack                       |
|<kbd>ctrl-l</kbd>               |toggles legend                          |
|<kbd>ctrl-c</kbd>,<kbd>esc</kbd>|quit the game                           |

## credits

[@gokceneraslan](https://github.com/gokceneraslan)

## license

same as the `quackle` license. see
[LICENSE](https://github.com/quackle/quackle/blob/master/LICENSE)
