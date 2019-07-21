package bus

type exchangeType int

const (
	fanoutExchange = iota
)

func (k exchangeType) string() string {
	return [...]string{"fanout",}[k]
}

