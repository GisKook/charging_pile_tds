package mq

import (
	"github.com/bitly/go-nsq"
	"github.com/giskook/charging_pile_tds/conf"
	"log"
)

type NsqSocket struct {
	conf      *conf.NsqConfiguration
	Consumers []*NsqConsumer
}

func NewNsqSocket(config *conf.NsqConfiguration) *NsqSocket {
	return &NsqSocket{
		conf:      config,
		Consumers: make([]*NsqConsumer, 0),
	}
}

func (socket *NsqSocket) Start() {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()

	socket.ConsumerStart()
}

func (socket *NsqSocket) ConsumerStart() {
	for _, ch := range socket.conf.Consumer.Channels {
		consumer_conf := &NsqConsumerConf{
			Addr:    socket.conf.Consumer.Addr,
			Topic:   socket.conf.Consumer.Topic,
			Channel: ch,
			Handler: nsq.HandlerFunc(func(message *nsq.Message) error {
				data := message.Body
				RecvNsq(socket, data)

				return nil
			}),
		}

		consumer := NewNsqConsumer(consumer_conf)
		consumer.Start()
		socket.Consumers = append(socket.Consumers, consumer)
	}

}

func (socket *NsqSocket) Stop() {
	for _, consumer := range socket.Consumers {
		consumer.Stop()
	}
}
