package models

import (
	"time"

	"github.com/google/uuid"
)

// User struct to describe User object.
type User struct {
	ID        uuid.UUID  `db:"id" json:"id" validate:"required,uuid"`
	Username  string     `db:"username" json:"username" validate:"required,username,lte=255"`
	Email     string     `db:"email" json:"email" validate:"required,email,lte=255"`
	Password  string     `db:"password" json:"password,omitempty" validate:"required,lte=255"`
	RoleID    string     `db:"role_id" json:"roleId" validate:"required"`
	Status    int        `db:"status" json:"status" validate:"required,len=1"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy *uuid.UUID `db:"updated_by" json:"updatedBy"`
	IsDeleted bool       `db:"is_deleted" json:"isDeleted"`
}

type Renew struct {
	RefreshToken string `json:"refresh_token"`
}

// SignIn struct to describe login user.
type SignIn struct {
	Username string `json:"username" validate:"required,username"`
	Password string `json:"password" validate:"required"`
}
