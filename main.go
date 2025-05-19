package main

import (
	"github.com/astaxie/beego"
	_ "testone/models"
	_ "testone/routers"
)

func main() {
	beego.AddFuncMap("prepage", ShowPrePage)
	beego.AddFuncMap("nextpage", ShowNextPage)
	beego.Run()
}

// 视图函数  上一页的
func ShowPrePage(pageIndex int) int {
	if pageIndex == 1 {
		return pageIndex
	}
	return pageIndex - 1
}

// 视图函数  下一页的
func ShowNextPage(pageIndex, pageCount int) int {
	if pageIndex == pageCount {
		return pageCount
	}
	return pageIndex + 1
}
