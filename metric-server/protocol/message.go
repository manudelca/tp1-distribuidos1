package protocol

type MsgType uint8

const (
	Success MsgType = iota
	FormatError
	ServerError
)

type Message struct {
	MsgType MsgType
	Message string
}
