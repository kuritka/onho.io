package qutils

import (
	"github.com/kuritka/onho.io/common/utils"
	"github.com/streadway/amqp"
	)

type IMessageConsumer interface {

}

type MessageConsumer struct {
	channel *amqp.Channel
	connection *amqp.Connection
}


type UniqueQueue struct{
	queue *amqp.Queue
	channel *amqp.Channel
}

type ChannelConsumer struct{
	queue *amqp.Queue
	channel *amqp.Channel
}

func NewMessageConsumer(connection *amqp.Connection,channel *amqp.Channel) *MessageConsumer{
	return &MessageConsumer{
		channel, connection,
	}
}



func (m *MessageConsumer) GetUniqueQueue() *UniqueQueue {
	//empty queue name = rabbit creates unique name for it
	q := GetQueue("", m.channel,true)
	return &UniqueQueue {
		q,
		m.channel,
	}
}

func (m *MessageConsumer) GetQueue(name string) *ChannelConsumer {
	return m.getQueue(name, true)
}

func (m *MessageConsumer) GetPersistentQueue(name string) *ChannelConsumer {
	return m.getQueue(name, false)
}


func (u *UniqueQueue) BindToFanout() *ChannelConsumer {
	return u.BindToExchange(Fanout)
}


func (u *UniqueQueue) BindToExchange(exchange string) *ChannelConsumer {
	//rebinding queue to fanout exchange
	err := u.channel.QueueBind( u.queue.Name,
		"",				//one queue could be bounded to one exchange several times and all bounds will work
		exchange,
		false,
		nil,
	)
	utils.FailOnError(err, "binding to fanout")
	return &ChannelConsumer{
		u.queue,u.channel,
	}
}


func (c *ChannelConsumer) ConsumeFromChannel() (<-chan amqp.Delivery, error) {
	return c.channel.Consume(c.queue.Name,
		"", true,false,false,false, nil)
}

//only one service could exclusively read from the channel
func (c *ChannelConsumer) ConsumeExclusiveFromChannel(autoAck bool) (<-chan amqp.Delivery, error) {
	return c.channel.Consume(c.queue.Name,
		"", autoAck,true,false,false, nil)
}



func (m *MessageConsumer) getQueue(name string, autoDelete bool) *ChannelConsumer {
	//empty queue name = rabbit creates unique name for it
	q := GetQueue(name, m.channel,autoDelete)
	return &ChannelConsumer{q,m.channel}
}