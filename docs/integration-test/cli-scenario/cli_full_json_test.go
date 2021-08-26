package cliscenario

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestCliFullJson(t *testing.T) {
	t.Run("command full json in/out test for mock driver", func(t *testing.T) {
		SetUpForCli()

		tc := TestCases{
			Name:                "healthy",
			CmdArgs:             []string{"healthy", "--config", "../conf/grpc_conf.yaml", "-i", "json", "-o", "json"},
			ExpectResStartsWith: `{"message":"cb-barista cb-ladybug"}`,
		}
		LadybugCmdTest(t, tc)

		tc = TestCases{
			Name: "create cluster",
			CmdArgs: []string{"cluster", "create", "--config", "../conf/grpc_conf.yaml", "-i", "json", "-o", "json", "-d",
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
		LadybugCmdTest(t, tc)

		tc = TestCases{
			Name:                "list cluster",
			CmdArgs:             []string{"cluster", "list", "--config", "../conf/grpc_conf.yaml", "-i", "json", "-o", "json", "--ns", "ns-unit-01"},
			ExpectResStartsWith: `{"kind":"ClusterList"`,
		}
		LadybugCmdTest(t, tc)

		tc = TestCases{
			Name:                "get cluster",
			CmdArgs:             []string{"cluster", "get", "--config", "../conf/grpc_conf.yaml", "-i", "json", "-o", "json", "--ns", "ns-unit-01", "--cluster", "cb-cluster"},
			ExpectResStartsWith: `{"name":"cb-cluster","kind":"Cluster"`,
		}
		LadybugCmdTest(t, tc)

		tc = TestCases{
			Name: "add node",
			CmdArgs: []string{"node", "add", "--config", "../conf/grpc_conf.yaml", "-i", "json", "-o", "json", "-d",
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
		LadybugCmdTest(t, tc)

		tc = TestCases{
			Name:                "list node",
			CmdArgs:             []string{"node", "list", "--config", "../conf/grpc_conf.yaml", "-i", "json", "-o", "json", "--ns", "ns-unit-01", "--cluster", "cb-cluster"},
			ExpectResStartsWith: `{"kind":"NodeList","items":[`,
		}
		res, err := LadybugCmdTest(t, tc)
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
			Name:                "get node",
			CmdArgs:             []string{"node", "get", "--config", "../conf/grpc_conf.yaml", "-i", "json", "-o", "json", "--ns", "ns-unit-01", "--cluster", "cb-cluster", "--node", nodeName},
			ExpectResStartsWith: `{"name":"` + nodeName + `","kind":"Node"`,
		}
		LadybugCmdTest(t, tc)

		tc = TestCases{
			Name:                "remove node",
			CmdArgs:             []string{"node", "remove", "--config", "../conf/grpc_conf.yaml", "-i", "json", "-o", "json", "--ns", "ns-unit-01", "--cluster", "cb-cluster", "--node", nodeName},
			ExpectResStartsWith: `{"kind":"Status","code":1,"message":"success"}`,
		}
		LadybugCmdTest(t, tc)

		tc = TestCases{
			Name:                "delete cluster",
			CmdArgs:             []string{"cluster", "delete", "--config", "../conf/grpc_conf.yaml", "-i", "json", "-o", "json", "--ns", "ns-unit-01", "--cluster", "cb-cluster"},
			ExpectResStartsWith: `{"kind":"Status","code":1,"message":"cluster cb-cluster has been deleted"}`,
		}
		LadybugCmdTest(t, tc)

		TearDownForCli()
	})
}
