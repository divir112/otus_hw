package apperror

import "errors"

var ErrorParseConfig = errors.New("can't parse config")
var ErrTimeEndLessStart = errors.New("time date end event less date start")
var ErrNotFound = errors.New("not found")
var ErrDateBusy = errors.New("date is busy")
var ErrIncorrectDate = errors.New("incorrect date")
