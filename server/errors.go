package server

type CmdError string

const (
	CmdErrorInvalidCommand CmdError = "invalid command"
	CmdErrorIDAlreadyTaken CmdError = "id already taken"
)
