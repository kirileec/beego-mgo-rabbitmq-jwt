package rabbitmq

//rabbitmq Message Queue functions

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

//if err output
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello world"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}

// public publish message
//
//@amqpURI amqp uri
//@exchange exchange name
//@key routing key
//@body content
func publish(amqpURI string, exchange string, key string, body []byte) {
	//dial to MQ
	log.Printf("dialing %q", amqpURI)
	connection, err := amqp.Dial(amqpURI)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer connection.Close()

	//create channel
	log.Printf("got Connection, getting Channel")
	channel, err := connection.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	log.Printf("publishing %dB body (%q)", len(body), body)

	err = channel.Publish(
		exchange, // exchange name
		key,      // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			DeliveryMode:    amqp.Persistent,
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            body,
		})
	failOnError(err, "Failed to publish a message")
}

//consumer method
//
//@amqpURI, amqp uri
//@exchange, exchange name
//@exchangeType, exchange type eg. direct|fanout|topic
//@queue, queue name
//@key routing key
//@callback callback function when new message received
func consumer(amqpURI string, exchange string, exchangeType string, queue string, key string, callback func([]byte)) {
	//dial to MQ
	log.Printf("dialing %q", amqpURI)
	connection, err := amqp.Dial(amqpURI)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer connection.Close()

	//create channel
	log.Printf("got Connection, getting Channel")
	channel, err := connection.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	//declare a exchange
	log.Printf("got Channel, declaring Exchange (%q)", exchange)
	err = channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	)
	failOnError(err, "Exchange Declare:")

	//declare a queue
	q, err := channel.QueueDeclare(
		"exchange.queue", // queue name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive:when a consumer close connection，delete the queue
		false,            // no-wait
		nil,              // arguments
	)
	failOnError(err, "Failed to declare a queue")

	//binding to exchange then this queue can receive message from exchange with routing key
	err = channel.QueueBind(
		q.Name,   // name of the queue
		key,      // routing key
		exchange, // exchange name
		false,    // noWait
		nil,      // arguments
	)
	failOnError(err, "Failed to bind a queue")

	// log.Printf("Queue bound to Exchange, starting Consume")
	//consume message
	msgs, err := channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	//create a go channel
	forever := make(chan bool)

	//run by gorountine
	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
			callback(d.Body) //call the callback
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	//waiting the mesage forever
	<-forever
}
