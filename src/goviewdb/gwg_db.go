package goviewdb

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

// 数据库对象DB,对xorm进行了包装 //注意驼峰法
type Db struct {
	ORM *xorm.Engine
}

var gwg_lock sync.Mutex
var GwgDb *Db

// 数据库表
// GWGPic 高微柜图像信息表
type GWGPic struct {
	Id        int64  // xorm默认自动递增
	Url       string //图片保存路径
	Camera    uint8
	Size      uint32
	ImageNo   uint8
	Time      string    // 船上时
	CreatedAt time.Time `xorm:"created"`
}

//生成数据库
func NewGWGDb(dbfile string) *Db {
	// 初始化数据库
	orm, _ := xorm.NewEngine("sqlite3", dbfile)
	//此时生成数据库文件db/gwgpic.db，以及表GWGPIC
	var err error
	if err = orm.Sync2(new(GWGPic)); err != nil {
		log.Fatalf("Fail to sync database %v\n", err)
	}
	GwgDb = new(Db)
	GwgDb.ORM = orm
	return GwgDb
}

// 保存图片
func (d *Db) SavePic(url string) {
	pic := new(GWGPic)
	pic.Url = url
	gwg_lock.Lock()
	defer gwg_lock.Unlock()
	_, err := d.ORM.Insert(pic)
	if err != nil {
		log.Fatalf("Fail to insert database %v\n", err)
	}
}

// 列出图片
func (d *Db) ListPic(limit int, pageno int, codeLike string) []GWGPic {
	var p []GWGPic
	var query string
	if codeLike == "" {
		query = fmt.Sprintf("select * from g_w_g_pic LIMIT %v OFFSET %v", limit, (pageno-1)*limit)
	} else {
		query = fmt.Sprintf("select * from g_w_g_pic where %v LIMIT %v OFFSET %v", codeLike, limit, (pageno-1)*limit)
	}
	gwg_lock.Lock()
	err := d.ORM.Sql(query).Find(&p)
	gwg_lock.Unlock()
	if err != nil {
		log.Fatalf("Fail to list database %v\n", err)
	}

	return p
}

// 得到总数
func (d *Db) CountPic() int64 {
	gwg_lock.Lock()
	total, err := d.ORM.Count(new(GWGPic))
	gwg_lock.Unlock()

	if err != nil {
		log.Fatalf("Fail to count database %v\n", err)
	}
	return total
}
