package user

// Валидация баланса, чтобы был больше нуля

/*func TrancBalancVal(z BalGetReturnModel) bool {
	if z.Balance >= 0 {
		return true
	} else {
		return false
	}
}*/

func TranBalVal(z TransferBalanceModel) bool {
	if z.Sum >= 0 {
		return true
	} else {
		return false
	}
}

func ReserBalVal(z ReserveCreateModel) bool {
	if z.Sum >= 0 {
		return true
	} else {
		return false
	}
}
