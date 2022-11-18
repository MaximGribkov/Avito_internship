package user

import (
	_ "avito/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Пути для запросов

const (
	test            = "/test"              // GET
	balanceAdd      = "/balance/add"       //POST
	balanceGet      = "/balance/get"       //GET
	balanceTransfer = "/balance/transfer"  //POST
	balanceList     = "/balance/list"      //GET
	reserveMake     = "/reserve/make"      //POST
	reserveAccept   = "/reserve/accept"    //POST
	reportCSV       = "/reserve/csv/:file" //GET
	report          = "/report"            //GET
)

func Router(router *gin.Engine) {
	router.GET(test, Welcome)
	router.GET(balanceGet, GetBal)
	router.GET(balanceList, ReportHist)
	router.GET(report, ReportSer)
	router.GET(reportCSV, OpenReport)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST(balanceAdd, AddBal)
	router.POST(balanceTransfer, TransBal)
	router.POST(reserveMake, ReserveBal)
	router.POST(reserveAccept, ConfBal)
}
