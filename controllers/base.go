package controllers

//base controller for auth

import (
	"SuperCenterServer/constants"
	"beego-mgo-rabbitmq-jwt/models"
	"beego-mgo-rabbitmq-jwt/utilities/mgodb"
	services "beego-mgo-rabbitmq-jwt/utilities/svc"
	"fmt"
	"time"

	"beego-mgo-rabbitmq-jwt/utilities/helper"

	"github.com/astaxie/beego"
	jwt "github.com/dgrijalva/jwt-go"
	log "github.com/goinggo/tracelog"
)

//BaseController base controller
type BaseController struct {
	beego.Controller
	services.Service        //every controller has a service
	SecureKey        string //do something special
	ControllerName   string
	ActionName       string
}

type Page struct {
	Content             interface{} `json:"content"` //page content data
	IsFirstPage         bool        `json:"isFirst"`
	IsLastPage          bool        `json:"isLast"`
	PageNum             int         `json:"pageNum"` //current page number
	PerSize             int         `json:"perSize"`
	PageCount           int         `json:"pageCount"`
	ElementCount        int         `json:"elementCount"`
	CurPageElementCount int         `json:"curPageElementCount"`
}

// PageUtil create Page Options
//
// @count content count
// @pageNo page number
// @pageSize page PerSize
// @list the list of content
func PageUtil(count int64, pageNo int, pageSize int, list interface{}) Page {
	tp := count / int64(pageSize)
	fmt.Println("count:", count, "pageNo:", pageNo, "pageSize:", pageSize)
	if count%int64(pageSize) > 0 {
		tp = count/int64(pageSize) + 1
	}
	return Page{PageNum: pageNo, PerSize: pageSize, PageCount: int(tp), ElementCount: int(count), IsFirstPage: pageNo == 1, IsLastPage: pageNo == int(tp), Content: list}
}

//Prepare execute when a controller create
func (base *BaseController) Prepare() {
	base.ControllerName, base.ActionName = base.GetControllerAndAction()
	sec := base.Ctx.Input.Header("Secure") //custom header
	useSecureKey := false                  //if use custom header value
	base.SecureKey = sec
	//read auth infomation
	tokenString := base.Ctx.Input.Header("Authorization")
	fmt.Println(tokenString)
	if base.SecureKey == string(constants.SecureKey[:]) {
		base.UserID = "test"
		base.UserName = "test"
		//get db connection | just check a db connection
		if err := base.Service.Prepare(); err != nil {
			log.Errorf(err, base.UserID, "Service.Prepare.SecureKey", base.Ctx.Request.URL.Path)
			base.ServeError(err)
		}
		base.Service.UserID = base.UserID
		base.Service.UserName = base.UserName
		useSecureKey = true
	}
	if useSecureKey {
		return
	}
	token, err := jwt.ParseWithClaims(tokenString, &models.MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return constants.Secret, nil
	})

	fmt.Println(token.Valid)

	if err != nil {
		fmt.Println("somthing error", err)
	}
	if token.Valid {
		fmt.Println("you looks good today!")
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("token validation error")
			base.ServeAuthErrorWithMsg(err, "validation error : token parse Malformed error")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			log.CompletedError(err, "validation error : token expired", "base.Prepare")
			base.ServeAuthErrorWithMsg(err, "validation error : token expired")
		} else {
			log.CompletedError(err, "validation error : token parse error", "base.Prepare")
			base.ServeAuthErrorWithMsg(err, "validation error : token parse error")
		}
	} else {
		log.CompletedError(err, "handle token error", "base.Prepare")
		base.ServeAuthErrorWithMsg(err, "handle token error")

	}

	if claims, ok := token.Claims.(*models.MyClaims); ok && token.Valid {
		base.UserID = claims.DevID
		base.UserName = claims.Name
		if base.UserID == "" {
			base.UserID = "unknown"
		}

		if err := base.Service.Prepare(); err != nil {
			log.Errorf(err, base.UserID, "BaseController.Prepare", base.Ctx.Request.URL.Path)
			base.ServeError(err)
			return
		}
		base.Service.UserID = base.UserID
		base.Service.UserName = base.UserName
	} else {
		fmt.Println(err)
	}

}

// Finish realease connection
func (base *BaseController) Finish() {
	defer func() {
		if base.MongoSession != nil {
			mgodb.CloseSession(base.UserID, base.MongoSession)
			base.MongoSession = nil
		}
	}()

	log.Completedf(base.UserID, "Finish", base.Ctx.Request.URL.Path)
}

// ServeError prepares and serves an Error exception.
func (base *BaseController) ServeError(err error) {
	base.Data["json"] = models.CustomException{
		Error:     "ServerError",
		Exception: err.Error(),
		Message:   "server err when handle data",
		Path:      base.Ctx.Input.Method() + " " + base.Ctx.Input.URL(),
		TimeStamp: time.Now()}

	base.Ctx.Output.SetStatus(500) //set internal server error
	base.ServeJSON()
}

// ServeErrorWithMsg Serve an error with a readable message
func (base *BaseController) ServeErrorWithMsg(err error, msg string) {
	base.Data["json"] = models.CustomException{
		Error:     "ServerError",
		Exception: err.Error(),
		Message:   msg,
		Path:      base.Ctx.Input.Method() + " " + base.Ctx.Input.URL(),
		TimeStamp: time.Now()}

	base.Ctx.Output.SetStatus(500) //set internal server error
	base.ServeJSON()
}

//ServeAuthError Seve Auth error message
func (base *BaseController) ServeAuthError(err error) {
	base.Data["json"] = models.CustomException{
		Error:     "AuthError",
		Exception: err.Error(),
		Message:   base.Ctx.Input.Header("Authorization"),
		Path:      base.Ctx.Input.Method() + " " + base.Ctx.Input.URL(),
		TimeStamp: time.Now()}
	base.Ctx.Output.SetStatus(401)
	base.ServeJSON()
}

//ServeAuthErrorWithMsg Seve Auth error with readable message
func (base *BaseController) ServeAuthErrorWithMsg(err error, msg string) {
	base.Data["json"] = models.CustomException{
		Error:     "AuthError",
		Exception: err.Error(),
		Message:   msg,
		Path:      base.Ctx.Input.Method() + " " + base.Ctx.Input.URL(),
		TimeStamp: time.Now()}
	base.Ctx.Output.SetStatus(401)
	base.ServeJSON()
}

//ServeSuccessWithMsg return some message when request successed with some message
func (base *BaseController) ServeSuccessWithMsg(msg string) {
	base.Data["json"] = models.CustomException{
		Error:     "",
		Exception: "",
		Message:   msg,
		Path:      base.Ctx.Input.Method() + " " + base.Ctx.Input.URL(),
		TimeStamp: time.Now()}
	base.Ctx.Output.SetStatus(200) //request success
	base.ServeJSON()
}

//HandleResult return data also can handle errors
// no need to check error in controller
func (base *BaseController) HandleResult(result interface{}, err error) {
	if err == nil {
		base.Data["json"] = &result
		base.ServeJSON()
	} else {
		base.ServeError(err)
	}
}

//HandleResultWithMsg return data also can handle errors with readable message
// no need to check error in controller
func (base *BaseController) HandleResultWithMsg(result interface{}, err error, msg string) {
	if err == nil {
		base.Data["json"] = &result
		base.ServeJSON()
	} else {
		log.CompletedError(err, msg, base.ControllerName+"."+base.ActionName)
		base.ServeErrorWithMsg(err, msg)
	}
}

//GetPageInfo get page info
func (base *BaseController) GetPageInfo() *helper.Paginator {
	return base.Data["paginator"].(*helper.Paginator)
}

//HandlePageResult return paged results
func (base *BaseController) HandlePageResult(result interface{}, err error, p *helper.Paginator) {
	if err == nil {

		page := PageUtil(p.Nums(), p.Page(), p.PerPageNums, result)

		base.Data["json"] = &page
		fmt.Println("paginator:", p)
		base.ServeJSON()
	} else {
		base.ServeError(err)
	}
}

//SetPaginator ...
func (base *BaseController) SetPaginator(per int, nums int) *helper.Paginator {
	p := helper.NewPaginator(base.Ctx.Request, per, nums)
	base.Data["paginator"] = p
	return p
}
