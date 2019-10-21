package test

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/sayems/golang.webdriver/selenium/pages"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tebeka/selenium"
)

var browser selenium.WebDriver
var page pages.Page

func Test_can_start_a_table_and_see_it_later(t *testing.T) {

	var err error

	// set browser as chrome
	caps := selenium.Capabilities(map[string]interface{}{"browserName": "chrome"})

	// remote to selenium server
	if browser, err = selenium.NewRemote(caps, ""); err != nil {
		fmt.Printf("Failed to open session: %s\n", err)
		return
	}

	Convey("Test should start a table and see it later", t, func() {

		Convey(` Edith want to look the telemetry data from spacecraft,
			she has heard about a cool new online goview app.
			She goes to check out its homepage `, func() {
			err = browser.Get("http://127.0.0.1:8080/")
			So(err == nil, ShouldBeTrue)
		})

		// use page to interact with browser
		page = pages.Page{Driver: browser}

		Convey("She noticed the title mention 'goview'", func() {
			title, _ := browser.Title()
			So(strings.Contains(title, "goview"), ShouldBeTrue)
		})

		Convey("She is invited to select a data table from a tree to view", func() {
			tree := page.FindElementByLinkText("WYG")
			tree.Click()
			item, err := tree.FindElement("xpath", "//span[contains(text(),'PK-CEH2.xml')]")
			item.Click()
			So(err == nil, ShouldBeTrue)
		})

		Convey("She selects the PK-CEH2.xml table", func() {
			table := page.FindElementByXpath("//div[@class='layui-tab' and @lay-filter='param-tab']")
			So(table, ShouldNotBeNil)
			name, _ := table.Text()
			So(name, ShouldContainSubstring, "PK-CEH2.xml")
		})

		Convey("When she hits enter, the page updates, and now the page shows a table named gcyctd", func() {})

		Convey("Just this time, gcyctd data is sent to the table, she sees the data varing in 2 fps", func() {
			p := simu_init_kafka()
			defer p.Close()

			table := page.FindElementByXpath("//div[@class='layui-tab' and @lay-filter='param-tab']")
			frame, _ := table.FindElement("xpath", "//iframe[contains(@src, 'PK-CEH2.xml')]")
			So(frame, ShouldNotBeNil)

			// total 10 seconds
			for i := 0; i < 50; i++ {
				simu_send_kafka(p, i)
				time.Sleep(time.Millisecond * 200)
			}
		})

		Convey(`Edith wonders whether the site will remember her table.Then she sees that the site has generated
			a unique URL for her -- there is some explanatory text to that effect. `, func() {})

		Convey("She visits that URL - her table is still there", func() {})

		Convey("Satisfied, she goes back to sleep", func() {})
	})

	browser.Quit()
}

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
