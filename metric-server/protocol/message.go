package protocol

type MsgType uint8

const (
	SUCCESS MsgType = iota
	FORMATERROR
	SERVERERROR
)

type Message struct {
	MsgType MsgType
	Message string
}
