package server

type CmdCode string

const (
	CmdMessage    CmdCode = "message"
	CmdDisconnect CmdCode = "disconnect"
)
