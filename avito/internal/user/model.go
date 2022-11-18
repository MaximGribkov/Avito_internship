package user

// Структуры для работы с бд

type BalGetModel struct {
	UserId int `json:"user_id"`
}

type BalGetReturnModel struct {
	Balance float64 `json:"balance"`
}

type BalAddModel struct {
	UserId  int     `json:"user_id"`
	Balance float64 `json:"balance"`
}

type BalAddReturnModel struct {
	Status string `json:"status"`
}

type ReserveCreateModel struct {
	UserId    int     `json:"user_id"`
	ServiceId int     `json:"service_id"`
	OrderId   int     `json:"order_id"`
	Sum       float64 `json:"sum"`
}

type ReserveConfirmModel struct {
	UserId    int     `json:"user_id"`
	ServiceId int     `json:"service_id"`
	OrderId   int     `json:"order_id"`
	Sum       float64 `json:"sum"`
	Command   string  `json:"command"`
}

type ReserveCreateReturnModel struct {
	Status string `json:"status"`
}

type TransferBalanceModel struct {
	From int     `json:"from"`
	To   int     `json:"to"`
	Sum  float64 `json:"sum"`
}

type TransferBalanceReturnModel struct {
	Status string `json:"status"`
}

type ReportServiceStructModel struct {
	Date string `json:"date"`
}

type ReportServiceStructReturnModel struct {
	Url string `json:"url"`
}

type ServiceNameModel struct {
	ServiceId   int
	ServiceName string
}
type ReportOperationRequestModel struct {
	UserId int    `json:"user_id"`
	Page   int    `json:"page"`
	Rows   int    `json:"rows"`
	Sort   string `json:"sort"`
}

type ReportOperationRequestTempModel struct {
	Time    string  `json:"time"`
	Money   float64 `json:"money"`
	Service string  `json:"service"`
}

type ReportOperationRequestReturnModel struct {
	Orders []ReportOperationRequestTempModel
}
