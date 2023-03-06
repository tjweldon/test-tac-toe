# Test-tac-toe

This is a ludicrously over-engineered test first implementation of tic tac toe (or naughts and crosses) built in golang.

## Project Goals

 - To practise some TDD in a strongly typed language.
 - To see how much more work it requires in terms of up front design to do this in go vs in python.
 - How using a language with a far less dynamic runtime than python affects the way tests are written and how that manifests the design of the software interfaces.
 - **Bonus Points:** Make the design sufficiently generic to acommodate other kinds of turn based games, and potentially other kinds of UI than the command line.

## Lessons Learned

Test first is a misnomer. You can't have a test fail at the compilation stage (it didn't even run), so really it's stub first, then test, then implementation. The thing about writing the stubs first is that it's often hard to know what stubs to write without attempting to write the client code, or having a strong idea of what the client code will be.

I found that writing stubs and tests impossible as a way to start, but thinking about what I needed in order to do that immensely valuable in terms of the design decisions it led me to, largely around dependencies. Had the project been written in python I would have been easily able to partially mock out the UI if it were just implemented as part of the top level Game struct. Maybe you can do unittest style patching in go, but I don't know how so instead I had the game service take a UIAdaptor that presented a generic interface and a ruleset with a similarly generic interface.

The UIAdaptor interface below is in principle agnostic to the i/o port that the application uses. It could operate over the network or, like the one I actually implemented, use the cli.
```golang
package adaptors

type UIAdaptor[MoveType, BoardType any] interface {
	OnGameStart() error
	OnTurnStart(playerName string, board BoardType) error
	ReceiveInput() MoveType
	OnInvalidMove(err error) error
	OnTurnEnd() error
	OnGameWon(board BoardType, winner string) error
	OnGameDrawn(board BoardType) error
}
```

The Rules interface is decomposed into commands which change the state (like a player making a move) and queries like IsOver or GetBoard that return some partial state information. We exclusively interact with the composed Rules interface but I found it helpful to separate these responisibilities in code, if only to minimise cognitive overhead.
```golang
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
```

The game loop then was written first and was implemented pretty much in advance of anything else because it let me validate the _utility_ (not the correctness, since there was no implementation) of the interfaces I had designed, and did not depend in any way whatsoever on the implementation.

I could then write stubs for a tic-tac-toe implementation as an implementation of the Rules interface, and express my expectations about the behavior of the rules in those tests.

Did I have fewer bugs in my code? no. Tests are as prone to bugs as the code itself and all a test can do is check what you tell it to check. If you get that wrong, then you implement the wrong behavior. What I can say for sure is that the testability constraint means that the implementation is better from the perspective of interface segregation and inversion of control. On the other hand it's way more complex than it needed to be. 

Ultimately I think it depends on the scope of the project and the details of the ecosystem as to whether the trade-offs make sense, as always.
