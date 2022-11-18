package user

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

//type Pool struct {
//	pool *pgxpool.Pool
//}

// Тут написаны все функции для работы с базой данных

var database *pgx.Conn

func SetDataBase(databaseNew *pgx.Conn) {
	database = databaseNew
}

func Hello() string {
	hello := "Welcome to my microservice for user blanc"
	return hello
}

func BalanceGet(data BalGetModel) BalGetReturnModel {
	// Получение баланса пользователя из бд
	var sum float64
	var returned BalGetReturnModel
	err := database.QueryRow(context.Background(), "select balance from users where user_id=$1", data.UserId).Scan(&sum)
	if err != nil {
		logrus.Info(nil)
	}
	returned.Balance = sum
	return returned
}

func AddBalance(data BalAddModel) BalAddReturnModel {
	// Пополнение баланса определенного пользователя
	var returnData BalAddReturnModel
	var isok bool
	err := database.QueryRow(context.Background(), "SELECT EXISTS (SELECT * FROM users WHERE user_id = $1)",
		data.UserId).Scan(&isok)
	if err != nil {
		returnData.Status = "false"
		logrus.Infof("database %s", err)
		return returnData
	}

	// Обновление информации о зачислении средств
	if isok {
		database.Exec(context.Background(), "UPDATE users SET balance=$2 + balance WHERE user_id = $1",
			data.UserId, data.Balance)

	} else {
		database.Exec(context.Background(), "INSERT INTO users (user_id ,balance, reserve) VALUES ($1, $2, 0)",
			data.UserId, data.Balance)
	}

	// Добавление записи в историю
	database.Exec(context.Background(), "INSERT INTO operation_history (amount , user_id, time_history) VALUES ($1, $2, $3)",
		data.Balance, data.UserId, time.Now())
	returnData.Status = "true"
	return returnData
}

func ReserveBalance(data ReserveCreateModel) ReserveCreateReturnModel {
	// Резерв средст у пользователя
	var returnData ReserveCreateReturnModel
	var isok bool
	var isOk bool
	var balance float64
	var reserve float64
	database.QueryRow(context.Background(), "SELECT EXISTS (SELECT * FROM users WHERE user_id= $1)",
		data.UserId).Scan(&isok)
	if !(isok) {
		returnData.Status = "User not found"
		logrus.Infof("User not found ")
		return returnData
	}

	// Услуга не найдена
	database.QueryRow(context.Background(), "SELECT EXISTS (SELECT * FROM services WHERE services_id = $1)",
		data.ServiceId).Scan(&isOk)
	if !(isOk) {
		returnData.Status = "Service not found"
		logrus.Info("Service not found")
		return returnData
	}

	// Заказ уже существует
	database.QueryRow(context.Background(), "SELECT EXISTS (SELECT * FROM orders WHERE order_id = $1)",
		data.OrderId).Scan(&isOk)
	if isOk {
		returnData.Status = "This order already exists"
		logrus.Info("This order already exists")
		return returnData
	}

	// Не хватает средств
	database.QueryRow(context.Background(), "SELECT balance, reserve FROM users WHERE user_id = $1",
		data.UserId).Scan(&balance, &reserve)
	if (balance - reserve - data.Sum) < 0 {
		returnData.Status = "Not enough funds"
		logrus.Info("Not enough funds")
		return returnData
	}

	// Ошибка записи
	g, err := database.Begin(context.Background())
	if err != nil {
		returnData.Status = "Recording failed"
		logrus.Infof("Recording failed %s", err)
		return returnData
	}

	// Создание записей в бд в случае соблюдения условий
	defer g.Rollback(context.Background())
	g.Exec(context.Background(), "UPDATE users SET reserve=reserve+$1 WHERE user_id = $2", data.Sum, data.UserId)
	returnData.Status = "reserve"
	g.Exec(context.Background(), "INSERT INTO orders (order_id, user_id, services_id, price, time_operation, status_order) "+
		"VALUES ($1, $2, $3, $4, $5, $6)", data.OrderId, data.UserId, data.ServiceId, data.Sum, time.Now(), returnData.Status)
	err = g.Commit(context.Background())
	if err != nil {
		returnData.Status = "Write error"
		logrus.Infof("Write error %s", err)
		return returnData
	}
	returnData.Status = "Successfully"
	logrus.Info("Successfully")
	return returnData
}

func Confirmation(data ReserveConfirmModel) ReserveCreateReturnModel {
	// Подтверждение операции при определенных условиях
	var returnData ReserveCreateReturnModel
	var isok bool
	status := "reserve"
	database.QueryRow(context.Background(), "SELECT EXISTS (SELECT * FROM orders WHERE order_id = $1 AND "+
		"user_id = $2 AND services_id = $3 AND price = $4 AND status_order = $5)",
		data.OrderId, data.UserId, data.ServiceId, data.Sum, status).Scan(&isok)
	if !(isok) {
		returnData.Status = "Order not found"
		logrus.Info("Order not found")
		return returnData
	}

	// Команда подтверждения и обновление информации в бд
	if data.Command == "approved" {
		z, err := database.Begin(context.Background())
		if err != nil {
			returnData.Status = "Write error"
			logrus.Infof("Write error %s", err)
			return returnData
		}
		defer z.Rollback(context.Background())
		z.Exec(context.Background(), "UPDATE users SET reserve = reserve - $2 WHERE user_id = $1",
			data.UserId, data.Sum)
		z.Exec(context.Background(), "UPDATE users SET balance = balance - $2 WHERE user_id = $1",
			data.UserId, data.Sum)
		result := "approved"
		z.Exec(context.Background(), "UPDATE orders SET status_order = $1 WHERE order_id = $2",
			result, data.OrderId)
		err = z.Commit(context.Background())
		if err != nil {
			returnData.Status = "Write error"
			logrus.Infof("Write error %s", err)
			return returnData
		}
	}

	// Команда отмены и обновление информации в бд
	if data.Command == "cancel" {
		z, err := database.Begin(context.Background())
		if err != nil {
			returnData.Status = "Write error"
			logrus.Infof("Write error %s", err)
			return returnData
		}
		defer z.Rollback(context.Background())
		z.Exec(context.Background(), "UPDATE users SET reserve = reserve - $2 WHERE user_id = $1",
			data.UserId, data.Sum)
		response := "cancel"
		z.Exec(context.Background(), "UPDATE orders SET status_order = $1 WHERE order_id = $2",
			response, data.OrderId)
		err = z.Commit(context.Background())
		if err != nil {
			returnData.Status = "Write error"
			logrus.Infof("Write error %s", err)
			return returnData
		}
	}
	returnData.Status = "Successfully"
	logrus.Info("Confirmation successfully")
	return returnData
}

func TransferBalance(data TransferBalanceModel) TransferBalanceReturnModel {
	// Перевод от одного пользователя к другому
	var returnData TransferBalanceReturnModel
	var isok bool
	var sum float64
	var reserve float64
	// проверка условий существования пользователей и наличия средств
	database.QueryRow(context.Background(), "SELECT EXISTS (SELECT * FROM users WHERE user_id = $1)",
		data.From).Scan(&isok)
	if !(isok) {
		returnData.Status = "User not found"
		logrus.Info("User not found")
		return returnData
	}
	database.QueryRow(context.Background(), "SELECT EXISTS (SELECT * FROM users WHERE user_id = $1)",
		data.To).Scan(&isok)
	if !(isok) {
		returnData.Status = "User not found"
		logrus.Info("User not found")
		return returnData
	}
	database.QueryRow(context.Background(), "SELECT balance, reserve FROM users WHERE user_id = $1",
		data.From).Scan(&sum, &reserve)
	if (sum - reserve - data.Sum) < 0 {
		returnData.Status = "Not enough funds"
		logrus.Info("Not enough funds")
		return returnData
	}

	// Запись в бд
	z, err := database.Begin(context.Background())
	if err != nil {
		returnData.Status = "Write error"
		logrus.Infof("Write error %s", err)
		return returnData
	}
	defer z.Rollback(context.Background())

	// Обновление информации в бд о совершенном переводе
	z.Exec(context.Background(), "UPDATE users SET balance = balance - $1 WHERE user_id = $2",
		data.Sum, data.From)
	z.Exec(context.Background(), "UPDATE users SET balance = balance + $1 WHERE user_id = $2",
		data.Sum, data.To)
	z.Exec(context.Background(), "INSERT INTO operation_history (amount, user_id, time_history) "+
		"VALUES ($1, $2, $3)", data.Sum, data.To, time.Now())
	errCommit := z.Commit(context.Background())
	if errCommit != nil {
		returnData.Status = "Write error"
		logrus.Infof("Write error %s", errCommit)
		logrus.Info(err)
		log.Fatal(errCommit)
		return returnData
	}
	returnData.Status = "Successfully"
	return returnData
}

func ReportService(data ReportServiceStructModel) ReportServiceStructReturnModel {
	// Создание отчета с форматом .csv
	var returnData ReportServiceStructReturnModel
	var serviceId int
	var nameService string
	var price float64
	var allMoney float64 = 0
	var rowsNum []ServiceNameModel

	// Преобразование входных параметров
	yearString := strings.Split(data.Date, "-")[0]
	yearInt, _ := strconv.Atoi(yearString)
	monthString := strings.Split(data.Date, "-")[1]
	monthInt, _ := strconv.Atoi(monthString)
	dateStart := time.Date(yearInt, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	dateEnd := dateStart.AddDate(0, 1, -1)

	// Создание файла, если файл есть с таким именем обновляет его
	rows, err := database.Query(context.Background(), "SELECT services_id, name_ser FROM services")
	if err != nil {
		logrus.Infof("Error find service %s", err)
	}
	dataFinal := [][]string{{"Service", "Revenue"}}
	fileName := "Report revenue " + data.Date + ".csv"
	errMkdir := os.MkdirAll("Report", 0644)
	if errMkdir != nil {
		panic(err)
	}
	allFile, err := os.OpenFile("Report/"+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}
	status := "approved"
	defer allFile.Close()

	// Заполнение файла
	for rows.Next() {
		errScan := rows.Scan(&serviceId, &nameService)
		if errScan != nil {
			logrus.Infof("Error scan database %s", errScan)
		}
		rowsNum = append(rowsNum, ServiceNameModel{ServiceId: serviceId, ServiceName: nameService})

	}
	rows.Close()
	for i := 0; i < len(rowsNum); i++ {
		rowsCont, errCount := database.Query(context.Background(), "SELECT price FROM orders WHERE services_id = "+
			"$1 AND status_order = $2 AND time_operation > $3 AND time_operation < $4",
			rowsNum[i].ServiceId, status, dateStart, dateEnd)
		if errCount != nil {
			logrus.Infof("Error for cont rows %s", errCount)
		}
		for rowsCont.Next() {
			errPrice := rowsCont.Scan(&price)
			allMoney = allMoney + price
			if errPrice != nil {
				logrus.Infof("Revenue calculation error %s", errPrice)
			}
		}
		s := fmt.Sprintf("%v", allMoney)
		dataReport := []string{rowsNum[i].ServiceName, s}
		dataFinal = append(dataFinal, dataReport)
		rowsCont.Close()
	}

	// Возвращается путь и имя файла
	w := csv.NewWriter(allFile)
	w.WriteAll(dataFinal)
	returnData.Url = "localhost:8080/csv/" + fileName
	return returnData
}

func ReportHistoryUser(data ReportOperationRequestModel) ReportOperationRequestReturnModel {
	// Создание отчета по операциям пользователя с возможностью сортировки
	var serviceId int
	var serviceName string
	var rowsFind pgx.Rows
	var errFind error
	var forReturn ReportOperationRequestTempModel
	var sum float64
	var created time.Time
	var returnData ReportOperationRequestReturnModel
	var numRows []ServiceNameModel
	// Берем имя пользователя и номер услуги для составления отчета
	rows, err := database.Query(context.Background(), "SELECT services_id, name_ser FROM services")
	if err != nil {
		logrus.Infof("database error %s", err)
	}
	for rows.Next() {
		errRow := rows.Scan(&serviceId, &serviceName)
		if errRow != nil {
			logrus.Infof("some error %s", errRow)
		}
		numRows = append(numRows, ServiceNameModel{ServiceId: serviceId, ServiceName: serviceName})

	}
	rows.Close()
	skip := data.Rows * (data.Page - 1)

	// Параметры сортировки
	if data.Sort == "price_up" {
		rowsFind, errFind = database.Query(context.Background(), "SELECT services_id, price, time_operation FROM orders WHERE user_id = $1 AND status_order = $2 ORDER BY price ASC OFFSET $3 LIMIT $4",
			data.UserId, "approved", skip, data.Rows)
		if errFind != nil {
			logrus.Infof("find orders error %s", err)
		}
	}
	if data.Sort == "price_down" {
		rowsFind, errFind = database.Query(context.Background(), "SELECT services_id, price, time_operation FROM orders WHERE user_id = $1 AND status_order = $2 ORDER BY price DESC OFFSET $3 LIMIT $4",
			data.UserId, "approved", skip, data.Rows)
		if errFind != nil {
			logrus.Infof("find orders error %s", err)
		}
	}
	if data.Sort == "created_up" {
		rowsFind, errFind = database.Query(context.Background(), "SELECT services_id, price, time_operation FROM orders WHERE user_id = $1 AND status_order = $2 ORDER BY time_operation ASC OFFSET $3 LIMIT $4",
			data.UserId, "approved", skip, data.Rows)
		if errFind != nil {
			logrus.Infof("find orders error %s", err)
		}
	}
	if data.Sort == "created_down" {
		rowsFind, errFind = database.Query(context.Background(), "SELECT services_id, price, time_operation FROM orders WHERE user_id = $1 AND status_order = $2 ORDER BY time_operation DESC OFFSET $3 LIMIT $4",
			data.UserId, "approved", skip, data.Rows)
		if errFind != nil {
			logrus.Infof("find orders error %s", err)
		}
	}
	// Вывод отсортированной информации
	for rowsFind.Next() {
		errRows := rowsFind.Scan(&serviceId, &sum, &created)
		if errRows != nil {
			log.Print("1")
			log.Print(errRows)
		}
		forReturn.Time = created.String()
		forReturn.Money = sum
		forReturn.Service = numRows[serviceId].ServiceName
		returnData.Orders = append(returnData.Orders, forReturn)
	}
	rowsFind.Close()
	return returnData
}
