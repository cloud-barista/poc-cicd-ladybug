package restscenario

import (
	"net/http"
	"testing"
)

func TestRestEmpty(t *testing.T) {
	t.Run("rest api empty test for mock driver", func(t *testing.T) {
		SetUpForRest()

		tc := TestCases{
			Name:                 "list cluster",
			EchoFunc:             "ListCluster",
			HttpMethod:           http.MethodGet,
			WhenURL:              "/ladybug/ns/:namespace/clusters",
			GivenQueryParams:     "",
			GivenParaNames:       []string{"namespace"},
			GivenParaVals:        []string{"ns-unit-01"},
			GivenPostData:        "",
			ExpectStatus:         http.StatusOK,
			ExpectBodyStartsWith: `{"kind":"ClusterList"`,
		}
		EchoTest(t, tc)

		tc = TestCases{
			Name:                 "get cluster",
			EchoFunc:             "GetCluster",
			HttpMethod:           http.MethodGet,
			WhenURL:              "/ladybug/ns/:namespace/clusters/:cluster",
			GivenQueryParams:     "",
			GivenParaNames:       []string{"namespace", "cluster"},
			GivenParaVals:        []string{"ns-unit-01", "cb-cluster"},
			GivenPostData:        "",
			ExpectStatus:         http.StatusNotFound,
			ExpectBodyStartsWith: `{"message":"/ns/ns-unit-01/clusters/cb-cluster not found"}`,
		}
		EchoTest(t, tc)

		tc = TestCases{
			Name:                 "list node",
			EchoFunc:             "ListNode",
			HttpMethod:           http.MethodGet,
			WhenURL:              "/ladybug/ns/:namespace/clusters/:cluster/nodes",
			GivenQueryParams:     "",
			GivenParaNames:       []string{"namespace", "cluster"},
			GivenParaVals:        []string{"ns-unit-01", "cb-cluster"},
			GivenPostData:        "",
			ExpectStatus:         http.StatusOK,
			ExpectBodyStartsWith: `{"kind":"NodeList","items":[`,
		}
		EchoTest(t, tc)

		tc = TestCases{
			Name:                 "get node",
			EchoFunc:             "GetNode",
			HttpMethod:           http.MethodGet,
			WhenURL:              "/ladybug/ns/:namespace/clusters/:cluster/nodes/:node",
			GivenQueryParams:     "",
			GivenParaNames:       []string{"namespace", "cluster", "node"},
			GivenParaVals:        []string{"ns-unit-01", "cb-cluster", "sample"},
			GivenPostData:        "",
			ExpectStatus:         http.StatusNotFound,
			ExpectBodyStartsWith: `{"message":"/ns/ns-unit-01/clusters/cb-cluster/nodes/sample not found"}`,
		}
		EchoTest(t, tc)

		tc = TestCases{
			Name:                 "remove node",
			EchoFunc:             "RemoveNode",
			HttpMethod:           http.MethodDelete,
			WhenURL:              "/ladybug/ns/:namespace/clusters/:cluster/nodes/:node",
			GivenQueryParams:     "",
			GivenParaNames:       []string{"namespace", "cluster", "node"},
			GivenParaVals:        []string{"ns-unit-01", "cb-cluster", "sample"},
			GivenPostData:        "",
			ExpectStatus:         http.StatusBadRequest,
			ExpectBodyStartsWith: `{"message":"cluster info not found"}`,
		}
		EchoTest(t, tc)

		tc = TestCases{
			Name:                 "delete cluster",
			EchoFunc:             "DeleteCluster",
			HttpMethod:           http.MethodDelete,
			WhenURL:              "/ladybug/ns/:namespace/clusters/:cluster",
			GivenQueryParams:     "",
			GivenParaNames:       []string{"namespace", "cluster"},
			GivenParaVals:        []string{"ns-unit-01", "cb-cluster"},
			GivenPostData:        "",
			ExpectStatus:         http.StatusOK,
			ExpectBodyStartsWith: `{"kind":"Status","code":1,"message":"cluster cb-cluster not found"}`,
		}
		EchoTest(t, tc)

		TearDownForRest()
	})

}
