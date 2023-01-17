package utils

func GetPagination(limit, page int32) (newLimit, offset int) {
	if page > 0 {
		offset = int(page - 1)
	}
	if limit > 0 {
		newLimit = int(limit)
	}
	offset = offset * newLimit
	return
}
