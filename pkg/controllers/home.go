package controllers

//HomeController : default controller
type HomeController struct {
	BaseController
}

//Home : default home page
func (h *HomeController) Home() {
	h.response(RESPOK, "hello world", nil)
}
