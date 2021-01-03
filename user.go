package main

// User ...
type User struct {
	ID          string      `json:"id"`
	UID         string      `json:"uid"`
	Avatar      *string     `json:"avatar"`
	LoginName   string      `json:"login_name"`
	Email       string      `json:"email"`
	DisplayName string      `json:"display_name"`
	FirstName   string      `json:"first_name"`
	LastName    string      `json:"last_name"`
	Company     string      `json:"company"`
	Groups      []GroupMini `json:"groups"`
	MobilePhone string      `json:"mobile_phone"`
	WorkPhone   string      `json:"work_phone"`
	HomePhone   string      `json:"home_phone"`
	IsActive    int         `json:"is_active"`
	Status      int         `json:"status"`
}

// GroupMini ...
type GroupMini struct {
	ID          string   `json:"id"`
	DisplayName string   `json:"display_name"`
	Rights      []string `json:"rights"`
}
