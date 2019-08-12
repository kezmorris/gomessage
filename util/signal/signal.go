package signal

// used by client to signal server, requesting a connection
type greet struct {
	userID string
}

// used by server to respond to greet message
type greetReply struct {
	port   int
	userID string
}
