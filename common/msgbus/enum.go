package msgbus


type exchangeType int

const (
	fanoutExchange = iota
	directExchange
)

func (k exchangeType) string() string {
	return [...]string{"fanout","direct"}[k]
}






type exchange int

const (
	serviceDiscoveryExchange = iota
	serviceEventExchange
	serviceCommandExchange
)

func (k exchange) string() string {
	return [...]string{"serviceDiscoveryExchange","serviceEventExchange","serviceCommandExchange"}[k]
}



