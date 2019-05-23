package controllers

//Home api page
type HomeController struct {
	BaseController
}

//Home api
func (h *HomeController) Home() {
	h.response(RESPOK, "hello world", nil)
}
