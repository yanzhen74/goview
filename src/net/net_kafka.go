package net

import (
	"fmt"
	"reflect"
	"time"

	"github.com/Shopify/sarama"
	"github.com/yanzhen74/goview/src/model"
)

type NetKafka struct {
	name               string
	config             *sarama.Config
	consumer           sarama.Consumer
	partition_consumer sarama.PartitionConsumer
	subscribers        *[]*model.FrameType
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
	this.subscribers = new([]*model.FrameType)

	return 1, nil
}

func (this *NetKafka) Subscribe(sub *model.FrameType) {
	if this.name == sub.DataType {
		*this.subscribers = append(*this.subscribers, sub)
	}
}

// receive net frame, parse and dispatch
func (this *NetKafka) Process() error {
	ticker := time.NewTicker(time.Millisecond * time.Duration(100))
	cases := init_cases(this.partition_consumer.Messages(),
		this.partition_consumer.Errors(),
		ticker)
	for {
		chose, value, _ := reflect.Select(cases)

		switch chose {
		case 0: // chan_msg
			msg := (value.Interface()).(*sarama.ConsumerMessage)
			fmt.Printf("msg offset: %d, partition: %d, timestamp: %s, value: %s\n",
				msg.Offset, msg.Partition, msg.Timestamp.String(), string(msg.Value))
			// to be fixed
			if len(cases) > 3+len((*(this.subscribers))) {
				cases = cases[:len(cases)-len((*(this.subscribers)))]
			}
			this.send_to_subscribers(&cases, string(msg.Value))

		case 1: // chan_err
			err := (value.Interface()).(*sarama.ConsumerError)
			fmt.Printf("err :%s\n", err.Error())
		case 2: // timer
			// to be deleted , just for test now
			//for _, c := range *this.subscribers {
			// c.NetChanFrame <- "hello world"
			// }
		default: // send ok
			cases = append(cases[:chose], cases[chose+1:]...)
		}
	}
}

func init_cases(
	chan_msg <-chan *sarama.ConsumerMessage,
	chan_err <-chan *sarama.ConsumerError,
	ticker *time.Ticker) (cases []reflect.SelectCase) {

	// chan msg
	selectcase := reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(chan_msg),
	}
	cases = append(cases, selectcase)

	// chan err
	selectcase = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(chan_err),
	}
	cases = append(cases, selectcase)

	// timer
	selectcase = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(ticker.C),
	}
	cases = append(cases, selectcase)

	return
}

func (this *NetKafka) send_to_subscribers(
	cases *[]reflect.SelectCase,
	send_value interface{}) {

	// 每个消费者，发送一次后必须删除
	for _, item := range *(this.subscribers) {
		selectcase := reflect.SelectCase{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf((*item).NetChanFrame),
			Send: reflect.ValueOf(send_value),
		}
		*cases = append(*cases, selectcase)
	}
}
