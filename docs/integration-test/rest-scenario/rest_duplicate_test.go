package restscenario

import (
	"net/http"
	"testing"
)

func TestRestDuplicate(t *testing.T) {
	t.Run("rest api duplicate test for mock driver", func(t *testing.T) {
		SetUpForRest()

		tc := TestCases{
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
			Name:             "create cluster(duplicate)",
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
			ExpectStatus:         http.StatusBadRequest,
			ExpectBodyStartsWith: `{"message":"MCIS already exists"}`,
		}
		EchoTest(t, tc)

		TearDownForRest()
	})

}
