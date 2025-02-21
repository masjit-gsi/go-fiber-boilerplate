package models

import (
	"time"

	"github.com/gofrs/uuid"
)

const tableNameAuthor = "authors"

type Author struct {
	ID        uuid.UUID  `db:"id" json:"id" gorm:"column:id"`
	Name      string     `db:"name" json:"name" gorm:"column:name"`
	Address   *string    `db:"address" json:"address" gorm:"column:address"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt" gorm:"column:created_at"`
	CreatedBy *uuid.UUID `db:"created_by" json:"createdBy" gorm:"column:created_by"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt" gorm:"column:updated_at"`
	UpdatedBy *uuid.UUID `db:"updated_by" json:"updatedBy" gorm:"column:updated_by"`
	IsDeleted bool       `db:"is_deleted" json:"isDeleted" gorm:"column:is_deleted"`
}

type AuthorRequest struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name" validate:"required,lte=255"`
	Address *string   `json:"address"`
	UserID  uuid.UUID `json:"-"`
}

func (*Author) TableName() string {
	return tableNameAuthor
}

var ColumnMappAuthor = map[string]interface{}{
	"id":         "id",
	"name":       "name",
	"address":    "address",
	"createdAt":  "created_at",
	"updatedAt":  "updated_at",
	"is_deleted": "is_deleted",
}

func (i *Author) BindFromRequest(req AuthorRequest) {
	var now = time.Now()
	if req.ID == uuid.Nil {
		newID, _ := uuid.NewV4()
		i.ID = newID
		i.CreatedAt = now
		i.CreatedBy = &req.UserID
		i.UpdatedAt = nil
	} else {
		i.ID = req.ID
		i.UpdatedAt = &now
		i.UpdatedBy = &req.UserID
	}

	i.Name = req.Name
	i.Address = req.Address
}

func (i *Author) SoftDelete(userID uuid.UUID) {
	i.UpdatedBy = &userID
}
