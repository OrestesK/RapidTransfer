package error

import (
	"fmt"
	"runtime"
)

func NewError(message string) error {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Errorf("[%s][%d] : %s", file, line, message)
}
