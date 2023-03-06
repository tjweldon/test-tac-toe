package cli

import (
	"fmt"
	"log"
	"strconv"
	"tictac/packages/adaptors"
	tictactoe "tictac/packages/rules/tic_tac_toe"
)


type CliAdaptor struct {}

func (c *CliAdaptor) OnGameStart() error {
	_, err := fmt.Println("Game on!")

	return err
}

func (c *CliAdaptor) OnTurnEnd() error {
	return nil
}

func (c *CliAdaptor) OnTurnStart(playerName string, board tictactoe.BoardDisplayType) error {
	if _, err := fmt.Println("Here's the board:", "\n"+fmt.Sprint(board)); err != nil {
		return err
	}

	if _, err := fmt.Println(playerName, "it's your turn"); err != nil {
		return err
	}

	return nil
}

func (c *CliAdaptor) ReceiveInput() tictactoe.Move {
	if _, err := fmt.Println("Please supply the column, then row for your move separated by a space (eg 0 2 for top right)"); err != nil {
		log.Fatal(err)
	}

	collectInput := func() (x, y int, err error) {
		var colStr, rowStr string
		_, err = fmt.Scanln(&colStr, &rowStr)
		if err != nil {
			log.Fatal(err)
		}
		inputError := fmt.Errorf(
			"please supply an input of the form 'x y' for x, y in 0-2\n",
		)

		if x, err = strconv.Atoi(colStr); err != nil {
			return x, y, inputError
		}
		if y, err = strconv.Atoi(rowStr); err != nil {
			return x, y, inputError
		}

		return x, y, nil
	}

	var (
		x, y int
		err error
	)
	for {
		y, x, err = collectInput()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		break
	} 
	
	return tictactoe.Move{X: x, Y: y}
}

func (c *CliAdaptor) OnInvalidMove(err error) error {
	_, err = fmt.Println(err.Error())
	return err
}

func (c *CliAdaptor) OnGameWon(board tictactoe.BoardDisplayType, winner string) error {
	if _, err := fmt.Println(board); err != nil {
		return err
	}
	if _, err := fmt.Println(winner, "wins!"); err != nil {
		return err
	}

	return nil
}

func (c *CliAdaptor) OnGameDrawn(board tictactoe.BoardDisplayType) error {
	if _, err := fmt.Println(board); err != nil {
		return err
	}
	if _, err := fmt.Println("It was a draw"); err != nil {
		return err
	}

	return nil
}

func NewTicTacToeCli() adaptors.UIAdaptor[tictactoe.Move, tictactoe.BoardDisplayType] {
	return &CliAdaptor{}
}


