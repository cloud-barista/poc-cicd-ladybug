package rgequalscenario

import (
	"net/http"
	"testing"

	gs "github.com/cloud-barista/poc-cicd-ladybug/integration-test/go-scenario"
	rs "github.com/cloud-barista/poc-cicd-ladybug/integration-test/rest-scenario"
)

func TestPrepareFull(t *testing.T) {
	t.Run("prepare full test for mock driver", func(t *testing.T) {
		gs.SetUpForGrpc()

		tc := rs.TestCases{
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
		rs.EchoTest(t, tc)

		tc = rs.TestCases{
			Name:                 "list cluster for rest",
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
		res, _ := rs.EchoTest(t, tc)
		RestPrepareResult["list cluster"] = res

		gtc := gs.TestCases{
			Name:     "list cluster for grpc",
			Instance: gs.McarApi,
			Method:   "ListCluster",
			Args: []interface{}{
				`{"namespace": "ns-unit-01"}`,
			},
			ExpectResStartsWith: `{"kind":"ClusterList"`,
		}
		res, _ = gs.MethodTest(t, gtc)
		GrpcPrepareResult["list cluster"] = res

		tc = rs.TestCases{
			Name:                 "list node for rest",
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
		res, _ = rs.EchoTest(t, tc)
		RestPrepareResult["list node"] = res

		gtc = gs.TestCases{
			Name:     "list node for grpc",
			Instance: gs.McarApi,
			Method:   "ListNode",
			Args: []interface{}{
				`{"namespace": "ns-unit-01", "cluster": "cb-cluster"}`,
			},
			ExpectResStartsWith: `{"kind":"NodeList","items":[`,
		}
		res, _ = gs.MethodTest(t, gtc)
		GrpcPrepareResult["list node"] = res

		gs.TearDownForGrpc()
	})

}
