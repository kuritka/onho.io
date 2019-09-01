package msgbus

type exchange int

const (
	serviceDiscoveryExchange = iota
	serviceEventExchange
	serviceCommandExchange
	exchangeWorkerQueue
)

func (k exchange) string() string {
	return [...]string{"serviceDiscoveryExchange", "serviceEventExchange", "serviceCommandExchange",""}[k]
}
