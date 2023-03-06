package tictactoe

import (
	"fmt"
	"log"
	"strings"

	"tictac/packages/rules"
	"tictac/packages/util"
)

type Move struct {X, Y int}

type BoardDisplayType string

type TicTacToeRules rules.Rules[Move, BoardDisplayType]

const Pt = "·"

type WinCondition [3]Move

func (w WinCondition) Check(s State) (winner util.Optional[string]) {
	winner = util.None[string]()
	tokenSet := w.BuildTokenSet(s)
	isWon := len(tokenSet) == 1
	if isWon {
		var winnerName string
		// the condition for a win is that there is only one key in the map 
		// so this just assigns the only key to winnerName
		for winnerName = range tokenSet {}
		winner = util.Some(winnerName)
	}

	if winner.Some() == Pt {
		return util.None[string]()
	}

	return winner
}

func (w WinCondition) BuildTokenSet(s State) map[string]struct{} {
	tokenSet := map[string]struct{} {}
	for _, cell := range w {
		token, err := s.GetCell(cell.X, cell.Y)
		if err != nil {
			log.Fatal(err)
		}

		tokenSet[token] = struct{}{}
	}

	return tokenSet
}

func (w WinCondition) IsEliminated(s State) bool {
	tokenSet := w.BuildTokenSet(s)
	distinctTokens := 0
	for token := range tokenSet {
		if token == Pt {
			continue
		}
		distinctTokens++
	}

	return distinctTokens > 1
}


func (w WinCondition) String() string {
	token := "#"
	var err error
	state := MakeEmptyState()
	for _, move := range w {
		state, err = state.ApplyMove(token, move)
		if err != nil {
			log.Fatal(err)
		}
	}

	return state.String()
}

func GetWinConditions() []WinCondition {
	return []WinCondition{
		// column wins
		{{0, 0},{0, 1},{0, 2}},
		{{1, 0},{1, 1},{1, 2}},
		{{2, 0},{2, 1},{2, 2}},
		
		// row wins
		{{0, 0},{1, 0},{2, 0}},
		{{0, 1},{1, 1},{2, 1}},
		{{0, 2},{1, 2},{2, 2}},
		
		// diagonals
		{{0, 0},{1, 1},{2, 2}},
		{{0, 2},{1, 1},{2, 0}},
	}
}

type State [3][3]string

func (s State) String() (result string) {
	for _, row := range s {
		result += strings.Join(row[:], "")+"\n"
	}

	return result
}

func (s State) Display() BoardDisplayType {
	return BoardDisplayType(s.String())
}

func (s State) GetCell(x, y int) (cellValue string, err error) {
	if x < 0 || y < 0 || x >= 3 || y >= 3 {
		return cellValue, fmt.Errorf("The cell coordinates (%d, %d) are not valid. The state is a zero indexed 3x3 grid", x, y)
	}

	return s[y][x], nil
}

func (s State) CountEmptyCells() (count int) {
	for _, row := range s {
		for _, cell := range row {
			if cell == Pt {
				count++
			}
		}
	}

	return count
}

func (s State) ApplyMove(player string, m Move) (result State, err error) {
	if cell, err := s.GetCell(m.X, m.Y); err != nil {
		return State{}, err
	} else if cell != Pt {
		return State{}, fmt.Errorf("This cell is occupied, choose another")
	}
	
	result, err = MakeState(s.String())
	result[m.Y][m.X] = player
	return result, err
}

func MakeState(raw string) (s State, err error) {
	rows := [3]string{}
	if n := copy(rows[:], strings.Split(raw, "\n")); n < 3 {
		return State{}, fmt.Errorf("Not enough rows provided, got %d", n)
	}

	for rIdx, row := range rows {
		if m := copy(s[rIdx][:], strings.Split(row, "")); m != 3 {
			return State{}, fmt.Errorf("row %d copied %d items, expected 3", rIdx, m)
		}
	}

	return s, nil
}

func MakeEmptyState() State {
	return State{
		{Pt, Pt, Pt},
		{Pt, Pt, Pt},
		{Pt, Pt, Pt},
	}
}

type TicTacToe struct {
	currentPlayer string
	state State
	nextPlayer util.Generator[string]
	winner util.Optional[string]
	winConditions []WinCondition
}

func InitialiseGame() rules.Rules[Move, BoardDisplayType] {
	game := &TicTacToe{}
	game.NewGame()

	return game
}

// Game Queries implementation

func (g *TicTacToe) GetBoard() BoardDisplayType {
	return g.state.Display()
}

func (g *TicTacToe) GetCurrentPlayer() string {
	return g.currentPlayer
}

func (g *TicTacToe) pruneOutcomes(unreachable []bool) {
	prunedOutcomes := []WinCondition{}
	for idx, condition := range g.winConditions {
		if unreachable[idx] {
			continue
		}
		prunedOutcomes = append(prunedOutcomes, condition)
	}
	g.winConditions = prunedOutcomes
}

func (g *TicTacToe) IsOver() (over bool) {
	unreachable := make([]bool, len(g.winConditions))
	for cIdx, condition := range g.winConditions {
		// if the win condition cannot possibly occur, we can prune it and needn't check if it has occurred
		if unreachable[cIdx] = condition.IsEliminated(g.state); unreachable[cIdx] {
			continue
		}
		
		// otherwiswe we check if that condition is met
		g.winner = condition.Check(g.state)
		if g.winner.IsSome() {
			over = true
			break
		}
	}

	if !over {
		g.pruneOutcomes(unreachable)
	}

	over = over || len(g.winConditions) == 0 || g.state.CountEmptyCells() == 0

	return over
}

func (g *TicTacToe) GetWinner() util.Optional[string] {return g.winner}

// Game commands

func (g *TicTacToe) NewGame() {
	g.state = MakeEmptyState()
	nextPlayer, err := util.LoopFrom("X", "O")
	if err != nil {
		log.Fatal(err)
	}
	g.nextPlayer = nextPlayer
	g.currentPlayer = g.nextPlayer()
	g.winner = util.None[string]()
	g.winConditions = GetWinConditions()
}

func (g *TicTacToe) MakeMove(m Move) (err error) {
	var nextState State
	if nextState, err = g.state.ApplyMove(g.currentPlayer, m); err == nil {
		g.state = nextState
		g.currentPlayer = g.nextPlayer()
	}

	return func() error {
		return err
	}()
}


