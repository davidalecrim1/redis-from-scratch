package internal

type Command interface{}

type SetCommand struct {
	Key, Val []byte
}

type GetCommand struct {
	Key []byte
}

type PingCommand struct{}

type HelloCommand struct {
	Value string
}

type ClientCommand struct {
	Value string
}
