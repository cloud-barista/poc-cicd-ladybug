package cliscenario

import (
	"fmt"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestCliFullYaml(t *testing.T) {
	t.Run("command full yaml in/out test for mock driver", func(t *testing.T) {
		SetUpForCli()

		tc := TestCases{
			Name:                "healthy",
			CmdArgs:             []string{"healthy", "--config", "../conf/grpc_conf.yaml", "-i", "yaml", "-o", "yaml"},
			ExpectResStartsWith: `message: cb-barista cb-ladybug`,
		}
		LadybugCmdTest(t, tc)

		tc = TestCases{
			Name: "create cluster",
			CmdArgs: []string{"cluster", "create", "--config", "../conf/grpc_conf.yaml", "-i", "yaml", "-o", "yaml", "-d", `
"namespace": "ns-unit-01"
"ReqInfo":
  "name": "cb-cluster"
  "controlPlane":
    - "connection": "mock-unit-config01"
      "count": 1
      "spec": "mock-vmspec-01"
  "worker":
    - "connection": "mock-unit-config01"
      "count": 1
      "spec": "mock-vmspec-01"
  "config":
    "kubernetes":
      "networkCni": "kilo"
      "podCidr": "10.244.0.0/16"
      "serviceCidr": "10.96.0.0/12"
      "serviceDnsDomain": "cluster.local"
`,
			},
			ExpectResStartsWith: `name: cb-cluster
kind: Cluster`,
		}
		LadybugCmdTest(t, tc)

		tc = TestCases{
			Name:    "list cluster",
			CmdArgs: []string{"cluster", "list", "--config", "../conf/grpc_conf.yaml", "-i", "yaml", "-o", "yaml", "--ns", "ns-unit-01"},
			ExpectResStartsWith: `kind: ClusterList
items:
- name: cb-cluster`,
		}
		LadybugCmdTest(t, tc)

		tc = TestCases{
			Name:    "get cluster",
			CmdArgs: []string{"cluster", "get", "--config", "../conf/grpc_conf.yaml", "-i", "yaml", "-o", "yaml", "--ns", "ns-unit-01", "--cluster", "cb-cluster"},
			ExpectResStartsWith: `name: cb-cluster
kind: Cluster
status: completed`,
		}
		LadybugCmdTest(t, tc)

		tc = TestCases{
			Name: "add node",
			CmdArgs: []string{"node", "add", "--config", "../conf/grpc_conf.yaml", "-i", "yaml", "-o", "yaml", "-d", `
"namespace": "ns-unit-01"
"cluster": "cb-cluster"
"ReqInfo":
  "worker":
    - "connection": "mock-unit-config01"
      "count": 1
      "spec": "mock-vmspec-01"
`,
			},
			ExpectResStartsWith: `kind: NodeList
items:
- name:`,
		}
		LadybugCmdTest(t, tc)

		tc = TestCases{
			Name:    "list node",
			CmdArgs: []string{"node", "list", "--config", "../conf/grpc_conf.yaml", "-i", "yaml", "-o", "yaml", "--ns", "ns-unit-01", "--cluster", "cb-cluster"},
			ExpectResStartsWith: `kind: NodeList
items:
- name:`,
		}
		res, err := LadybugCmdTest(t, tc)
		nodeName := "undefined"
		if err == nil {
			nodeList := make(map[interface{}]interface{})
			err = yaml.Unmarshal([]byte(res), &nodeList)
			if err == nil {
				for _, m := range nodeList["items"].([]interface{}) {
					nodeInfo := m.(map[interface{}]interface{})
					nodeName = fmt.Sprintf("%v", nodeInfo["name"])
				}
			}
		}

		tc = TestCases{
			Name:    "get node",
			CmdArgs: []string{"node", "get", "--config", "../conf/grpc_conf.yaml", "-i", "yaml", "-o", "yaml", "--ns", "ns-unit-01", "--cluster", "cb-cluster", "--node", nodeName},
			ExpectResStartsWith: `name: ` + nodeName + `
kind: Node
credential:`,
		}
		LadybugCmdTest(t, tc)

		tc = TestCases{
			Name:    "remove node",
			CmdArgs: []string{"node", "remove", "--config", "../conf/grpc_conf.yaml", "-i", "yaml", "-o", "yaml", "--ns", "ns-unit-01", "--cluster", "cb-cluster", "--node", nodeName},
			ExpectResStartsWith: `kind: Status
code: 1
message: success`,
		}
		LadybugCmdTest(t, tc)

		tc = TestCases{
			Name:    "delete cluster",
			CmdArgs: []string{"cluster", "delete", "--config", "../conf/grpc_conf.yaml", "-i", "yaml", "-o", "yaml", "--ns", "ns-unit-01", "--cluster", "cb-cluster"},
			ExpectResStartsWith: `kind: Status
code: 1
message: cluster cb-cluster has been deleted`,
		}
		LadybugCmdTest(t, tc)

		TearDownForCli()
	})
}
