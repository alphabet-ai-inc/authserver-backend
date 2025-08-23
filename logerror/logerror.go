package logerror

import (
	"fmt"
	"runtime"
)

func LogError(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("Error at %s:%d: %v\n", file, line, err)
	}
}
