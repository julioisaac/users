package usersAPI

import (
	"github.com/gofiber/fiber/v2"
	users "github.com/julioisaac/users/user"
)

const HandlerPath string = "/v1/users"

func RegisterRoutes(router fiber.Router) {
	router.Get("", GetUsersHandler)
}

type userParams struct {
	Page         uint
	Limit        uint
	UsersID      int64
	FirstName    string
	LastName     string
	EmailAddress string
	ParentUserID int64
}

func (up userParams) toUsersFilters() users.Filters {
	p := users.Pagination{
		Page:  up.Page,
		Limit: up.Limit,
	}
	return users.Filters{
		Pagination:   &p,
		UserID:       up.UsersID,
		FirstName:    up.FirstName,
		LastName:     up.LastName,
		EmailAddress: up.EmailAddress,
		ParentUserID: up.ParentUserID,
	}
}

// @Description holds the recovered user
type userResponse struct {
	UsersID      int64  `json:"users_id,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	EmailAddress string `json:"email_address,omitempty"`
	ParentUserID int64  `json:"parent_user_id,omitempty"`
}

// @Description holds the recovered users
type usersResponse []userResponse

func toUsersResponse(us users.Users) usersResponse {
	responses := usersResponse{}
	for _, u := range us {

		responses = append(responses, userResponse{
			UsersID:      u.UserID,
			FirstName:    u.FirstName,
			LastName:     u.LastName,
			EmailAddress: u.EmailAddress,
			ParentUserID: u.ParentUserID,
		})
	}
	return responses
}
