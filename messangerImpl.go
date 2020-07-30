package main

import "github.com/Shopify/sarama"

type KafkaMessanger struct {
}

func (kafka KafkaMessanger) PublishMessage(message string) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder([]byte(message)),
	}
	Producer.Input() <- msg
}

func (kafka KafkaMessanger) ConsumeMessage() {

}
