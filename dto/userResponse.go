package dto

type UserResponse struct {
	Id         string `json:"user_id"`
	Username   string `json:"username"`
	CustomerId string `json:"customer_id"`
	Email      string `json:"email"`
	Role       string `json:"role"`
}
