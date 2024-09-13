package server

type CmdCode string

const (
	CmdMessage    CmdCode = "message"
	CmdDisconnect CmdCode = "disconnect"
	CmdClose      CmdCode = "close"
	CmdID         CmdCode = "id"
	CmdAction     CmdCode = "action"
)

type CommandOptions struct {
	Code  CmdCode
	Input string
}

type Command struct {
	Code  CmdCode
	Input string
}

func NewCommand(options CommandOptions) *Command {
	return &Command{Code: CmdCode(options.Code), Input: options.Input}
}
