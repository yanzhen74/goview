package goviewnet

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/yanzhen74/goview/src/goviewdb"

	"github.com/Shopify/sarama"
	"github.com/yanzhen74/goview/src/model"
)

type NetKafkaGWG NetKafka

func (this *NetKafkaGWG) Init(config *model.NetWork) (int, error) {
	fmt.Printf("init kafka-gwg network")
	return init_kafka((*NetKafka)(this), config)
}

func init_kafka(this *NetKafka, config *model.NetWork) (int, error) {

	conf := sarama.NewConfig()
	conf.Consumer.Return.Errors = true
	conf.Version = sarama.V0_10_2_1
	conf.Consumer.MaxWaitTime = time.Duration(30) * time.Millisecond

	ips := strings.Split(config.NetWorkIP, ";")

	// consumer
	consumer, err := sarama.NewConsumer(ips, conf)
	if err != nil {
		fmt.Printf("consumer kafka create error %s\n", err.Error())
		return -1, err
	}

	partition_consumer, err := consumer.ConsumePartition(config.NetWorkPort, 0, sarama.OffsetNewest)
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

func (this *NetKafkaGWG) Subscribe(sub *model.FrameType) {
}

// receive net frame, parse and dispatch
func (this *NetKafkaGWG) Process() error {

	image := new(GWG_IMAGE)
	image.Image = make([]byte, 1024*1024)
	image.Size = 0

	// just test, to be deleted
	for i := 0; i < 100; i++ {
		image.Camera = (byte)(i)
		image.ImageNo = (byte)(i)
		gwg_image_save(image)
		time.Sleep(time.Second * 2)
	}

	can := make([]byte, 12)

	// image stream flag;
	// -1: no image stream; 0: new image stream;
	// 1: end of image, means new image or end of image stream;
	// 2: new image no.1; 3: new image no.2;
	pinfo := GWG_PINFO{pcan: 0, pimage: 0, image_stream_flag: 0}

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
			// process 2048 ccsds package
			gwg_image_process_ccsds(msg.Value, can, &pinfo, image)
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

type GWG_IMAGE struct {
	ImageNo byte
	Size    int
	Camera  byte //0x11:A1,0x22:A2;0x33:B
	TimeSec uint32
	TimeMs  uint16
	Image   []byte
}

type GWG_PINFO struct {
	pcan              int
	pimage            int
	image_stream_flag int
}

func gwg_image_process_ccsds(buffer []byte, can []byte, pinfo *GWG_PINFO, image *GWG_IMAGE) error {
	s_flag := gwg_image_process_2048(buffer)
	// 2000B data
	data := buffer[2048-2000-6 : 2048-6]
	pdata := 0

	if s_flag == 0x1 { //begin
		pinfo.image_stream_flag = 0

		// first frame
		x := copy(can, data)
		pinfo.pcan = 0
		pdata += x

		// image = new(GWG_IMAGE)
		if !set_gwg_first_can(can, image) {
			return errors.New("set first can error")
		}
		pinfo.pimage = 0
		pinfo.image_stream_flag = 2

		// second frame
		x = copy(can, data[pdata:])
		pinfo.pcan = 0
		pdata += x

		if !set_gwg_second_can(can, image) {
			return errors.New("set first can error")
		}

		pinfo.image_stream_flag = 3
	}

	if pinfo.image_stream_flag == -1 {
		return errors.New("set first can error")
	}

	for {
		x := copy(can[pinfo.pcan:], data[pdata:])
		pinfo.pcan += x
		pdata += x
		if pinfo.pcan == 12 {
			pinfo.pcan = 0
			info := get_can_info(can)
			if info.end_flag == 0x5 { // image stream end
				// s_flag must be 0x2
				// pinfo.image_stream_flag must be 1
				pinfo.image_stream_flag = -1
				// can data 8Byte filled by 0xAB
				// last of this 2048 data is all fill stuff, should be 0xCC
				break
			}
			if pinfo.image_stream_flag == 1 { // first frame

				//image = new(GWG_IMAGE)
				if set_gwg_first_can(can, image) {
					pinfo.pimage = 0
					pinfo.image_stream_flag = 2
				}
				continue
			}
			if pinfo.image_stream_flag == 2 { // second frame
				if set_gwg_second_can(can, image) {
					pinfo.image_stream_flag = 3
				}
				continue
			}

			if pinfo.pimage+8 < image.Size {
				y := copy(image.Image[pinfo.pimage:], can[4:])
				pinfo.pimage += y
			} else if pinfo.pimage < image.Size { // last 2048 frame
				y := copy(image.Image[pinfo.pimage:], can[4:4+image.Size-pinfo.pimage])
				pinfo.pimage += y
				// can data 8Byte last filled by 0xAA
			} else if pinfo.pimage == image.Size { // last this image's frame, I don't care
				gwg_image_save(image)
				// can data 8Byte filled by 0xAA
				pinfo.image_stream_flag = 1 // new image maybe; or image stream end
			} else {
				return errors.New("image size too long")
			}

		}
		if pdata == 2000 {
			break
		}
	}
	return nil
}

func set_gwg_first_can(can []byte, image *GWG_IMAGE) bool {
	image.Image = make([]byte, 1024*1024)
	info := get_can_info(can)
	image.ImageNo = info.image_no

	image.Size = int(can[4]) + int(can[5])*256 + int(can[6])*256*256 + int(can[7])*256*256*256
	image.Camera = can[8]
	if image.Size < 3000 { // something wrong
		fmt.Println("image.Size too small to be a valid image")
		return false
	}
	return true
}

func set_gwg_second_can(can []byte, image *GWG_IMAGE) bool {
	image.TimeSec = uint32(can[4]) + uint32(can[5])*256 + uint32(can[6])*256*256 + uint32(can[7])*256*256*256
	image.TimeMs = uint16(can[8]) + uint16(can[9])*256
	return true
}

func gwg_image_save(image *GWG_IMAGE) {
	name := fmt.Sprintf("data/gwg/%08d_%08d.jpeg", image.ImageNo, image.Camera)
	file, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, os.ModePerm)

	if err != nil {
		defer file.Close()
		return
	}

	file.Write(image.Image[:image.Size])

	file.Close()

	// write to database
	goviewdb.GwgDb.SavePic(name)

}

func gwg_image_process_2048(frame []byte) (s_flag byte) {
	p := 6
	p += 10
	p += 16
	p += 2

	// 2b子类型应用数据包序列标识，01-begin，10-end，00-mid，11-no
	s_flag = frame[p] >> 6

	return s_flag

}

type caninfo struct {
	image_no byte
	end_flag byte
	str_rtr  byte
	ide      byte
	rtr      byte
	frame_no int32
}

func get_can_info(can []byte) (info caninfo) {
	info = caninfo{}
	// 12B can
	fmt.Println(can)
	// 1B image no
	info.image_no = can[0]
	// 3b end flag when == 101b
	info.end_flag = can[1] >> 5
	// 1b STR/RTR
	info.str_rtr = (can[1] >> 4) & 0x01
	// 1b IDE
	info.ide = (can[1] >> 3) & 0x01
	// 18b frame no
	var frame_no int32 = int32(can[1] & 0x07)
	frame_no >>= 8
	frame_no += int32(can[2])
	frame_no >>= 8
	frame_no += int32(can[3] & 0xFE)
	info.frame_no = frame_no

	// 1b RTR
	info.rtr = (can[3] & 0x01)

	// 8B image
	can[0] = 0

	return info
}
