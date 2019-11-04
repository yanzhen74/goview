package functional

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Shopify/sarama"
)

func simu_init_kafka() (p sarama.SyncProducer) {

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second
	p, err := sarama.NewSyncProducer([]string{"10.211.55.2:9092"}, config)
	if err != nil {
		log.Printf("sarama.NewSyncProducer err, message=%s \n", err)
		return nil
	}

	return p
}

func simu_send_kafka(p sarama.SyncProducer, i int) {
	topic := "RTM"
	srcValue0 := "RTM_XJYH_PK-CEH2_Result\t.\tindex=%d\n1 aa 233;2 bb 55;3 00000000 53.78;4 39a8 %d;5 55aa %d"
	srcValue1 := "RTM_WYG_PK-CEH2_Result\t.\tindex=%d\n1 11 233;2 22 55;3 33 53.78;4 22cc %d;5 ffee %d"
	var value string
	if i%2 == 0 {
		value = fmt.Sprintf(srcValue0, i, i, i)
	} else {
		value = fmt.Sprintf(srcValue1, i, i, i)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}
	part, offset, err := p.SendMessage(msg)
	if err != nil {
		log.Printf("send message(%s) err=%s \n", value, err)
	} else {
		fmt.Fprintf(os.Stdout, value+"发送成功，partition=%d, offset=%d \n", part, offset)
	}
}
