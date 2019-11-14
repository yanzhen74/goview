package test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestBinaryFile(t *testing.T) {
	file, err := os.OpenFile("test", os.O_CREATE|os.O_RDWR, os.ModePerm)

	if err != nil {
		defer file.Close()
		return
	}

	// write 2048 bytes
	var frame []byte = make([]byte, 2048)
	for i := 0; i < 2048; i++ {
		frame[i] = (byte)(i & 0xFF)
	}

	fmt.Println(frame)

	file.Write(frame)

	file.Close()

}

func TestReadBFile(t *testing.T) {
	file, err := os.OpenFile("test", os.O_RDONLY, os.ModePerm)

	if err != nil {
		defer file.Close()
		return
	}

	// read 2048 bytes
	var buffer bytes.Buffer
	io.CopyN(&buffer, file, 2048)
	_bytes := buffer.Bytes()

	fmt.Println(_bytes)

}

func TestBitOp(t *testing.T) {
	a := []byte{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	b := []byte{255, 254, 253, 252, 251}
	fmt.Println(int(b[0]) * 255)
	copy(a[:3], b)
	fmt.Println(a)
	b[0] = 33
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(a[:10])
	n := copy(a[8:], b)
	fmt.Println(n, a)

	var frame []byte = make([]byte, 2048)
	for i := 0; i < 2048; i++ {
		frame[i] = (byte)(i & 0xFF)
	}

	gwg_image_write_2048(frame)
}

func TestArrayPara(t *testing.T) {
	can := make([]byte, 12)
	fill_can(can)
	fmt.Println(can)
}
func fill_can(can []byte) {
	for i := 0; i < len(can); i++ {
		can[i] = 0x22
	}
}

type GWG_PINFO struct {
	pcan              int
	pimage            int
	image_stream_flag int
}

func TestStructPara(t *testing.T) {
	info := GWG_PINFO{pcan: 3, pimage: 20, image_stream_flag: 5}

	fill_info(&info)
	fmt.Println(info)
}

func fill_info(info *GWG_PINFO) {
	info.image_stream_flag = 44
	info.pcan = 66
	info.pimage = 55
}

func TestPrintf(t *testing.T) {
	fmt.Printf("%08d_%08d.jepg\n", 22, 33)
}
func TestProcessGWG(t *testing.T) {
	gwg_image_process("")
}

func TestProcessGWG2Images(t *testing.T) {
	gwg_image_process("2048-0102.dat")
}

type GWG_IMAGE struct {
	ImageNo byte
	Size    int
	Camera  byte //0x11:A1,0x22:A2;0x33:B
	TimeSec uint32
	TimeMs  uint16
	Image   []byte
}

func gwg_image_process(fileName string) (err error) {

	if "" == fileName {
		fileName = "2048.dat"
	}
	file, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)

	if err != nil {
		return
	}
	defer file.Close()

	// read 2048 bytes
	BUFFERSIZE := 2048
	buffer := make([]byte, BUFFERSIZE)

	image := new(GWG_IMAGE)
	image.Image = make([]byte, 1024*1024)
	image.Size = 0

	can := make([]byte, 12)

	// image stream flag;
	// -1: no image stream; 0: new image stream;
	// 1: end of image, means new image or end of image stream;
	// 2: new image no.1; 3: new image no.2;
	pinfo := GWG_PINFO{pcan: 0, pimage: 0, image_stream_flag: 0}

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n != BUFFERSIZE {
			break
		}

		gwg_image_process_ccsds(buffer, can, &pinfo, image)
	}

	return nil
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
	name := fmt.Sprintf("%d_%d.jpeg", image.ImageNo, image.Camera)
	file, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, os.ModePerm)

	if err != nil {
		defer file.Close()
		return
	}

	file.Write(image.Image[:image.Size])

	file.Close()
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

func gwg_image_write_2048(frame []byte) {

	// bit op
	// 6B CCSDS header
	p := 6
	// 10B CCSDS vice header
	p += 10
	// 16B 载荷源包 header
	p += 16

	// 2B 应用数据类型 c8c8
	var app_type byte = 0xc8
	frame[p] = app_type
	frame[p+1] = app_type
	p += 2

	// 2b 子类型应用数据包序列标志, 01-begin, 10-end, 00-mid, 11-no
	var s_flag byte = 0x1

	frame[p] = s_flag << 6
	fmt.Println(frame[p])

	// 14b 序列号
	var s_no uint16 = 0xEEDD
	frame[p] |= (byte)((s_no >> 8) & (0xFF >> 2))
	frame[p+1] = (byte)(s_no & 0xFF)
	fmt.Printf("%x %x\n", frame[p], 0xEE&(0xFF>>2))
	fmt.Println(frame[p] == (s_flag<<6)|0xEE&(0xFF>>2))
	fmt.Println(frame[p+1] == 0xDD)
	p += 2

	// 4B + 2B 内部时间码
	var sec uint32 = 0x00112233
	frame[p] = (byte)(sec >> 24)
	frame[p+1] = (byte)(sec >> 16)
	frame[p+2] = (byte)(sec >> 8)
	frame[p+3] = (byte)(sec & 0xFF)
	var usec uint16 = 0x0055
	frame[p+4] = (byte)(usec >> 8)
	frame[p+5] = (byte)(usec & 0xff)
	p += 6

	// 2000B data

	p += 2000

	// 2B crc
	p += 2

	// 4B crc
	p += 4
	fmt.Println(p)
}

func TestRBitOp(t *testing.T) {
	var frame []byte = make([]byte, 2000)

	for i := 0; i < 2000; i++ {
		frame[i] = (byte)(i & 0xFF)
	}

	var image []byte = make([]byte, 0, 1024*1024)
	i := 0
	for ; i < 2000/12; i++ {
		// 12B CAN数据包
		index := i * 12
		if !process_can(frame[index:index+12], &image) {
			break
		}
		fmt.Println(frame[index : index+12])
	}
	fmt.Println(2000-12*i, 8)
	// 8B 0xCC填充
	fmt.Println(frame[12*i:])
	// image
	fmt.Println(image)
}

func process_can(can []byte, image *[]byte) (over bool) {
	info := caninfo{}
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
	*image = append(*image, can[4:]...)
	can[0] = 0

	return info.end_flag != 0x5
}
