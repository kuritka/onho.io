package msgbus

type DiscoveryRequest struct{
	CommandQueue string
	ServiceGuid string
}

type Message struct {
	Name    string
	Message string
}