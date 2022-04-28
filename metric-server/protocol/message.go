package protocol

type MsgType uint8

const (
	Success MsgType = iota
	SuccessQueryResult
	FormatError
	ServerError
)

type Message struct {
	MsgType MsgType
	Message string
}
