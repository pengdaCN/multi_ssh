package tools

import "fmt"

func PanicToErr(fn func()) (err error) {
	defer func() {
		if e := recover(); e != nil {
			if _e, ok := e.(error); ok {
				err = _e
			}
			err = fmt.Errorf("panic info: %v", e)
		}
	}()
	fn()
	return
}
