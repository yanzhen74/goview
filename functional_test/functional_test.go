package test

import (
	"flag"
	"fmt"
	"strings"
	"testing"

	"github.com/sayems/golang.webdriver/selenium/pages"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tebeka/selenium"
)

var driver selenium.WebDriver
var page pages.Page

func TestMain(t *testing.T) {

	var err error

	// set browser as chrome
	caps := selenium.Capabilities(map[string]interface{}{"browserName": "chrome"})
	const (
		seleniumPath = `/usr/local/bin/chromedriver`
		port         = 9515
	)

	// remote to selenium servrr
	if driver, err = selenium.NewRemote(caps, ""); err != nil {
		fmt.Printf("Failed to open session: %s\n", err)
		return
	}

	page = pages.Page{Driver: driver}

	err = driver.Get("http://127.0.0.1:8080/")
	if err != nil {
		fmt.Printf("Failed to load page: %s\n", err)
		return
	}
	Convey("TestMain should open a web page with title goview", t, func() {
		title, _ := driver.Title()
		So(strings.Contains(title, "goview"), ShouldBeTrue)
	})

	fmt.Println(driver.PageSource())

	flag.Parse()
	//r := m.Run()

	driver.Quit()
	//os.Exit(r)

}
