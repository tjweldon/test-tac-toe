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


