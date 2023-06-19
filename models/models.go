package models

type User struct {
	ID       int    `json:"-"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserOrder struct {
	ID      int
	userID  int
	orderID int
}

type Order struct {
	ID         int     `json:"-"`
	Number     int     `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual"`
	UploadedAt string  `json:"uploaded_at"`
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Withdraw struct {
	Order       int    `json:"order"`
	Sum         int    `json:"sum"`
	ProcessedAt string `json:"processed_at"`
}

type Number struct {
	Order   int    `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}
