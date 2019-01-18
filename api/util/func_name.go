package util

import "runtime"

// FuncName return current function   name
func FuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
