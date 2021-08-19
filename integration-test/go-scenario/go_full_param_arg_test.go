package goscenario

import (
	"encoding/json"
	"fmt"
	"testing"

	api "github.com/cloud-barista/poc-cicd-ladybug/src/grpc-api/request"
)

func TestGoFullParamArg(t *testing.T) {
	t.Run("go api full test for mock driver by parameter args style", func(t *testing.T) {
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
			Method:   "CreateClusterByParam",
			Args: []interface{}{
				&api.ClusterCreateRequest{
					Namespace: "ns-unit-01",
					Item: api.ClusterReq{
						Name: "cb-cluster",
						ControlPlane: []api.NodeConfig{
							api.NodeConfig{Connection: "mock-unit-config01", Count: 1, Spec: "mock-vmspec-01"},
						},
						Worker: []api.NodeConfig{
							api.NodeConfig{Connection: "mock-unit-config01", Count: 1, Spec: "mock-vmspec-01"},
						},
						Config: api.Config{
							Kubernetes: api.Kubernetes{
								NetworkCni:       "kilo",
								PodCidr:          "10.244.0.0/16",
								ServiceCidr:      "10.96.0.0/12",
								ServiceDnsDomain: "cluster.local",
							},
						},
					},
				},
			},
			ExpectResStartsWith: `{"name":"cb-cluster","kind":"Cluster"`,
		}
		MethodTest(t, tc)

		tc = TestCases{
			Name:     "list cluster",
			Instance: McarApi,
			Method:   "ListClusterByParam",
			Args: []interface{}{
				"ns-unit-01",
			},
			ExpectResStartsWith: `{"kind":"ClusterList"`,
		}
		MethodTest(t, tc)

		tc = TestCases{
			Name:     "get cluster",
			Instance: McarApi,
			Method:   "GetClusterByParam",
			Args: []interface{}{
				"ns-unit-01",
				"cb-cluster",
			},
			ExpectResStartsWith: `{"name":"cb-cluster","kind":"Cluster"`,
		}
		MethodTest(t, tc)

		tc = TestCases{
			Name:     "add node",
			Instance: McarApi,
			Method:   "AddNodeByParam",
			Args: []interface{}{
				&api.NodeCreateRequest{
					Namespace: "ns-unit-01",
					Cluster:   "cb-cluster",
					Item: api.NodeReq{
						Worker: []api.NodeConfig{
							api.NodeConfig{Connection: "mock-unit-config01", Count: 1, Spec: "mock-vmspec-01"},
						},
					},
				},
			},
			ExpectResStartsWith: `{"kind":"NodeList","items":[`,
		}
		MethodTest(t, tc)

		tc = TestCases{
			Name:     "list node",
			Instance: McarApi,
			Method:   "ListNodeByParam",
			Args: []interface{}{
				"ns-unit-01",
				"cb-cluster",
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
			Method:   "GetNodeByParam",
			Args: []interface{}{
				"ns-unit-01",
				"cb-cluster",
				nodeName,
			},
			ExpectResStartsWith: `{"name":"` + nodeName + `","kind":"Node"`,
		}
		MethodTest(t, tc)

		tc = TestCases{
			Name:     "remove node",
			Instance: McarApi,
			Method:   "RemoveNodeByParam",
			Args: []interface{}{
				"ns-unit-01",
				"cb-cluster",
				nodeName,
			},
			ExpectResStartsWith: `{"kind":"Status","code":1,"message":"success"}`,
		}
		MethodTest(t, tc)

		tc = TestCases{
			Name:     "delete cluster",
			Instance: McarApi,
			Method:   "DeleteClusterByParam",
			Args: []interface{}{
				"ns-unit-01",
				"cb-cluster",
			},
			ExpectResStartsWith: `{"kind":"Status","code":1,"message":"cluster cb-cluster has been deleted"}`,
		}
		MethodTest(t, tc)

		TearDownForGrpc()
	})
}
