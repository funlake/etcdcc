package controllers

type HomeController struct {
	BaseController
}
func (h *HomeController) Home(){
	h.response(RESPOK,"hello world",nil)
}