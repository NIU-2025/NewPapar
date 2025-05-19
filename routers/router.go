package routers

import (
	"github.com/astaxie/beego"
	"testone/controllers"
)

func init() {
	//注册页面路由
	beego.Router("/register", &controllers.UserController{}, "get:Register;post:HandleRegister")

	//登陆页面路由
	beego.Router("/login", &controllers.UserController{}, "get:LoginArticle;post:HandleLogin")

	//展示网站主页
	beego.Router("/showArticleList", &controllers.ArticleController{}, "get:ShowArticleList")

	//添加文章
	beego.Router("/addArticle", &controllers.ArticleController{}, "get:ShowAddArticle;post:HandleAddArticle")

	//文章详情页面
	beego.Router("/showArticleDetail", &controllers.ArticleController{}, "get:ShowArticleDetail")

	//添加文章分类
	beego.Router("/addType", &controllers.ArticleController{}, "get:ShowAddType;post:HandleAddType")

	//删除文章
	beego.Router("/deleteArticle", &controllers.ArticleController{}, "get:DeleteArticle")

	//更新文章数据
	beego.Router("/updataArticle", &controllers.ArticleController{}, "get:UpdataArticle;post:HandleUpdataArticle")

	//退出网站页面
	beego.Router("/logout", &controllers.ArticleController{}, "get:Logout")
}
