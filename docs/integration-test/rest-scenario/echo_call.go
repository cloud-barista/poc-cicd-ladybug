package restscenario

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/cloud-barista/poc-cicd-ladybug/src/rest-api/router"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var funcs = map[string]interface{}{
	"Healthy":       router.Healthy,
	"ListCluster":   router.ListCluster,
	"CreateCluster": router.CreateCluster,
	"GetCluster":    router.GetCluster,
	"DeleteCluster": router.DeleteCluster,
	"ListNode":      router.ListNode,
	"AddNode":       router.AddNode,
	"GetNode":       router.GetNode,
	"RemoveNode":    router.RemoveNode,
}

func Call(m map[string]interface{}, name string, params ...interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(m[name])
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

func EchoTest(t *testing.T, tc TestCases) (string, error) {

	var (
		body string = ""
		err  error  = nil
	)

	t.Run(tc.Name, func(t *testing.T) {
		e := echo.New()
		var req *http.Request = nil
		if tc.GivenPostData != "" {
			req = httptest.NewRequest(tc.HttpMethod, "/"+tc.GivenQueryParams, bytes.NewBuffer([]byte(tc.GivenPostData)))
		} else {
			req = httptest.NewRequest(tc.HttpMethod, "/"+tc.GivenQueryParams, nil)
		}
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath(tc.WhenURL)
		if tc.GivenParaNames != nil {
			c.SetParamNames(tc.GivenParaNames...)
			c.SetParamValues(tc.GivenParaVals...)
		}

		res, err := Call(funcs, tc.EchoFunc, c)
		if assert.NoError(t, err) {
			if res != nil && !res[0].IsNil() {
				he, ok := res[0].Interface().(*echo.HTTPError)
				if ok { // echo.NewHTTPError() 로 에러를 리턴했을 경우
					assert.Equal(t, tc.ExpectStatus, he.Code)
					body = fmt.Sprintf("%v", he.Message)
				} else { // err 로 에러를 리턴했을 경우
					body = fmt.Sprintf("%v", res[0])
				}
				if tc.ExpectBodyStartsWith != "" {
					if !assert.True(t, strings.HasPrefix(body, tc.ExpectBodyStartsWith)) {
						fmt.Fprintf(os.Stderr, "\n                Not Equal(echo.NewHTTPError): \n"+
							"                  Expected Start With: %s\n"+
							"                  Actual  : %s\n", tc.ExpectBodyStartsWith, body)
					}
				}
				if tc.ExpectBodyContains != "" {
					if !assert.True(t, strings.Contains(body, tc.ExpectBodyContains)) {
						fmt.Fprintf(os.Stderr, "\n                Not Equal(echo.NewHTTPError): \n"+
							"                  Expected Contains: %s\n"+
							"                  Actual  : %s\n", tc.ExpectBodyContains, body)
					}
				}
				if tc.ExpectBodyStartsWith == "" && tc.ExpectBodyContains == "" {
					if !assert.True(t, "" == body) {
						fmt.Fprintf(os.Stderr, "\n                Not Equal(echo.NewHTTPError): \n"+
							"      Expected StartWith/Contains: %s\n"+
							"      Actual  : %s\n", tc.ExpectBodyStartsWith, body)
					}
				}
			} else {
				assert.Equal(t, tc.ExpectStatus, rec.Code)
				body = rec.Body.String()
				if tc.ExpectBodyStartsWith != "" {
					if !assert.True(t, strings.HasPrefix(body, tc.ExpectBodyStartsWith)) {
						fmt.Fprintf(os.Stderr, "\n                Not Equal(echo.Context): \n"+
							"                  Expected Start With: %s\n"+
							"                  Actual  : %s\n", tc.ExpectBodyStartsWith, body)
					}
				}
				if tc.ExpectBodyContains != "" {
					if !assert.True(t, strings.Contains(body, tc.ExpectBodyContains)) {
						fmt.Fprintf(os.Stderr, "\n                Not Equal(echo.Context): \n"+
							"                  Expected Contains: %s\n"+
							"                  Actual  : %s\n", tc.ExpectBodyContains, body)
					}
				}
				if tc.ExpectBodyStartsWith == "" && tc.ExpectBodyContains == "" {
					if !assert.True(t, "" == body) {
						fmt.Fprintf(os.Stderr, "\n                Not Equal(echo.Context): \n"+
							"      Expected StartWith/Contains: %s\n"+
							"      Actual  : %s\n", tc.ExpectBodyStartsWith, body)
					}
				}
			}
		}
	})

	return body, err
}
