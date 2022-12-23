package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	conn net.Conn
	nick string
	room *Room
	cmds *chan command
}

func readFromClient(c *client) {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		args := strings.Fields(msg)

		switch args[0] {
		case cmd_nick.String():
			*c.cmds <- command{
				client: c,
				id:     cmd_nick,
				msg:    strings.Join(args[1:], " "),
			}
		case cmd_rooms.String():
			*c.cmds <- command{
				client: c,
				id:     cmd_rooms,
				msg:    strings.Join(args[1:], " "),
			}
		case cmd_join.String():
			*c.cmds <- command{
				client: c,
				id:     cmd_join,
				msg:    strings.Join(args[1:], " "),
			}
		case cmd_msg.String():
			*c.cmds <- command{
				client: c,
				id:     cmd_msg,
				msg:    strings.Join(args[1:], " "),
			}
		case cmd_quit.String():
			*c.cmds <- command{
				client: c,
				id:     cmd_quit,
			}
		default:
			c.msg(fmt.Sprintf("invalid command :(\nvalid commands list:%s\n", welcomeMsg))
		}
	}
}

func (c *client) err(msg string) {
	fmt.Fprintf(c.conn, msg)
}

func (c *client) msg(message string) {
	fmt.Fprintf(c.conn, message)
}
