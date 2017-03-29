// @APIVersion 1.0.0
// @Title 宇宙超级无敌中心服务器
// @Description 既然你诚心诚意地发问了 那我就大发慈悲地告诉你
// @Contact linx@llinx.me
// @TermsOfServiceUrl http://llinx.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"beego-mgo-rabbitmq-jwt/controllers"

	"github.com/astaxie/beego"
)

func init() {
	//only write router like this can swagger read the api
	ns :=
		beego.NewNamespace("/v1",
			beego.NSNamespace("/fake",
				beego.NSInclude(
					&controllers.FakeController{},
				),
			),
		)
	beego.AddNamespace(ns)

}
