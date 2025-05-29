package models

import (
	"github.com/go-playground/validator/v10"
)

type User struct {
	ID           int    `json:"id" validate:"required,gt=0"`
	FirstName    string `json:"first_name" validate:"required,min=1,max=50"`
	LastName     string `json:"last_name" validate:"required,min=1,max=50"`
	Email        string `json:"email" validate:"required,email"`
	CreatedAt    int64  `json:"created_at" validate:"required,gt=0"`
	DeletedAt    int64  `json:"deleted_at"`
	MergedAt     int64  `json:"merged_at"`
	ParentUserID int    `json:"parent_user_id"`
}

func (u *User) Validate() error {
	validate := validator.New()

	/* - Out of scope , but I could add custom validations here

	if u.DeletedAt > 0 && u.DeletedAt < u.CreatedAt {
		return fmt.Errorf("deleted_at cannot be before created_at")
	}

	if u.MergedAt > 0 && u.MergedAt < u.CreatedAt {
		return fmt.Errorf("merged_at cannot be before created_at")
	}

	if u.ParentUserID > 0 && u.ParentUserID == u.ID {
		return fmt.Errorf("user cannot be their own parent")
	}
	*/
	return validate.Struct(u)
}
