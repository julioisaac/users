//go:build unit

package user_test

import (
	"context"
	"errors"
	"github.com/julioisaac/users/providers/cache"
	cache_mock "github.com/julioisaac/users/test/mocks/cache"
	user_mock "github.com/julioisaac/users/test/mocks/user"
	users "github.com/julioisaac/users/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_service_FindByFilters(t *testing.T) {
	type fields struct {
		repo  users.Repository
		cache cache.Cache
	}
	type args struct {
		ctx     context.Context
		filters users.Filters
	}

	tests := []struct {
		name        string
		fields      fields
		args        args
		want        users.Users
		expectedErr error
		wantErr     bool
	}{
		{
			name: "should return users from cache",
			fields: fields{
				repo: new(user_mock.Repository),
				cache: func() cache.Cache {
					cacheMock := new(cache_mock.Cache)
					users := users.Users{{
						UserID:       12,
						FirstName:    "Julio",
						LastName:     "Penha",
						EmailAddress: "julioisaac7@gmail.com",
					}}
					u, _ := users.ToBytes()

					cacheMock.On("Get", mock.Anything, mock.Anything).
						Return(u, nil)
					return cacheMock
				}(),
			},
			args: args{
				ctx:     context.Background(),
				filters: users.Filters{},
			},
			want: users.Users{{
				UserID:       12,
				FirstName:    "Julio",
				LastName:     "Penha",
				EmailAddress: "julioisaac7@gmail.com",
			}},
			wantErr: false,
		},
		{
			name: "should return error trying to get user from database",
			fields: fields{
				repo: func() users.Repository {
					repoMock := new(user_mock.Repository)
					repoMock.On("GetUserByParams", mock.Anything, users.Filters{}).
						Return(nil, errors.New("failed getting user"))
					return repoMock
				}(),
				cache: func() cache.Cache {
					cacheMock := new(cache_mock.Cache)
					cacheMock.On("Get", mock.Anything, mock.Anything).
						Return(nil, nil)
					return cacheMock
				}(),
			},
			args: args{
				ctx:     context.Background(),
				filters: users.Filters{},
			},
			expectedErr: errors.New("failed getting user"),
			wantErr:     true,
		},
		{
			name: "should found user in database and add to cache",
			fields: fields{
				repo: func() users.Repository {
					repoMock := new(user_mock.Repository)
					users := users.Users{{
						UserID:       12,
						FirstName:    "Julio",
						LastName:     "Penha",
						EmailAddress: "julioisaac7@gmail.com",
					}}

					repoMock.On("GetUserByParams", mock.Anything, mock.Anything).
						Return(users, nil)
					return repoMock
				}(),
				cache: func() cache.Cache {
					cacheMock := new(cache_mock.Cache)
					cacheMock.On("Get", mock.Anything, mock.Anything).
						Return(nil, nil)

					users := users.Users{{
						UserID:       12,
						FirstName:    "Julio",
						LastName:     "Penha",
						EmailAddress: "julioisaac7@gmail.com",
					}}
					u, _ := users.ToBytes()

					cacheMock.On("Add", mock.Anything, mock.Anything, u, mock.Anything).
						Return(nil)
					return cacheMock
				}(),
			},
			args: args{
				ctx:     context.Background(),
				filters: users.Filters{},
			},
			want: users.Users{{
				UserID:       12,
				FirstName:    "Julio",
				LastName:     "Penha",
				EmailAddress: "julioisaac7@gmail.com",
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := users.NewService(tt.fields.repo, tt.fields.cache)

			got, err := s.FindByFilters(tt.args.ctx, tt.args.filters)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equalf(t, tt.want, got, "FindByFilters(%v, %v)", tt.args.ctx, tt.args.filters)
			}
		})
	}
}
