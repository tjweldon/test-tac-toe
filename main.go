package main

import (
	"tictac/packages/cli"
	"tictac/packages/game"
	tictactoe "tictac/packages/rules/tic_tac_toe"
)

func main () {
	g := game.NewGame(
		cli.NewTicTacToeCli(),
		tictactoe.InitialiseGame,
	)
	g.Run()
}
