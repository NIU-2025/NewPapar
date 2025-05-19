package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

// 配置用户表
type User struct {
	Id   int
	Name string
	Pwd  string

	Articles []*Article `orm:"reverse(many)"`
}

//用户与文章是多对多的关系

// 文章信息表
type Article struct {
	Id       int       `orm:"pk;auto"`
	ArtiName string    `orm:"size(20)"`
	ArtiTime time.Time `orm:"auto_now"`
	Acount   int       `orm:"default(0);null"`
	Acontent string    `orm:"size(500)"`
	Aimg     string    `orm:"size(100)"`

	ArticleType *ArticleType `orm:"rel(fk)"` //外健
	Users       []*User      `orm:"rel(m2m)"`
}

//文章信息表和文章类型表是  多对一关系

// 文章类型表
type ArticleType struct {
	Id       int    `orm:"pk;auto"`
	TypeName string `orm:"size(20)"`

	Articles []*Article `orm:"reverse(many)"` //代表有类型有许多的文章与其对应
}

// 连接使用数据库
func init() {
	//ORM操作数据库
	//1.连接数据库
	//不知道为什么使用root用户无法登陆，于是我使用的是/etc/mysql中debian.cnf文件中的用户名和密码
	orm.RegisterDataBase("default", "mysql", "debian-sys-maint:4pcyN6SocSRJLCd6@tcp(127.0.0.1:3306)/niu")

	//2.创建表
	orm.RegisterModel(new(User), new(Article), new(ArticleType)) //注册表
	orm.RunSyncdb("default", false, true)                        //生成表
}
