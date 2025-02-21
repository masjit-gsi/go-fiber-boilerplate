package repository

import (
	"bytes"

	"github.com/fiber-go-template/app/models"
	"github.com/fiber-go-template/config/logger"
	"github.com/fiber-go-template/database"
	"github.com/fiber-go-template/helper/pagination"
)

var (
	authorQuery = struct {
		Select string
		Count  string
	}{
		Select: `SELECT id, name, address, created_at, created_by, updated_at, updated_by, is_deleted  
				FROM authors `,
		Count: `select count(id) from authors `,
	}
)

type AuthorRepository interface {
	ResolveAll(req models.StandardRequest) (data pagination.Response, err error)
}

type AuthorRepositoryDB struct {
	DB database.DBConn
}

func NewAuthorRepository(db database.DBConn) AuthorRepository {
	return &AuthorRepositoryDB{
		DB: db,
	}
}

func (r *AuthorRepositoryDB) ResolveAll(req models.StandardRequest) (data pagination.Response, err error) {
	var params []interface{}
	var query bytes.Buffer
	query.WriteString(" WHERE coalesce(is_deleted) = false ")

	if req.Keyword != "" {
		query.WriteString(" AND ")
		query.WriteString(" concat(name, address) ilike ? ")
		params = append(params, "%"+req.Keyword+"%")
	}

	// Get count data
	queryCount := r.DB.Query().Rebind(authorQuery.Count + query.String())
	var totalData int
	err = r.DB.Query().QueryRow(queryCount, params...).Scan(&totalData)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if totalData < 1 {
		data.Items = make([]interface{}, 0)
		return
	}

	// Mapping column sorting
	query.WriteString("order by " + models.ColumnMappAuthor[req.SortBy].(string) + " " + req.SortType + " ")

	// Set Offset, Pagesize / limit
	offset := (req.PageNumber - 1) * req.PageSize
	query.WriteString("limit ? offset ? ")
	params = append(params, req.PageSize)
	params = append(params, offset)

	// Rebind params to query
	rawQuery := query.String()
	rawQuery = r.DB.Query().Rebind(authorQuery.Select + rawQuery)
	rows, err := r.DB.Query().Queryx(rawQuery, params...)
	if err != nil {
		return
	}

	// Mapping to data model
	for rows.Next() {
		var items models.Author
		err = rows.StructScan(&items)
		if err != nil {
			return
		}

		data.Items = append(data.Items, items)
	}

	// Generate meta pagination
	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)

	return
}
