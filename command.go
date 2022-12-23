package main

type commandID int

const (
	cmd_nick commandID = iota
	cmd_rooms
	cmd_join
	cmd_msg
	cmd_quit
)

func (c commandID) String() string {
	switch c {
	case cmd_nick:
		return "/nick"
	case cmd_rooms:
		return "/rooms"
	case cmd_join:
		return "/join"
	case cmd_msg:
		return "/msg"
	case cmd_quit:
		return "/quit"
	}
	return "unknown"
}

type command struct {
	client *client
	id     commandID
	msg    string
}
