package main

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/gocql/gocql"
)

const (
	kafkaConn        = "localhost:9092"
	topic            = "cust"
	cassandraCluster = "127.0.0.1"
)

var Session *gocql.Session
var Producer sarama.AsyncProducer
var Consumer sarama.Consumer

const (
	cassandraIP = "127.0.0.1"
	keyspace    = "customer"
)

func init() {
	var err error
	cluster := gocql.NewCluster(cassandraIP)
	cluster.Keyspace = keyspace
	Session, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	fmt.Println("Cassandra connection done")

	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Consumer.Return.Errors = true
	Producer, _ = sarama.NewAsyncProducer([]string{kafkaConn}, config)
	Consumer, _ = sarama.NewConsumer([]string{kafkaConn}, config)
}
