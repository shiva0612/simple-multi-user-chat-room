package main

type Room struct {
	name    string
	members map[string]*client
}

func (r *Room) broadcast(msg string) {
	for _, client := range r.members {
		client.msg(msg)
	}
}
