package net

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/yanzhen74/goview/src/model"
)

type NetKafka struct {
	name               string
	config             *sarama.Config
	consumer           sarama.Consumer
	partition_consumer sarama.PartitionConsumer
	subscribers        map[string]*[]*model.FrameType
}

func (this *NetKafka) Init(config *model.NetWork) (int, error) {
	fmt.Printf("init kafka network")

	conf := sarama.NewConfig()
	conf.Consumer.Return.Errors = true
	conf.Version = sarama.V0_11_0_2
	conf.Consumer.MaxWaitTime = time.Duration(30) * time.Millisecond

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
	this.subscribers = make(map[string]*[]*model.FrameType)

	return 1, nil
}

func (this *NetKafka) Subscribe(sub *model.FrameType) {
	if this.name == sub.DataType {
		frame_type := get_frame_type(sub)
		log.Println(frame_type)
		if _, ok := this.subscribers[frame_type]; !ok {
			this.subscribers[frame_type] = new([]*model.FrameType)
		}
		*(this.subscribers[frame_type]) = append(*(this.subscribers[frame_type]), sub)
	}
}

// todo add all types of frame_type title here
func get_frame_type(sub *model.FrameType) (frame_type string) {
	types := make([]string, 0, 0)
	if sub.DataType == "RTM" {
		types = []string{sub.DataType, sub.MissionID, sub.SubAddressName, "Result"}
	} else {
		types = []string{sub.DataType, sub.MissionID, sub.PayloadName, sub.SubAddressName, "Result"}
	}

	frame_type = strings.Join(types, "_")
	return
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
			// match the subscribers
			matched_subscribers := this.get_matched_subscribers(string(msg.Value))
			if matched_subscribers == nil {
				continue
			}

			// to be fixed
			if len(cases) > 3+len((*(matched_subscribers))) {
				cases = cases[:len(cases)-len((*(matched_subscribers)))]
			}
			this.send_to_subscribers(&cases, string(msg.Value), matched_subscribers)

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
func (this *NetKafka) get_matched_subscribers(msg string) (matched_subscribers *[]*model.FrameType) {
	lines := strings.Split(msg, "\n")
	if len(lines) <= 0 {
		return nil
	}
	titles := strings.Split(lines[0], "\t")
	if len(titles) <= 0 {
		return nil
	}
	frame_type := titles[0]
	if len(frame_type) <= 0 {
		return nil
	}

	if v, ok := this.subscribers[frame_type]; ok {
		return v
	}

	return nil
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
	send_value interface{},
	matched_subscribers *[]*model.FrameType) {

	// 每个消费者，发送一次后必须删除
	for _, item := range *(matched_subscribers) {
		selectcase := reflect.SelectCase{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf((*item).NetChanFrame),
			Send: reflect.ValueOf(send_value),
		}
		*cases = append(*cases, selectcase)
	}
}
