package user

// Тут описаны все хендлеры для запросов

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

// @Summary Test
// @Tags get
// @Description test function
// @Accept json
// @Produce json
// @Success 200 {integer} integer
// @Failure 400 {object} errors.ResponseError
// @Router /test [get]
func Welcome(h *gin.Context) {
	// Приветсвенная функция
	var reqBody BalGetModel
	if err := h.BindJSON(&reqBody); err != nil {
		h.JSON(http.StatusBadRequest, err)
		logrus.Infof("some error %s", err)
	}
	ok := Hello()
	h.JSON(http.StatusOK, ok)
}

// @Summary GetBalance
// @Tags get
// @Description get balance user
// @Accept json
// @Produce json
// @Param input body BalGetModel true "user_id"
// @Success 200 {integer} integer
// @Failure 400 {object} errors.ResponseError
// @Router /balance/get [get]
func GetBal(h *gin.Context) {
	// Получение баланса определенного пользователя
	var reqBody BalGetModel
	if err := h.BindJSON(&reqBody); err != nil {
		h.JSON(http.StatusBadRequest, err)
		logrus.Infof("some error %s", err)
	}
	ok := BalanceGet(reqBody)
	h.JSON(http.StatusOK, ok)

}

// @Summary addBalanceUser
// @Tags post
// @Description replenishment of the user's account
// @Accept json
// @Produce json
// @Param input body BalAddModel true "user_id"
// @Success 204 {integer} integer
// @Failure 400 {object} errors.ResponseError
// @Router /balance/add [post]
func AddBal(h *gin.Context) {
	// Зачисление средств на баланс
	var reqBody BalAddModel
	if err := h.BindJSON(&reqBody); err != nil {
		h.JSON(http.StatusBadRequest, err)
		logrus.Infof("some error %s", err)
	}
	ok := AddBalance(reqBody)
	h.JSON(http.StatusOK, ok)
}

// @Summary transactionBalance
// @Tags post
// @Description transfer from user to user
// @Accept json
// @Produce json
// @Param input body TransferBalanceModel true "from"
// @Success 204 {integer} integer
// @Failure 400 {object} errors.ResponseError
// @Router /balance/transfer [post]
func TransBal(h *gin.Context) {
	// Перевод от одного пользователя к другому при условии наличия средств
	var reqBody TransferBalanceModel
	if err := h.BindJSON(&reqBody); err != nil {
		h.JSON(http.StatusBadRequest, err)
		logrus.Infof("some error %s", err)
	}
	if TranBalVal(reqBody) {
		ok := TransferBalance(reqBody)
		h.JSON(http.StatusOK, ok)
	} else {
		h.JSON(http.StatusBadRequest, "some error")
	}
}

// @Summary reserveBalnce
// @Tags post
// @Description balance reservation
// @Accept json
// @Produce json
// @Param input body ReserveCreateModel true "user_id"
// @Success 204 {integer} integer
// @Failure 400 {object} errors.ResponseError
// @Router /reserve/make [post]
func ReserveBal(h *gin.Context) {
	// Резервация средств на счетепри условии наличия средств
	var reqBody ReserveCreateModel
	if err := h.BindJSON(&reqBody); err != nil {
		h.JSON(http.StatusBadRequest, err)
		logrus.Infof("some error %s", err)
	}
	if ReserBalVal(reqBody) {
		ok := ReserveBalance(reqBody)
		h.JSON(http.StatusOK, ok)
	} else {
		h.JSON(http.StatusBadRequest, "some error")
	}
}

// @Summary confrumBalance
// @Tags post
// @Description confirmation or cancellation of the operation
// @Accept json
// @Produce json
// @Param input body ReserveConfirmModel true "user_id"
// @Success 204 {integer} integer
// @Failure 400 {object} errors.ResponseError
// @Router /reserve/accept [post]
func ConfBal(h *gin.Context) {
	// Подтверждение операции
	var reqBody ReserveConfirmModel
	if err := h.BindJSON(&reqBody); err != nil {
		h.JSON(http.StatusBadRequest, err)
		logrus.Infof("some error %s", err)
	}
	ok := Confirmation(reqBody)
	h.JSON(http.StatusOK, ok)
}

// @Summary reportService
// @Tags get
// @Description servicec report
// @Accept json
// @Produce json
// @Param input body ReportServiceStructModel true "user_id"
// @Success 200 {integer} integer
// @Failure 400 {object} errors.ResponseError
// @Router /report [get]
func ReportSer(h *gin.Context) {
	// Создание отчета по операциям
	var reqBody ReportServiceStructModel
	if err := h.BindJSON(&reqBody); err != nil {
		h.JSON(http.StatusBadRequest, err)
		logrus.Infof("some error %s", err)
	}
	ok := ReportService(reqBody)
	h.JSON(http.StatusOK, ok)
}

// @Summary operationReport
// @Tags get
// @Description open report
// @Accept json
// @Produce json
// @Success 200 {integer} integer
// @Failure 400 {object} errors.ResponseError
// @Router /reserve/csv/:file [get]
func OpenReport(h *gin.Context) {
	// Открытие отчета по операциям
	url := h.Param("file")
	if _, err := os.Stat("Report/" + url); os.IsNotExist(err) {
		h.JSON(http.StatusBadRequest, err)
		logrus.Infof("some error %s", err)
	}
	h.File("Report/" + url)
}

// @Summary reportHistory
// @Tags get
// @Description user transaction history
// @Accept json
// @Produce json
// @Param input body ReportOperationRequestModel true "user_id"
// @Success 200 {integer} integer
// @Failure 400 {object} errors.ResponseError
// @Router /balance/list [get]
func ReportHist(h *gin.Context) {
	// История транзакций пользователя
	var reqBody ReportOperationRequestModel
	if err := h.BindJSON(&reqBody); err != nil {
		h.JSON(http.StatusBadRequest, err)
		logrus.Infof("some error %s", err)
	}
	ok := ReportHistoryUser(reqBody)
	h.JSON(http.StatusOK, ok)
}
