package parse

import (
	"io/ioutil"

	"github.com/kataras/golog"
	"gopkg.in/yaml.v2"
)

var (
	DBConfig DB
)

type DB struct {
	Master DBConfigInfo `yaml:"Master"`
	Slave  DBConfigInfo `yaml:"Slave"`
}

type DBConfigInfo struct {
	Dialect      string `yaml:"dialect"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Database     string `yaml:"database"`
	Charset      string `yaml:"charset"`
	ShowSql      bool   `yaml:"showSql"`
	LogLevel     string `yaml:"logLevel"`
	MaxOpenConns int    `yaml:"maxOpenConns"`
	MaxIdleConns int    `yaml:"maxIdleConns"`

	//ParseTime       bool   `yaml:"parseTime"`
	//MaxIdleConns    int    `yaml:"maxIdleConns"`
	//MaxOpenConns    int    `yaml:"maxOpenConns"`
	//ConnMaxLifetime int64  `yaml:"connMaxLifetime: 10"`
	//Sslmode         string `yaml:"sslmode"`
}

func DBSettingParse() {
	golog.Info("@@@ Init db conf")
	//data, err := ioutil.ReadFile("conf/db.yml")
	////golog.Infof("%s", data)
	//if err != nil {
	//	golog.Fatalf("@@@ %s", err)
	//}
	//err = yaml.Unmarshal(data, &DBConfig)
	//if err != nil {
	//	golog.Fatalf("@@@ Unmarshal db config error!! %s", err)
	//}

	dbData, err := ioutil.ReadFile("config/db.yaml")
	if err != nil {
		golog.Fatalf("Error. %s", err)
	}
	if err = yaml.Unmarshal([]byte(dbData), &DBConfig); err != nil {
		golog.Fatalf("Error. %s", err)
	}

	//golog.Info(DBConfig)
}
