package controller

import (
	"bytes"
	"encoding/gob"
	"github.com/gorilla/websocket"
	"github.com/kuritka/onho.io/common/dto"
	"github.com/kuritka/onho.io/common/qutils"
	"github.com/kuritka/onho.io/common/utils"
	"github.com/kuritka/onho.io/services/webadmin/model"
	"github.com/streadway/amqp"
	"net/http"
	"sync"
)

type websocketController struct {
	conn             *amqp.Connection
	ch               *amqp.Channel
	sockets          []*websocket.Conn
	mutex            sync.Mutex
	upgrader         websocket.Upgrader
	publisher        *qutils.MessageProvider
	sourceSubscriber *qutils.MessageConsumer
	messageSubscriber *qutils.MessageConsumer
}



func newWebsocketController(options Options) *websocketController {
	wsc := new(websocketController)
	wsc.conn, wsc.ch = qutils.GetChannel(options.QueueConnectionString)
	wsc.upgrader = websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}
	wsc.publisher = qutils.NewGobMessageProvider()
	wsc.sourceSubscriber = qutils.NewMessageConsumer(wsc.conn, wsc.ch)
	wsc.messageSubscriber = qutils.NewMessageConsumer(wsc.conn, wsc.ch)
	go wsc.listenForSources()
	go wsc.listenForMessages()
	return wsc
}



func (wsc *websocketController) handleMessage(w http.ResponseWriter, r *http.Request) {
	socket, _ := wsc.upgrader.Upgrade(w,r, nil)
	wsc.addSocket(socket)
	go wsc.listenForDiscoveryRequest(socket)
}

func (wsc *websocketController) addSocket(socket *websocket.Conn) {
	wsc.mutex.Lock()
	wsc.sockets = append(wsc.sockets , socket)
	wsc.mutex.Unlock()
}

func (wsc *websocketController) removeSocket(socket *websocket.Conn) {
	wsc.mutex.Lock()
	socket.Close()
	for i := range wsc.sockets {
		if wsc.sockets[i] == socket {
			wsc.sockets = append(wsc.sockets[:i], wsc.sockets[i+1:]...)
		}
	}
	wsc.mutex.Unlock()
}



func (wsc *websocketController) listenForDiscoveryRequest(socket *websocket.Conn) {
	msg := message{}
	err := socket.ReadJSON(&msg)
	if err != nil {
		wsc.removeSocket(socket)
		return
	}
	if msg.Type =="discover"{
		//it is just trigger request, message it self could be empty, coordinator respond by list of sensord
		wsc.publisher.AsAmqpMessage("").
			PublishToCustomExchange(wsc.ch,"", qutils.WebAppDiscoveryQueue)
	}
}


func (wsc *websocketController) listenForSources(){
	msgs, err:= wsc.sourceSubscriber.GetUniqueQueue().
		BindToExchange(qutils.WebAppSourceExchange).
		ConsumeFromChannel()
	utils.FailOnError(err, "reading sources from " + qutils.WebAppSourceExchange)
	for msg := range msgs {
		sensor, _ := model.GetSensorByName(string(msg.Body))
		//send to the client
		wsc.sendMessage(message{
			Type: "source",
			Data: sensor,
		})
	}
}

func (wsc *websocketController) sendMessage(msg message) {
	socketsToRemove := []*websocket.Conn{}
	for _, socket := range wsc.sockets {
		err := socket.WriteJSON(msg)
		if err != nil {
			socketsToRemove = append(socketsToRemove, socket)
		}
	}

	for _, socket := range wsc.sockets {
		wsc.removeSocket(socket)
	}
}

func (wsc *websocketController) listenForMessages() {
	msgs, _:= wsc.messageSubscriber.GetUniqueQueue().
		BindToExchange(qutils.WebAppReadingsExchange).
		ConsumeFromChannel()

	for msg := range msgs {
		buf := bytes.NewBuffer(msg.Body)
		dec := gob.NewDecoder(buf)
		sm := dto.SensorMessage{}
		dec.Decode(&sm)

		wsc.sendMessage(message{ Type: "reading" , Data: sm, })
	}
}

//message sending and receiving from websockets
type message struct {
	Type string `json:"type"`
	Data interface{} `json:"Data"`
}



