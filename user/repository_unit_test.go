//go:build unit

package user_test

import (
	users "github.com/julioisaac/users/user"
	"reflect"
	"testing"
)

func Test_buildDynamicSQL(t *testing.T) {
	type args struct {
		filters users.Filters
	}
	tests := []struct {
		name          string
		args          args
		queryExpected string
		argsExpected  []interface{}
	}{
		{
			name: "should create dynamic query",
			args: func() args {
				filters := users.Filters{
					Pagination: &users.Pagination{
						Page:  0,
						Limit: 10,
					},
					UserID:    12,
					FirstName: "Julio",
				}

				return args{filters}
			}(),
			queryExpected: `SELECT user_id, first_name, last_name, email_address, created_at, deleted_at, merged_at, parent_user_id FROM "users" WHERE user_id = $1 AND first_name = $2 ORDER BY created_at DESC LIMIT 10 OFFSET 0 `,
			argsExpected:  []interface{}{int64(12), "Julio"},
		},
		{
			name: "should create dynamic query with default fields",
			args: func() args {
				filters := users.Filters{
					Pagination: &users.Pagination{
						Page:  0,
						Limit: 10,
					},
					UserID:   10,
					LastName: "Penha",
				}

				return args{filters}
			}(),
			queryExpected: `SELECT user_id, first_name, last_name, email_address, created_at, deleted_at, merged_at, parent_user_id FROM "users" WHERE user_id = $1 AND last_name = $2 ORDER BY created_at DESC LIMIT 10 OFFSET 0 `,
			argsExpected:  []interface{}{int64(10), "Penha"},
		},
		{
			name: "should ignore invalid fields when creating dynamic query",
			args: func() args {
				filters := users.Filters{
					Pagination: &users.Pagination{
						Page:  0,
						Limit: 10,
					},
					UserID:   10,
					LastName: "Penha",
				}

				return args{filters}
			}(),
			queryExpected: `SELECT user_id, first_name, last_name, email_address, created_at, deleted_at, merged_at, parent_user_id FROM "users" WHERE user_id = $1 AND last_name = $2 ORDER BY created_at DESC LIMIT 10 OFFSET 0 `,
			argsExpected:  []interface{}{int64(10), "Penha"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := users.BuildDynamicSQL(tt.args.filters)
			if got != tt.queryExpected {
				t.Errorf("BuildDynamicSQL() \n got = %v, \n want = %v", got, tt.queryExpected)
			}
			if !reflect.DeepEqual(got1, tt.argsExpected) {
				t.Errorf("BuildDynamicSQL() \n got1 = %v, \n want = %v", got1, tt.argsExpected)
			}
		})
	}
}
