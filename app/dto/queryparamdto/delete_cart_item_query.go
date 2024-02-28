package queryparamdto

import (
	"halodeksik-be/app/apperror"
	"strconv"
	"strings"
)

type DeleteCartItem struct {
	ProductIds string `form:"product_ids" validate:"required"`
}

func (q DeleteCartItem) ToSliceOfInt64() ([]int64, error) {
	idsStr := strings.Split(q.ProductIds, ",")
	ids := make([]int64, 0)
	for _, idStr := range idsStr {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, apperror.ErrInvalidIntInString
		}
		ids = append(ids, id)
	}
	return ids, nil
}
