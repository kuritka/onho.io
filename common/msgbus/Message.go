package msgbus

type DiscoveryRequest struct{
	CommandHandlers []string
	CommandQueue string
	ServiceGuid string
}

type Message struct {
	Name    string
	Message string
}