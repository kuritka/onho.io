//package bus
//
//import (
//	"crypto/rand"
//	"fmt"
//	"github.com/kuritka/onho.io/common/utils"
//	"github.com/rs/zerolog/log"
//	"github.com/streadway/amqp"
//)
//
//
//type IProvider interface {
//	Register(name string) *RegistredProviderImpl
//	Close()
//}
//
//
//type IRegistredProvider interface {
//	PublishEvent(name string) *RegistredProviderImpl
//	PublishCommand(name string) *RegistredProviderImpl
//	Listen()
//}
//
//
//type ProviderImpl struct {
//	connection *amqp.Connection
//	channel *amqp.Channel
//	qmanager *queueManagerImpl
//}
//
//
//
//type RegistredProviderImpl struct {
//	provider *ProviderImpl
//}
//
//
//type RegisterInfo struct {
//	Name string
//	Type string
//}
//
//
//func NewProvider(connection *amqp.Connection,channel *amqp.Channel) *ProviderImpl{
//	utils.FailOnNil(connection,"connection cannot be nil")
//	utils.FailOnNil(channel,"channel cannot be nil")
//	qmanager := newQueueManager(connection, channel)
//	return &ProviderImpl{
//		connection,
//		  channel,
//		  qmanager,
//	}
//}
//
//func (p *ProviderImpl) Register(name string) *RegistredProviderImpl {
//	utils.FailOnEmptyString(name, "name cannot be nil")
//
//	p.qmanager.
//		createExchangeIfNotExists(serviceDiscoveryExchange, fanoutExchange,true).
//		sendDiscoveryEvent(name)
//
//	//	createQueue(name, true).
//	//	bindToQueue().
//	return &RegistredProviderImpl {
//		p,
//	}
//}
//
//func (p *ProviderImpl) Listen() {
//
//
//}
//
//
////func ListenDiscovery(p *RegistredProviderImpl){
////	p.provider..
////}
//
//
//func (p *ProviderImpl) Close(){
//	var err error
//	err = p.channel.Close()
//	utils.FailOnError(err, "unable to close channel")
//	err = p.connection.Close()
//	utils.FailOnError(err, "unable to close connection")
//	log.Debug().Msg("connection closed")
//}
//
//
//func getGuid() (string,error) {
//	b := make([]byte, 16)
//	_, err := rand.Read(b)
//	if err != nil {
//		return "",err
//	}
//	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
//		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
//	return uuid, nil
//}