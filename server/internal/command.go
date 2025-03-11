package internal

type Command interface{}

type SetCommand struct {
	Key, Val []byte
}

type SetCommandWithExpiration struct {
	SetCommand
	ExpireMiliseconds int
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

type EchoCommand struct {
	Value string
}
