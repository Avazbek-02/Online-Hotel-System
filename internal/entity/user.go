package entity

type User struct {
	ID            string `json:"id"`
	FullName      string `json:"fullname"`
	UserName      string `json:"username"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Password_hash string `json:"password"`
	UserType      string `json:"user_type"`
	UserRole      string `json:"user_role"`
	UserStatus    string `json:"user_status"`
	Gender        string `json:"gender"`
	AccessToken   string `json:"access_token"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type UserSingleRequest struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	UserRole string `json:"user_role"`
}

type UserList struct {
	Items []User `json:"users"`
	Count int    `json:"count"`
}
