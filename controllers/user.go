package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"testone/models"
)

type UserController struct {
	beego.Controller
}

// 展示登陆页面
func (this *UserController) LoginArticle() {
	userName := this.Ctx.GetCookie("username")
	if userName == "" {
		this.Data["checked"] = ""
		this.Data["userName"] = ""
	} else {
		this.Data["checked"] = "checked"
		this.Data["userName"] = userName
	}
	this.TplName = "login.html"
}

// 处理登陆页面数据
func (this *UserController) HandleLogin() {
	//获取数据
	userName := this.GetString("username")
	passWord := this.GetString("password")

	//判断数据
	if userName == "" || passWord == "" {
		beego.Info("数据不完整")
		this.TplName = "login.html"
		return
	}

	//处理数据
	o := orm.NewOrm()
	var user models.User
	user.Name = userName

	//判断用户名是否存在并且密码是否正确
	err := o.Read(&user, "Name")
	if err != nil {
		this.Data["err_login"] = "用户名不存在"
		this.TplName = "login.html"
		return
	}
	if user.Pwd != passWord {
		this.Data["err_login"] = "密码不正确"
		this.TplName = "login.html"
		return
	}

	//记住用户名操作
	remember := this.GetString("remember")

	if remember == "on" {
		this.Ctx.SetCookie("username", userName, 100)
	} else {
		this.Ctx.SetCookie("username", userName, -1)
	}

	this.SetSession("username", userName)
	//返回页面
	this.Redirect("/showArticleList", 302)
}

// 展示注册页面
func (this *UserController) Register() {
	this.TplName = "register.html"
}

// 处理注册页面数据
func (this *UserController) HandleRegister() {
	//获取数据
	userName := this.GetString("username")
	passWord := this.GetString("password")

	if userName == "" || passWord == "" {
		beego.Info("数据不完整")
		this.TplName = "register.html"
		return
	}

	//处理数据
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	user.Pwd = passWord

	o.Insert(&user)

	//返回页面
	this.Redirect("/login", 302)
}
