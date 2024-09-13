package server

import (
	"errors"
	"fmt"
	"strings"
)

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

func CommandsFromTags(data string) []*Command {
	var cmds []*Command

	tags := strings.Split(data, "|>")
	for _, tag := range tags {
		cmd := NewCommand(CommandOptions{Code: CmdMessage})
		if strings.HasPrefix(tag, "<") {
			tagParts := strings.Split(tag, ":")
			cmd = NewCommand(
				CommandOptions{
					Code:  CmdCode(tagParts[0][1:]),
					Input: tagParts[1],
				},
			)
		}
		cmds = append(cmds, cmd)
	}

	return cmds
}

func (cmd *Command) Run(c *Connection) error {
	if serverInstance == nil {
		return errors.New("missing server instance: please start server")
	}

	if cmd.Code == CmdAction {
		actionParts := strings.SplitN(cmd.Input, " ", 1)
		cmd.Input = ""
		cmd.Code = CmdCode(actionParts[0])
		if len(actionParts) > 1 {
			cmd.Input = actionParts[1]
		}
	}

	switch cmd.Code {
	case CmdDisconnect, CmdClose:
		serverInstance.closeAndDeleteConnection(c)
		return nil
	case CmdMessage:
		for id, conn := range serverInstance.connections {
			if id != c.ID {
				if cmd.Input != "" {
					output := fmt.Sprintf(">message:%s", cmd.Input)
					if err := serverInstance.send(conn, c.ID, output); err != nil {
						serverInstance.closeAndDeleteConnection(conn)
						return err
					}
				}
			}
		}
	case CmdID:
		oldName := c.ID
		name := cmd.Input
		if serverInstance.hasID(name) {
			output := fmt.Sprintf(">message:%s", CmdErrorIDAlreadyTaken)
			if err := serverInstance.send(c, "", output); err != nil {
				serverInstance.closeAndDeleteConnection(c)
				return err
			}
			return nil
		}
		serverInstance.updateConnection(c, func(c *Connection) {
			c.ID = name
		})
		output := fmt.Sprintf(">message:%s|>>id:%s|>>old_id:%s", "id changed", name, oldName)
		if err := serverInstance.broadcast(output); err != nil {
			serverInstance.closeAndDeleteConnection(c)
			return err
		}
	default:
		output := fmt.Sprintf(">message:%s", CmdErrorInvalidCommand)
		if err := serverInstance.send(c, "", output); err != nil {
			serverInstance.closeAndDeleteConnection(c)
			return err
		}
	}

	return nil
}
