package ui

import (
	"fmt"
	"github.com/kr/pretty"
	"reflect"
	"testing"
)

func TestLogic(t *testing.T) {
	var (
		colA int     = 1
		colB float64 = 2.1
	)
	dset := newDataSet()
	dset.add(1, []reflect.Value{
		reflect.ValueOf(colA),
		reflect.ValueOf(colB),
	})
	colA = 32
	colB = 0.001
	dset.add(2, []reflect.Value{
		reflect.ValueOf(colA),
		reflect.ValueOf(colB),
	})
	fmt.Printf("%#v\n", pretty.Formatter(dset.columns[0].slice.Index(1).Int()))
}
