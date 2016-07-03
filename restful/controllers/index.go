package controllers

type IndexController  struct{
	BaseController

}

func (this *IndexController) Index()  {
	this.TplName = "index.html"
	this.Render()
}
