package rabbitmq

// helper of Message Queue

import (
	"github.com/astaxie/beego"
)

// SendMsgToClient send message to client
//
//@bodyMsg message body
//@key routing key (client identify id)
func SendMsgToClient(bodyMsg []byte, key string) {
	amqpURI := beego.AppConfig.String("rabbitmq::url")
	exchangeName := beego.AppConfig.String("rabbitmq::exchangename")
	publish(amqpURI, exchangeName, key, bodyMsg)
}

// BroadCast send message to all client
//
//@bodyMsg message body to broadcast
func BroadCast(bodyMsg []byte) {
	// amqpURI := beego.AppConfig.String("rabbitmq::url")
	// exchangeName := beego.AppConfig.String("rabbitmq::exchangename")

	// var service svc.Service
	// if err := service.Prepare(); err != nil {
	// 	log.Error(err, "get db conn failed", "BroadCast")
	// 	return
	// }
	// results, err := services.GetClients(&service)
	// if err != nil {
	// 	log.CompletedError(err, "get clients failed", "services.GetClients")
	// 	return
	// }

	// for _, result := range results {
	// 	log.Info("BroadCast Message", "BroadCast", "key: %s", result.routingkey)
	// 	publish(amqpURI, exchangeName, result.Base.Serial, bodyMsg)
	// }

}

// BeginListenMq start listening message
//
//@callback callback that handle the received message
func BeginListenMq(callback func([]byte)) {
	//get the configurations
	queue := beego.AppConfig.String("rabbitmq::queuename")
	amqpURI := beego.AppConfig.String("rabbitmq::url")
	exchangeName := beego.AppConfig.String("rabbitmq::exchangename")
	exchangeType := beego.AppConfig.String("rabbitmq::exchangetype")
	key := beego.AppConfig.String("rabbitmq::routingkey")
	// start a thread for listen
	go consumer(amqpURI, exchangeName, exchangeType, queue, key, callback)
}
