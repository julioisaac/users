package user

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/julioisaac/users/domainerror"
	"github.com/julioisaac/users/utils"
	"reflect"
	"strings"
)

var (
	// ErrNotFound when a user was not found
	ErrNotFound = domainerror.New(domainerror.GetUserNotFoundError, "user not found", nil)
)

type Repository interface {
	Create(ctx context.Context, u User) error
	GetUserByParams(ctx context.Context, f Filters) (Users, error)
}

var Insert = `INSERT INTO "users" (user_id, first_name, last_name, email_address, created_at, deleted_at, merged_at, parent_user_id) VALUES($1, $2, $3, $4, $5, $6, $7, $8);`

type defaultUserRepo struct {
	pgDb *sql.DB
}

func NewUserRepo(db *sql.DB) Repository {
	return &defaultUserRepo{pgDb: db}
}

func (d defaultUserRepo) Create(ctx context.Context, u User) error {
	args := []interface{}{
		u.UserID,
		u.FirstName,
		u.LastName,
		u.EmailAddress,
		u.CreatedAt,
		u.DeletedAt,
		u.MergedAt,
		u.ParentUserID,
	}

	_, err := d.pgDb.ExecContext(ctx, Insert, args...)

	if ctx.Err() != nil {
		return ctx.Err()
	}

	if err != nil {
		return err
	}

	return nil
}

func (d defaultUserRepo) GetUserByParams(ctx context.Context, f Filters) (Users, error) {
	sqlQuery, args := BuildDynamicSQL(f)

	rs, err := d.pgDb.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	result := Users{}
	for rs.Next() {
		u := User{}

		err := rs.Scan(
			&u.UserID,
			&u.FirstName,
			&u.LastName,
			&u.EmailAddress,
			&u.CreatedAt,
			&u.DeletedAt,
			&u.MergedAt,
			&u.ParentUserID)

		if err != nil {
			continue
		}
		result = append(result, u)
	}

	if err := rs.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func BuildDynamicSQL(f Filters) (string, []interface{}) {
	var whereConditions []string
	var args []interface{}

	fType := reflect.TypeOf(f)

	idx := 0
	for i := 0; i < fType.NumField(); i++ {
		field := fType.Field(i)
		fieldValue := reflect.ValueOf(f).Field(i).Interface()

		switch field.Name {
		case "Pagination":
		default:
			if !utils.IsNilOrEmpty(fieldValue) {
				idx++
				whereConditions = append(whereConditions, fmt.Sprintf("%s = $%d", field.Tag.Get("json"), idx))
				args = append(args, fieldValue)
			}
		}
	}

	sqlQuery := `SELECT user_id, first_name, last_name, email_address, created_at, deleted_at, merged_at, parent_user_id FROM "users"`
	if len(whereConditions) > 0 {
		sqlQuery += " WHERE " + strings.Join(whereConditions, " AND ")
	}
	sqlQuery += fmt.Sprintf(` ORDER BY created_at DESC LIMIT %d OFFSET %d `, f.Limit, f.Page)

	return sqlQuery, args
}
