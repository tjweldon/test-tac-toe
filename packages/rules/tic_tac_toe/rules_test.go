package tictactoe

import (
	"testing"
	"tictac/packages/rules"
	"tictac/packages/util"
)

func MustInitialise(constructor func() rules.Rules[Move, BoardDisplayType], t *testing.T) TicTacToeRules {
	game := constructor() 
	if game == nil {
		t.Logf("The game was nil, expected a new Game instance")
		t.FailNow()
	}

	return game
}

func MustAllowMove(game TicTacToeRules, t *testing.T, move Move) {
	if err := game.MakeMove(move); err != nil {
		t.Fatal(err)
		t.FailNow()
	}
}

func MustNotAllowMove(game TicTacToeRules, t *testing.T, move Move) {
	if err := game.MakeMove(move); err == nil {
		t.Logf("Expected an error when attempting an illegal move, got none")
		t.FailNow()
	}
}

func MustBeOver(game TicTacToeRules, t *testing.T) {
	if !game.IsOver() {
		t.Logf("The game should have been over, it is not")
		t.FailNow()
	}
}

func MustNotBeOver(game TicTacToeRules, t *testing.T) {
	if game.IsOver() {
		t.Logf("The game should not be over, it was")
		t.FailNow()
	}
}

func MustHaveWinner(game TicTacToeRules, t *testing.T, victor util.Optional[string]) {
	MustBeOver(game, t)
	if game.GetWinner().IsNone() != victor.IsNone() {
		t.Logf("The game was a draw, expected '%s' to win", victor.Some())
		t.FailNow()
	}

	if game.GetWinner().IsSome() && game.GetWinner().Some() != victor.Some() {
		t.Logf("The game was won by '%s', expected '%s' to win", game.GetWinner().Some(), victor.Some())
		t.FailNow()
	}
}

func TestGame_StartsWithCrosses(t *testing.T) {
	game := MustInitialise(InitialiseGame, t)

	if current := game.GetCurrentPlayer(); current != "X" {
		t.Logf("The first player was '%s', expected 'X'", current)
		t.FailNow()
	}

	if isOver := game.IsOver(); isOver {
		t.Logf("The game should not be over")
	}
}

type MoveTestCase struct {
	Name string
	Move Move
	ExpectedState State
}

var MoveTestCases = []MoveTestCase{
	{
		Name: "Top Left",
		Move: Move{0, 0},
		ExpectedState: State{
			{"X", Pt, Pt},
			{Pt, Pt, Pt},
			{Pt, Pt, Pt},
		},
	},
	{
		Name: "Bottom Left",
		Move: Move{0, 2},
		ExpectedState: State{
			{Pt, Pt, Pt},
			{Pt, Pt, Pt},
			{"X", Pt, Pt},
		},
	},
	{
		Name: "Top Right",
		Move: Move{2, 0},
		ExpectedState: State{
			{Pt, Pt, "X"},
			{Pt, Pt, Pt},
			{Pt, Pt, Pt},
		},
	},
}

func TestGame_PlacesCrossCorrectly_WhenFirstMoveIsProvided(t *testing.T) {
	game := MustInitialise(InitialiseGame, t)
	for _, testCase := range MoveTestCases {
		game.NewGame()
		MustAllowMove(game, t, testCase.Move)
		expectedBoard := testCase.ExpectedState.Display()
		if board := game.GetBoard(); board != expectedBoard {
			t.Logf("Test Case %s\nboard state was \n%s\nexpected \n%s\n", testCase.Name, board, expectedBoard)
			t.Fail()
		}
	}
}



func TestGame_DoesNotAllowMove_WhenOutOfBoundsMoveIsChosen(t *testing.T) {
	game := MustInitialise(InitialiseGame, t)

	game.NewGame()
	MustNotAllowMove(game, t, Move{-1, 0})
}

func TestGame_CurrentPlayerHasChanged_WhenMoveHasBeenMade(t *testing.T) {
	game := MustInitialise(InitialiseGame, t)

	game.NewGame()
	initialPlayer := game.GetCurrentPlayer()
	MustAllowMove(game, t, Move{0, 0})

	if game.GetCurrentPlayer() == initialPlayer {
		t.Logf("The current player had not changed after one move")
		t.FailNow()
	}

	MustAllowMove(game, t, Move{0, 1})

	if game.GetCurrentPlayer() != initialPlayer {
		t.Logf("The current player should have been %s after two moves, got %s", initialPlayer, game.GetCurrentPlayer())
		t.FailNow()
	}
}

func TestGame_DoesNotAllowMove_WhenOccupiedCellIsChosen(t *testing.T) {
	game := MustInitialise(InitialiseGame, t)
	MustAllowMove(game, t, Move{0, 0})
	MustNotAllowMove(game, t, Move{0, 0})
}

func TestGame_IsOver_WhenAllCellsAreFilled(t *testing.T) {
	game := MustInitialise(InitialiseGame, t)
	
	for i := range [3]struct{}{} {
		for j := range [3]struct{}{} {
			MustAllowMove(game, t, Move{i, j})
		}
	}

	MustBeOver(game, t)
}

type GameFixture struct {
	Moves []Move
	ExpectedWinner util.Optional[string]
}

var GameScenarios = []GameFixture{
	{
		Moves: []Move{
			{0, 0}, {1, 0},
			{0, 1}, {1, 1},
			{0, 2}, // X wins with a column here
		},
		ExpectedWinner: util.Some("X"),
	},
	{
		Moves: []Move{
			{0, 0}, {0, 1},
			{1, 0}, {1, 1},
			{2, 0}, // X wins with a row here
		},
		ExpectedWinner: util.Some("X"),
	},
	{
		Moves: []Move{
			{0, 0}, {0, 2},
			{1, 0}, {1, 1},
			{0, 1}, {2, 0}, // O wins with a diagonal here
		},
		ExpectedWinner: util.Some("O"),
	},
	{
		Moves: []Move{
			{0, 0}, {0, 1},
			{1, 0}, {2, 0},
			{1, 1}, {1, 2},
			{2, 1}, {2, 2},
			{0, 2}, // draw
		},
		ExpectedWinner: util.None[string](),
	},
}

func TestGame_HasExpectedResult_WhenMoveSequencePlaysOut(t *testing.T) {
	for _, fixture := range GameScenarios {
		game := MustInitialise(InitialiseGame, t)
		for _, move := range fixture.Moves {
			MustAllowMove(game, t, move)
		}

		MustHaveWinner(game, t, fixture.ExpectedWinner)
	}
}

