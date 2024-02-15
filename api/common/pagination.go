package common

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

const (
	defaultPageLimit  = 50
	maxPageLimit      = 200
	pageQueryArgsKey  = "page"
	limitQueryArgsKey = "limit"
)

// GetPaginationProperties retrieve page and limit params from query string parameters
func GetPaginationProperties(c *fiber.Ctx) (uint, uint, error) {
	var page uint
	var limit uint = defaultPageLimit

	pArg := c.Query(pageQueryArgsKey)
	if pArg != "" {
		p, err := strconv.ParseUint(pArg, 10, 64)
		if err != nil {
			return 0, 0, errors.New("invalid page param")
		}
		page = uint(p)
	}

	lArg := c.Query(limitQueryArgsKey)
	if lArg != "" {
		l, err := strconv.ParseUint(lArg, 10, 64)
		if err != nil {
			return 0, 0, errors.New("invalid limit param")
		}
		limit = uint(l)
		if limit > maxPageLimit {
			limit = maxPageLimit
		}
	}

	return page, limit, nil
}
