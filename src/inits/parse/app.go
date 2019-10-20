package parse

import (
	// "go-iris/main/inits/bindata/conf"

	"io/ioutil"

	"github.com/kataras/golog"
	"gopkg.in/yaml.v2"
)

var (
	// 解析app.yml中的Other项
	O Other
)

type Other struct {
	Other Inner `yaml:"Other"`
}
type Inner struct {
	Port       string   `yaml:"Port"`
	IgnoreURLs []string `yaml:"IgnoreURLs"`
	JWTTimeout int64    `yaml:"JWTTimeout"`
	LogLevel   string   `yaml:"LogLevel"`
	Secret     string   `yaml:"Secret"`
}

func AppOtherParse() {
	golog.Info("@@@ Init app conf")
	//c := iris.YAML("conf/app.yml")

	// conf := new(module.Yaml)
	appData, err := ioutil.ReadFile("config/app.yaml")

	// appData, err := conf.Asset("conf/app.yml")
	if err != nil {
		golog.Fatalf("Error. %s", err)
	}
	if err = yaml.Unmarshal([]byte(appData), &O); err != nil {
		golog.Fatalf("Error. %s", err)
	}

}
