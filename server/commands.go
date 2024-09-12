package server

type CmdCode string

const (
	CmdMessage    CmdCode = "message"
	CmdDisconnect CmdCode = "disconnect"
	CmdClose      CmdCode = "close"
	CmdID         CmdCode = "id"
)
