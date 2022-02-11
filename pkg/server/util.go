package server

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/yockii/qscore/pkg/util"
)

func ParsePaginationInfoFromQuery(ctx *fiber.Ctx) (limit, offset int, orderBy string, err error) {
	sizeStr := ctx.Query("limit", "10")
	offsetStr := ctx.Query("offset", "0")
	limit, err = strconv.Atoi(sizeStr)
	if err != nil {
		return
	}
	offset, err = strconv.Atoi(offsetStr)
	if err != nil {
		return
	}
	if limit < -1 || limit > 1000 {
		limit = 10
	}
	if offset < -1 {
		offset = 0
	}
	orderBy = ctx.Query("orderBy") // orderBy=xxx-desc,yyy-asc,zzz
	if orderBy != "" {
		obs := strings.Split(orderBy, ",")
		ob := ""
		for _, s := range obs {
			kds := strings.Split(s, "-")
			ob += ", " + util.SnakeString(strings.TrimSpace(kds[0]))
			if len(kds) == 2 {
				d := strings.ToLower(kds[1])
				if d == "desc" {
					ob += " DESC"
				}
			}
		}
		orderBy = ob
	}
	return
}
