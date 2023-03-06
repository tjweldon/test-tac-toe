package rules

import (
	"tictac/packages/util"
)

type GameQueries[BoardType comparable] interface {
	GetBoard() BoardType
	GetCurrentPlayer() string
	IsOver() bool
	GetWinner() util.Optional[string]
}

type GameCommands[MoveType any] interface {
	NewGame()
	MakeMove(move MoveType) error
}

type Rules[MoveType any, BoardType comparable] interface {
	GameCommands[MoveType]
	GameQueries[BoardType]
}
