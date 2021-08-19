package goscenario

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Call(any interface{}, name string, params []interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(any).MethodByName(name)
	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of params is not adapted.")
		return
	}

	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)

	return
}

func MethodTest(t *testing.T, tc TestCases) (string, error) {

	var (
		res string = ""
		err error  = nil
	)

	t.Run(tc.Name, func(t *testing.T) {

		rv, err := Call(tc.Instance, tc.Method, tc.Args)
		if assert.NoError(t, err) {
			if rv != nil && !rv[1].IsNil() {
				res = fmt.Sprintf("%v", rv[1])
				msgSplit := strings.SplitAfter(res, "method: ")
				if len(msgSplit) > 0 {
					res = msgSplit[len(msgSplit)-1]
				}
			} else {
				res = fmt.Sprintf("%v", rv[0])

				if strings.HasPrefix(res, "{") {
					dst := new(bytes.Buffer)
					err = json.Compact(dst, []byte(res))
					assert.NoError(t, err)
					res = dst.String()
				}
			}
			if tc.ExpectResStartsWith != "" {
				if !assert.True(t, strings.HasPrefix(res, tc.ExpectResStartsWith)) {
					fmt.Fprintf(os.Stderr, "\n                Not Equal: \n"+
						"                  Expected Start With: %s\n"+
						"                  Actual  : %s\n", tc.ExpectResStartsWith, res)
				}
			}
			if tc.ExpectResContains != "" {
				if !assert.True(t, strings.Contains(res, tc.ExpectResContains)) {
					fmt.Fprintf(os.Stderr, "\n                Not Equal: \n"+
						"                  Expected Contains: %s\n"+
						"                  Actual  : %s\n", tc.ExpectResContains, res)
				}
			}
			if tc.ExpectResStartsWith == "" && tc.ExpectResContains == "" {
				if !assert.True(t, "" == res) {
					fmt.Fprintf(os.Stderr, "\n                Not Equal: \n"+
						"      Expected StartWith/Contains: %s\n"+
						"      Actual  : %s\n", tc.ExpectResStartsWith, res)
				}
			}

		}

	})

	return res, err
}
