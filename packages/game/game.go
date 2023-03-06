package game

import (
	"log"
	"tictac/packages/adaptors"
	"tictac/packages/rules"
)

// MustDo attempts to invoke the function f, fatal if an error is returned
func MustDo(f func() error) {
	if err := f(); err != nil {
		log.Fatal(err)
	}
}

// WillRetryDo attempts to do the passed function f. On error, will pass the error to the handler provided and retry.
// If the handler returns an error it is fatal
func WillRetryDo(f func() error, failHandler func(error) error) {
	for err := f(); err != nil; err = f() {
		MustDo(func() error { return failHandler(err) })
	}
}

// MustGet attempts to obtain the result of calling get, fatal if the error is not nil
func MustGet[T any](get func() (T, error)) (result T) {
	var err error
	if result, err = get(); err != nil {
		log.Fatal(err)	
	}

	return result
}

// WillRetryGet attempts to retrieve the result using get, if the error is not nil it is passed to the handler provided and the
// the get is retried. If the failHandler returns an error it is fatal
func WillRetryGet[T any](get func() (T, error), failHandler func(error) error) T {
	result, err := get()
	for ; err != nil; result, err = get() {
		MustDo(func() error { return failHandler(err) })
	}

	return result
}

type Game[MoveType any, BoardType comparable] struct {
	UI adaptors.UIAdaptor[MoveType, BoardType]
	State rules.Rules[MoveType, BoardType]
}

func NewGame[MoveType any, BoardType comparable](
	UI adaptors.UIAdaptor[MoveType, BoardType], 
	constructor func() rules.Rules[MoveType, BoardType],
) *Game[MoveType, BoardType] {
	return &Game[MoveType, BoardType]{
		UI: UI,
		State: constructor(),
	}
}

func (g *Game[MoveType, BoardType]) Run() {
	MustDo(g.UI.OnGameStart)
	for !g.State.IsOver() {
		currentPlayer := g.State.GetCurrentPlayer()
		board := g.State.GetBoard()

		onTurnStart := func() error {
			return g.UI.OnTurnStart(currentPlayer, board)
		}
		MustDo(onTurnStart)
		
		doPlayerMove := func () error {
			return g.State.MakeMove(g.UI.ReceiveInput())
		}

		WillRetryDo(doPlayerMove, g.UI.OnInvalidMove)

		MustDo(g.UI.OnTurnEnd)
	}

	winnerOpt := g.State.GetWinner()
	var onGameEnd func() error
	if winner, isDraw := winnerOpt.Some(), winnerOpt.IsNone(); isDraw {
		onGameEnd = func () error {
			return g.UI.OnGameDrawn(g.State.GetBoard())
		}
	} else {
		onGameEnd = func() error {
			return g.UI.OnGameWon(g.State.GetBoard(), winner)
		}
	}

	MustDo(onGameEnd)
}
