package goscenario

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGoFullDocArg(t *testing.T) {
	t.Run("go api full test for mock driver by doccument args style", func(t *testing.T) {
		SetUpForGrpc()

		tc := TestCases{
			Name:                "healthy",
			Instance:            McarApi,
			Method:              "Healthy",
			Args:                nil,
			ExpectResStartsWith: `{"message":"cb-barista cb-ladybug"}`,
		}
		MethodTest(t, tc)

		tc = TestCases{
			Name:     "create cluster",
			Instance: McarApi,
			Method:   "CreateCluster",
			Args: []interface{}{
				`{
					"namespace":  "ns-unit-01",
					"ReqInfo": {
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
						}
				}`,
			},
			ExpectResStartsWith: `{"name":"cb-cluster","kind":"Cluster"`,
		}
		MethodTest(t, tc)

		tc = TestCases{
			Name:     "list cluster",
			Instance: McarApi,
			Method:   "ListCluster",
			Args: []interface{}{
				`{"namespace": "ns-unit-01"}`,
			},
			ExpectResStartsWith: `{"kind":"ClusterList"`,
		}
		MethodTest(t, tc)

		tc = TestCases{
			Name:     "get cluster",
			Instance: McarApi,
			Method:   "GetCluster",
			Args: []interface{}{
				`{"namespace": "ns-unit-01", "cluster": "cb-cluster"}`,
			},
			ExpectResStartsWith: `{"name":"cb-cluster","kind":"Cluster"`,
		}
		MethodTest(t, tc)

		tc = TestCases{
			Name:     "add node",
			Instance: McarApi,
			Method:   "AddNode",
			Args: []interface{}{
				`{
					"namespace":  "ns-unit-01",
					"cluster":  "cb-cluster",
					"ReqInfo": {
							"worker": [{
								"connection": "mock-unit-config01",
								"count": 1,
								"spec": "mock-vmspec-01"
							}]
						}
				}`,
			},
			ExpectResStartsWith: `{"kind":"NodeList","items":[`,
		}
		MethodTest(t, tc)

		tc = TestCases{
			Name:     "list node",
			Instance: McarApi,
			Method:   "ListNode",
			Args: []interface{}{
				`{"namespace": "ns-unit-01", "cluster": "cb-cluster"}`,
			},
			ExpectResStartsWith: `{"kind":"NodeList","items":[`,
		}
		res, err := MethodTest(t, tc)
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
			Name:     "get node",
			Instance: McarApi,
			Method:   "GetNode",
			Args: []interface{}{
				`{"namespace": "ns-unit-01", "cluster": "cb-cluster", "node": "` + nodeName + `"}`,
			},
			ExpectResStartsWith: `{"name":"` + nodeName + `","kind":"Node"`,
		}
		MethodTest(t, tc)

		tc = TestCases{
			Name:     "remove node",
			Instance: McarApi,
			Method:   "RemoveNode",
			Args: []interface{}{
				`{"namespace": "ns-unit-01", "cluster": "cb-cluster", "node": "` + nodeName + `"}`,
			},
			ExpectResStartsWith: `{"kind":"Status","code":1,"message":"success"}`,
		}
		MethodTest(t, tc)

		tc = TestCases{
			Name:     "delete cluster",
			Instance: McarApi,
			Method:   "DeleteCluster",
			Args: []interface{}{
				`{"namespace": "ns-unit-01", "cluster": "cb-cluster"}`,
			},
			ExpectResStartsWith: `{"kind":"Status","code":1,"message":"cluster cb-cluster has been deleted"}`,
		}
		MethodTest(t, tc)

		TearDownForGrpc()
	})

}
