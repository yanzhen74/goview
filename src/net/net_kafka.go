package net

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/yanzhen74/goview/src/model"
)

type NetKafka struct {
	name               string
	config             *sarama.Config
	consumer           sarama.Consumer
	partition_consumer sarama.PartitionConsumer
	subscribers        *[]chan string
}

func (this *NetKafka) Init(config *model.NetWork) (int, error) {
	fmt.Printf("init kafka network")

	conf := sarama.NewConfig()
	conf.Consumer.Return.Errors = true
	conf.Version = sarama.V0_11_0_2

	// consumer
	consumer, err := sarama.NewConsumer([]string{config.NetWorkIP}, conf)
	if err != nil {
		fmt.Printf("consumer kafka create error %s\n", err.Error())
		return -1, err
	}

	partition_consumer, err := consumer.ConsumePartition(config.NetWorkName, 0, sarama.OffsetNewest)
	if err != nil {
		fmt.Printf("try create partition_consumer error %s\n", err.Error())
		return -1, err
	}

	this.name = config.NetWorkName
	this.config = conf
	this.consumer = consumer
	this.partition_consumer = partition_consumer
	this.subscribers = new([]chan string)

	return 1, nil
}

func (this *NetKafka) Subscribe(name string, sub chan string) {
	if this.name == name {
		*this.subscribers = append(*this.subscribers, sub)
	}
}

func (this *NetKafka) Process() error {
	for {
		select {
		case msg := <-this.partition_consumer.Messages():
			fmt.Printf("msg offset: %d, partition: %d, timestamp: %s, value: %s\n",
				msg.Offset, msg.Partition, msg.Timestamp.String(), string(msg.Value))
		case err := <-this.partition_consumer.Errors():
			fmt.Printf("err :%s\n", err.Error())
		case <-time.After(100 * time.Millisecond):
			// to be deleted , just for test now
			for _, c := range *this.subscribers {
				c <- "hello world"
			}
		}
	}
}
