package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"path"
	"testone/models"
	"time"
)

type ArticleController struct {
	beego.Controller
}

// 展示网站主页
func (this *ArticleController) ShowArticleList() {
	//session判断是否获取到用户名
	userName := this.GetSession("username")
	if userName == nil {
		this.Redirect("/login", 302)
		return
	}

	o := orm.NewOrm()
	//遍历数据库的文章表
	qs := o.QueryTable("Article")
	var articles []models.Article
	//获取所有数据
	/*_, err := qs.All(&articles) //注意  这里传递是切片
	if err != nil {
		beego.Info("查询数据出错")
	}*/

	//分页操作
	var count int64
	//一面展示5页
	pageSize := 5
	//总记录数
	//count, _ := qs.Count()

	//上一页  下一页操作
	pageIndex, err := this.GetInt("pageIndex") //获取页码位置，首页为1
	if err != nil {
		pageIndex = 1
	}
	//获取数据
	//获取数据库中的部分数据  参数  获取几条  从那里开始获取
	start := (pageIndex - 1) * pageSize
	qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articles)

	//关联文章类型表中的数据

	//1.根据文章类型显示页数
	typeName := this.GetString("select")
	if typeName == "" {
		count, _ = qs.Count() //刚刚登陆页面，还没有选择文章类型，则统计表的总记录
	} else {
		//选中文章类型，进行过滤的展示
		count, _ = qs.Limit(pageSize, start).RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).Count()
	}
	//天花板函数 向上取整
	pageCount := math.Ceil(float64(count) / float64(pageSize)) //总页数

	//查询文章类型表中的数据
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"] = types

	//2.根据文章类型显示文章
	if typeName == "" {
		beego.Info("33333")
		qs.Limit(pageSize, start).All(&articles) //刚刚登陆页面，还没有选择文章类型,则显示所哟文章
	} else {
		qs.Limit(pageSize, start).RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).All(&articles)
	}

	//获取登陆用户的用户名
	name := this.GetSession("username")
	this.Data["name"] = name.(string)

	this.Data["typeName"] = typeName
	this.Data["count"] = count
	this.Data["pageCount"] = int(pageCount)
	this.Data["pageIndex"] = pageIndex
	this.Data["articles"] = articles
	this.TplName = "index.html"
}

// 展示添加文章页面
func (this *ArticleController) ShowAddArticle() {
	//关联文章类型表中数据
	o := orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)

	//获取登陆用户的用户名
	name := this.GetSession("username")
	this.Data["name"] = name.(string)

	this.Data["types"] = types
	this.TplName = "add.html"
}

// 上传文件的函数
func UploadFile(this *beego.Controller, filePath string) string {
	file, head, err := this.GetFile(filePath)
	/*if head.Filename == "" {
		beego.Info("未选择文件")
		this.Data["errmsg"] = "未选择文件"
		return ""
	}*/
	/*fileName := this.GetString(filePath)
	if fileName == "" {
		beego.Info("未选择文件")
		this.Data["errmsg"] = "未选择文件"
		return ""
	}*/

	if err != nil {
		this.Data["errmsg"] = "上传文件失败"
		this.TplName = "add.html"
		return ""
	}
	defer file.Close()
	//判断文件大小
	if head.Size > 5000000 {
		this.Data["errmsg"] = "文件太大"
		this.TplName = "add.html"
		return ""
	}
	//判断文件后缀是否是可选型式
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" {
		this.Data["errmsg"] = "文件格式不正确"
		this.TplName = "add.html"
		return ""
	}
	//判断上传的文件是否重名
	filePath = time.Now().Format("2006-01-02-15:06:07") + ext
	//存储  存储文件名   存储位置
	this.SaveToFile(filePath, "./static/img"+filePath)

	//返回存储位置
	return "/static/img" + filePath
}

// 处理添加文章页面数据
func (this *ArticleController) HandleAddArticle() {
	//获取数据
	articleName := this.GetString("article_title_name")
	articleContent := this.GetString("article_content")

	if articleName == "" || articleContent == "" {
		beego.Info("数据不完整")
		this.Data["errmsg"] = "数据不完整"
		this.TplName = "add.html"
	}

	//封装上传文件的函数
	filePath := UploadFile(&this.Controller, "uploadname")
	if filePath == "" {
		this.TplName = "add.html"
		return
	}

	//将获取到的添加文章数据插入到数据库中
	o := orm.NewOrm()
	var article models.Article
	article.ArtiName = articleName
	article.Acontent = articleContent
	article.Aimg = filePath

	//关联文章类型表  给文章添加类型
	//获取数据
	TypeName := this.GetString("select")
	beego.Info(TypeName)
	//添加数据
	var articleType models.ArticleType
	articleType.TypeName = TypeName
	//添加分类数据之前先查询 分类是否存在
	o.Read(&articleType, "TypeName")
	article.ArticleType = &articleType

	o.Insert(&article) //插入数据库

	//返回页面
	this.Redirect("/showArticleList", 302)
}

// 展示文章详情页面
func (this *ArticleController) ShowArticleDetail() {
	//获取数据
	id, err := this.GetInt("articleId")
	if err != nil {
		beego.Info("传递的链接错误")
	}

	//遍历数据库 找到特定的数据
	//将文章信息表中的文章类型和文章类型表关联起来
	o := orm.NewOrm()
	var article models.Article
	o.QueryTable("Article").RelatedSel("ArticleType").Filter("Id", id).One(&article)

	//增加阅读次数
	article.Acount += 1
	o.Update(&article)

	//关联用户与文章关系，也就是最近浏览显示真实用户
	m2m := o.QueryM2M(&article, "Users")    //多对多查询 需要插入的表  表中那个字段需要插入
	userName := this.GetSession("username") //获取用户名
	if userName == nil {
		this.Redirect("/showArticleList", 302)
		return
	}
	var user models.User
	user.Name = userName.(string)
	//插入之前 先查询数据库中是否存在用户
	o.Read(&user, "Name")

	//插入
	m2m.Add(user)

	//插入之后 查询
	var users []models.User
	o.QueryTable("User").Filter("Articles__Article__Id", id).Distinct().All(&users)
	this.Data["users"] = users

	//获取登陆用户的用户名
	name := this.GetSession("username")
	this.Data["name"] = name.(string)

	this.Data["article"] = article
	this.TplName = "article_content.html"
}

// 展示文章分类页面
func (this *ArticleController) ShowAddType() {
	//查询文章类型表中的数据
	o := orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)

	//获取登陆用户的用户名
	name := this.GetSession("username")
	this.Data["name"] = name.(string)

	this.Data["types"] = types
	this.TplName = "addType.html"
}

// 处理添加文章分类数据
func (this *ArticleController) HandleAddType() {
	//获取数据
	typeName := this.GetString("typeName")
	if typeName == "" {
		this.Data["nonull"] = "添加分类不能为空"
		beego.Info("添加分类不能为空")
		this.TplName = "addType.html"
		return
	}

	//处理数据，也就是将文章分类插入到数据库中
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Insert(&articleType)

	//返回页面
	this.Redirect("/addType", 302)
}

// 删除文章
func (this *ArticleController) DeleteArticle() {
	id, err := this.GetInt("articleId")
	if err != nil {
		beego.Info("删除页面错误")
		this.TplName = "index.html"
		return
	}

	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.Delete(&article)

	this.Redirect("/showArticleList", 302)
}

// 展示更新文章页面
func (this *ArticleController) UpdataArticle() {
	//获取数据
	id, err := this.GetInt("articleId")
	if err != nil {
		beego.Info("传递的链接错误")
	}

	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.Read(&article)

	//需要后续改
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"] = types

	//获取登陆用户的用户名
	name := this.GetSession("username")
	this.Data["name"] = name.(string)

	this.Data["article"] = article
	this.TplName = "updata.html"
}

// 处理更新文章数据
func (this *ArticleController) HandleUpdataArticle() {
	//获取数据
	id, err := this.GetInt("articleId")
	articleName := this.GetString("article_title_name")
	articleContent := this.GetString("article_content")
	filePath := UploadFile(&this.Controller, "uploadname")

	//校验数据
	if articleName == "" && articleContent == "" && filePath == "" {
		beego.Info("更新失败")
		this.TplName = "updata.html"
		return
	}
	/*o := orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"] = types*/

	//处理数据
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	//更新之前需要查询数据库中是否存在该数据
	err = o.Read(&article)
	if err != nil {
		beego.Info("数据不存在")
		return
	}

	//更新
	article.ArtiName = articleName
	article.Acontent = articleContent
	article.Aimg = filePath
	o.Update(&article)

	//返回页面
	this.Redirect("/showArticleList", 302)
}

// 退出网站页面
func (this *ArticleController) Logout() {
	//获取session
	userName := this.GetString("username")

	//退出
	this.DelSession(userName)

	this.Redirect("/login", 302)
}
