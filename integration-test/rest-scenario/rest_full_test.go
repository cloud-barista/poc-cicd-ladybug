package restscenario

import (
	"encoding/json"
	"fmt"
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

		tc = TestCases{
			Name:             "create cluster",
			EchoFunc:         "CreateCluster",
			HttpMethod:       http.MethodPost,
			WhenURL:          "/ladybug/ns/:namespace/clusters",
			GivenQueryParams: "",
			GivenParaNames:   []string{"namespace"},
			GivenParaVals:    []string{"ns-unit-01"},
			GivenPostData: `{
													"name": "cb-cluster",
													"controlPlane" : [{
														"connection": "mock-unit-config01",
														"count": 1,
														"spec": "mock-vmspec-01"
													}],
													"worker": [{
														"connection": "mock-unit-config01",
														"count": 1,
														"spec": "mock-vmspec-01"
													}],
													"config": {
														"kubernetes": {
															"networkCni": "kilo",
															"podCidr": "10.244.0.0/16",
															"serviceCidr": "10.96.0.0/12",
															"serviceDnsDomain": "cluster.local"
														}
													}
											}`,
			ExpectStatus:         http.StatusOK,
			ExpectBodyStartsWith: `{"name":"cb-cluster","kind":"Cluster"`,
		}
		EchoTest(t, tc)

		tc = TestCases{
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
			ExpectStatus:         http.StatusOK,
			ExpectBodyStartsWith: `{"name":"cb-cluster","kind":"Cluster"`,
		}
		EchoTest(t, tc)

		tc = TestCases{
			Name:             "add node",
			EchoFunc:         "AddNode",
			HttpMethod:       http.MethodPost,
			WhenURL:          "/ladybug/ns/:namespace/clusters/:cluster/nodes",
			GivenQueryParams: "",
			GivenParaNames:   []string{"namespace", "cluster"},
			GivenParaVals:    []string{"ns-unit-01", "cb-cluster"},
			GivenPostData: `{
													"worker": [{
														"connection": "mock-unit-config01",
														"count": 1,
														"spec": "mock-vmspec-01"
													}]
											}`,
			ExpectStatus:         http.StatusOK,
			ExpectBodyStartsWith: `{"kind":"NodeList","items":[`,
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
		res, err := EchoTest(t, tc)
		nodeName := "undefined"
		if err == nil {
			nodeList := make(map[string]interface{})
			err = json.Unmarshal([]byte(res), &nodeList)
			if err == nil {
				for _, m := range nodeList["items"].([]interface{}) {
					nodeInfo := m.(map[string]interface{})
					nodeName = fmt.Sprintf("%v", nodeInfo["name"])
				}
			}
		}

		tc = TestCases{
			Name:                 "get node",
			EchoFunc:             "GetNode",
			HttpMethod:           http.MethodGet,
			WhenURL:              "/ladybug/ns/:namespace/clusters/:cluster/nodes/:node",
			GivenQueryParams:     "",
			GivenParaNames:       []string{"namespace", "cluster", "node"},
			GivenParaVals:        []string{"ns-unit-01", "cb-cluster", nodeName},
			GivenPostData:        "",
			ExpectStatus:         http.StatusOK,
			ExpectBodyStartsWith: `{"name":"` + nodeName + `","kind":"Node"`,
		}
		EchoTest(t, tc)

		tc = TestCases{
			Name:                 "remove node",
			EchoFunc:             "RemoveNode",
			HttpMethod:           http.MethodDelete,
			WhenURL:              "/ladybug/ns/:namespace/clusters/:cluster/nodes/:node",
			GivenQueryParams:     "",
			GivenParaNames:       []string{"namespace", "cluster", "node"},
			GivenParaVals:        []string{"ns-unit-01", "cb-cluster", nodeName},
			GivenPostData:        "",
			ExpectStatus:         http.StatusOK,
			ExpectBodyStartsWith: `{"kind":"Status","code":1,"message":"success"}`,
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
			ExpectBodyStartsWith: `{"kind":"Status","code":1,"message":"cluster cb-cluster has been deleted"}`,
		}
		EchoTest(t, tc)

		TearDownForRest()
	})

}
