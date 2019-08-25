package msgbus

type exchange int

const (
	serviceDiscoveryExchange = iota
	serviceEventExchange
	serviceCommandExchange
)

func (k exchange) string() string {
	return [...]string{"serviceDiscoveryExchange", "serviceEventExchange", "serviceCommandExchange"}[k]
}
