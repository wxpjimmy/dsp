package metric

import (
	"net/http"
	"runtime"
	"reflect"
	"github.com/rcrowley/go-metrics"
	"time"
	"fmt"
)

func Decorate(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	name := GetFunctionName(f)
	t := metrics.GetOrRegisterTimer(name, nil)
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			t.UpdateSince(start)
			fmt.Println(start, time.Now())
		}()
		f(w, r)
	}
}

func Dc(f func(args ...interface{})) func(...interface{}) {
	name := GetFunctionName(f)
	t := metrics.GetOrRegisterTimer(name, nil)
	return func(args ...interface{}) {
		start := time.Now()
		defer func() {
			t.UpdateSince(start)
			fmt.Println(start, time.Now())
		}()
		f(args)
	}
}

func TimedFunc(f func()) func() {
	name := GetFunctionName(f)
	t := metrics.GetOrRegisterTimer(name, nil)
	return func() {
		start := time.Now()
		defer func() {
			t.UpdateSince(start)
			fmt.Println(start, time.Now())
		}()
		f()
	}
}


func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
