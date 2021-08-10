package restscenario

import (
	"net/http"
	"testing"
)

func TestRestFull(t *testing.T) {
	t.Run("rest api full test for mock driver", func(t *testing.T) {
		SetUpForRest()

		tc := TestCases{
			Name:                 "healthy",
			EchoFunc:             "Healthy",
			HttpMethod:           http.MethodGet,
			WhenURL:              "/ladybug/healthy",
			GivenQueryParams:     "",
			GivenParaNames:       nil,
			GivenParaVals:        nil,
			GivenPostData:        "",
			ExpectStatus:         http.StatusOK,
			ExpectBodyStartsWith: `cb-barista cb-ladybug`,
		}
		EchoTest(t, tc)

		TearDownForRest()
	})

}
