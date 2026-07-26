package dto

import "time"

type UserResponse struct {
	ID uint `json:"id"`

	Name      string     `json:"name"`
	LastName  string     `json:"lastName"`
	Contact   string     `json:"contact"`
	Email     string     `json:"email"`
	LastLogin *time.Time `json:"lastLogin,omitempty"`

	Business *BusinessResponse `json:"business,omitempty"`
	Role     *RoleResponse     `json:"role,omitempty"`
}
