package utils

import (
	"fmt"
)

func SafeGo(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("SafeGo error is", err)
			}
		}()
		fn()
	}()
}
