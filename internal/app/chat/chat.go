package chat

type GenericResponse struct {
	Text string
}

type ChatProvider interface {
	ChatProcess(message string) (GenericResponse, error)
}
