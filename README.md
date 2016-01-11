# bingo

bingo is a scrabble game for terminal. it uses [quackle](https://github.com/quackle/quackle) as the
engine.

## installation

`go get github.com/igungor/bingo`

## how to play

human vs. cpu. cpu always wins!

- `ctrl-s` toggles scores.
- input box accepts quackle-format
  * `H2 NABER` means start from `H2` and place the word `NABER` from top to bottom.
  * `3B HELO` means start from `3B` and place the word `HELO` from left to right.
- type `-` to pass your turn.

## there are bugs

not just bugs, there are many shortcomings as well. it was just a proof-of-concept thingy but it's
a playable game right now. i play and lose every damn game.

## credits

[@gokceneraslan](https://github.com/gokceneraslan)

## license

same as the `quackle` license. see
[LICENSE](https://github.com/quackle/quackle/blob/master/LICENSE)
