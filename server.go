package main

import (
	"fmt"
	"net"
)

var (
	welcomeMsg = []byte(
		`
		welcome to my chat server :), commands you can execute:

		/nick <name> 		#for giving yourself a name
		/rooms <room_name>  #list available rooms
		/join <room_name> 	#create a room/join the room
		/msg <msg> 			#send a message
		/quit   			#quit the server


	` + "\n")
)

type Server struct {
	cmds  chan command
	rooms map[string]*Room
}

func createServer() *Server {

	return &Server{
		cmds:  make(chan command),
		rooms: make(map[string]*Room),
	}
}

func (s *Server) createClient(conn net.Conn) *client {
	c := &client{
		conn: conn,
		nick: "Anonymous",
		room: nil,
		cmds: &s.cmds,
	}

	c.conn.Write(welcomeMsg)
	go readFromClient(c)

	return c
}

func listenForCommands(s *Server) {
	for {
		for cmd := range s.cmds {
			switch cmd.id {
			case cmd_nick:
				s.nick(cmd)
			case cmd_rooms:
				s.listRooms(cmd)
			case cmd_join:
				s.join(cmd)
			case cmd_msg:
				s.broadcast(cmd)
			case cmd_quit:
				s.quit(cmd)
			}
		}
	}
}

func (s *Server) nick(cmd command) {
	cmd.client.msg(fmt.Sprintf("ok i will call u %s :)\n", cmd.msg))
	cmd.client.nick = cmd.msg
}
func (s *Server) listRooms(cmd command) {
	rooms := make([]string, 0)
	for _, room := range s.rooms {
		rooms = append(rooms, room.name)
	}
	cmd.client.msg(fmt.Sprintf("available rooms: %s\n", rooms))
}
func (s *Server) join(cmd command) {
	roomName := cmd.msg
	//client already in another room, new client
	if cmd.client.room != nil {
		cmd.client.room.broadcast(fmt.Sprintf("%s has quit the room\n", cmd.client.nick))
	}

	//room might not be present...creating room
	room, ok := s.rooms[roomName]
	if !ok {
		room = &Room{
			name:    roomName,
			members: make(map[string]*client),
		}
		s.rooms[roomName] = room
	}
	cmd.client.room = room
	room.members[cmd.client.conn.RemoteAddr().String()] = cmd.client
	cmd.client.msg(fmt.Sprintf("you have joined %s\n", roomName))
	room.broadcast(fmt.Sprintf("%s: has joined this room\n", cmd.client.nick))

}
func (s *Server) broadcast(cmd command) {
	cmd.client.room.broadcast(fmt.Sprintf("%s: %s\n", cmd.client.nick, cmd.msg))
}
func (s *Server) quit(cmd command) {
	room := cmd.client.room
	if room != nil {
		room.broadcast(fmt.Sprintf("%s has quit the room\n", cmd.client.nick))
		delete(room.members, cmd.client.conn.RemoteAddr().String())
	}

	cmd.client.msg("sad to see u go :(\n")
	cmd.client.conn.Close()
}
