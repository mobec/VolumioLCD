package logger

import (
	"fmt"
	"log"
)

// // LogHTTP to log http operations, such as GET, PUSH and DELETE
// func LogHTTP(inner http.Handler, name string) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()

// 		inner.ServeHTTP(w, r)

// 		log.Printf(
// 			"\033[1;37m[HTTP]\033[0;37m %s\t%s\t%s\t%s",
// 			r.Method,
// 			r.RequestURI,
// 			name,
// 			time.Since(start),
// 		)
// 	})
// }

// Warningf logs events that require attention, but should not interrupt excecution
func Warningf(format string, v ...interface{}) {
	log.Printf(
		"\033[1;33m[WARNING]\033[0;37m %s",
		fmt.Sprintf(format, v...),
	)
}

// Infof logs uncritical events
func Infof(format string, v ...interface{}) {
	log.Printf(
		"\033[1;36m[INFO]\033[0;37m %s",
		fmt.Sprintf(format, v...),
	)
}

// Errorf prints error information and panics with the error string
func Errorf(format string, v ...interface{}) {
	errstr := fmt.Sprintf("\033[1;31m[ERROR]\033[1;37m %s\033[0;37m", fmt.Sprintf(format, v...))
	log.Printf(errstr)
	panic(errstr)
}

// AssertInfof asserts the expression with a info
func AssertInfof(expr bool, format string, v ...interface{}) {
	if expr == false {
		Infof(format, v...)
	}
}

// AssertWarningf asserts the expression with a warning
func AssertWarningf(expr bool, format string, v ...interface{}) {
	if expr == false {
		Warningf(format, v...)
	}
}

// Assertf asserts the expression with an error. Should be used to notify programmer
func Assertf(expr bool, format string, v ...interface{}) bool {
	if expr == false {
		Errorf(format, v...)
		return true
	}
	return false
}
