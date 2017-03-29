package controllers

//fake controller for testing

import (
	"github.com/astaxie/beego"
)

//FakeController not a BaseController
type FakeController struct {
	beego.Controller
}

// GenFakeKey generate a JWT Token
func (s *FakeController) GenFakeKey() {

}
