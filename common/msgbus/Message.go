package msgbus

type discoveryRequest struct{
	message string
	commandHandlers []string
}

type Message struct {
	Name    string
	Message string
}