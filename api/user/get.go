package usersAPI

import (
	"github.com/gofiber/fiber/v2"
	"github.com/julioisaac/users/api/common"
	"github.com/julioisaac/users/ioc"
	"net/http"
)

// GetUsersHandler @Summary Get Users
// @Description  Get users in cache or db
// @Tags         users
// @Accept       json
// @Produce      json
// @Param users_id query string false "users_id"
// @Param first_name query string false "first_name"
// @Param last_name query string false "last_name"
// @Param email_address query string false "email_address"
// @Param parent_user_id query string false "parent_user_id"
// @Param page query string false "page"
// @Param limit query string false "limit"
// @Success 200 {object} usersResponse
// @Failure 400 {object} common.ApiError "when some value of the request is invalid"
// @Failure 404 {object} common.ApiError "when the request was not found"
// @Failure 500 {object} common.ApiError "when something was wrong when processing request"
// @Router       /api/v1/users [get]
func GetUsersHandler(c *fiber.Ctx) error {

	page, limit, err := common.GetPaginationProperties(c)
	if err != nil {
		return err
	}

	f := userParams{
		UsersID:      int64(c.QueryInt("users_id", 0)),
		FirstName:    c.Query("first_name"),
		LastName:     c.Query("last_name"),
		EmailAddress: c.Query("email_address"),
		ParentUserID: int64(c.QueryInt("parent_user_id", 0)),
		Page:         page,
		Limit:        limit,
	}

	userService := ioc.UserService()
	result, err := userService.FindByFilters(c.UserContext(), f.toUsersFilters())
	if err != nil {
		return common.WriteError(c, err)
	}

	return common.WriteResponse(c, toUsersResponse(result), http.StatusOK)

}
