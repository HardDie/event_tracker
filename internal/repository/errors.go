package repository

import "errors"

var (
	ErrorInternal          = errors.New("internal error")
	ErrorEventTypeNotExist = errors.New("event type not exist")
	ErrorEventNotExist     = errors.New("event not exist")
)
