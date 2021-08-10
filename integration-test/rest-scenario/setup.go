package restscenario

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"os/exec"

	"bou.ke/monkey"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"github.com/cloud-barista/poc-cicd-ladybug/src/core/service"
	lb_conf "github.com/cloud-barista/poc-cicd-ladybug/src/utils/config"
	"github.com/cloud-barista/poc-cicd-ladybug/src/utils/lang"

	sshrun "github.com/cloud-barista/cb-spider/cloud-control-manager/vm-ssh"
	cbstore "github.com/cloud-barista/cb-store"
)

type TestCases struct {
	Name                 string
	EchoFunc             string
	HttpMethod           string
	WhenURL              string
	GivenQueryParams     string
	GivenParaNames       []string
	GivenParaVals        []string
	GivenPostData        string
	ExpectStatus         int
	ExpectBodyStartsWith string
	ExpectBodyContains   string
}

var (
	holdStdout *os.File = nil
)

func init() {
	logrus.SetLevel(logrus.ErrorLevel)

	if flag.Lookup("app-root") == nil {
		lb_conf.Config.AppRootPath = flag.String("app-root", lang.NVL(os.Getenv("APP_ROOT"), ""), "application root path")
		lb_conf.Config.RootURL = flag.String("root-url", lang.NVL(os.Getenv("BASE_URL"), "/ladybug"), "root url")
		lb_conf.Config.SpiderUrl = flag.String("spider-url", lang.NVL(os.Getenv("SPIDER_URL"), "http://localhost:1024/spider"), "cb-spider service end-point url")
		lb_conf.Config.TumblebugUrl = flag.String("tumblebug-url", lang.NVL(os.Getenv("TUMBLEBUG_URL"), "http://localhost:1323/tumblebug"), "cb-tumblebug service end-point url")
		lb_conf.Config.Username = flag.String("basic-auth-username", lang.NVL(os.Getenv("BASIC_AUTH_USERNAME"), "default"), "rest-api basic auth usernmae")
		lb_conf.Config.Password = flag.String("basic-auth-password", lang.NVL(os.Getenv("BASIC_AUTH_PASSWORD"), "default"), "rest-api basic auth password")
		lb_conf.Config.LoglevelHTTP = flag.Bool("log-http", os.Getenv("LOG_HTTP") == "true", "The logging http data")
	}
}

func SetUpForRest() {

	holdStdout = os.Stdout
	os.Stdout = os.NewFile(0, os.DevNull)

	cbstore.InitData()

	/**
	** Backend Server Setup
	**/
	client := resty.New().SetCloseConnection(true)

	cmd := exec.Command("./stop.sh")
	cmd.Dir = "../backend"
	cmd.Run()

	cmd = exec.Command("./start.sh")
	cmd.Dir = "../backend"
	cmd.Start()

	backendFlg := false
	for i := 0; i < 10; i++ {
		//fmt.Printf("backend server waiting... \n")
		time.Sleep(time.Second * 2)

		_, err := client.R().
			Get("http://localhost:31024/spider/")

		if err == nil {
			backendFlg = true
			break
		}
	}

	if !backendFlg {
		log.Fatalf("backend server failed...\n")
	}

	client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"DriverName":"mock-unit-driver01","ProviderName":"MOCK", "DriverLibFileName":"mock-driver-v1.0.so"}`).
		Post("http://localhost:31024/spider/driver")

	client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"CredentialName":"mock-unit-credential01","ProviderName":"MOCK", "KeyValueInfoList": [{"Key":"MockName", "Value":"mock_unit_name00"}]}`).
		Post("http://localhost:31024/spider/credential")

	client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"RegionName":"mock-unit-region01","ProviderName":"MOCK", "KeyValueInfoList": [{"Key":"Region", "Value":"default"}]}`).
		Post("http://localhost:31024/spider/region")

	client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"ConfigName":"mock-unit-config01","ProviderName":"MOCK", "DriverName":"mock-unit-driver01", "CredentialName":"mock-unit-credential01", "RegionName":"mock-unit-region01"}`).
		Post("http://localhost:31024/spider/connectionconfig")

	auth := "default:default"
	encAuth := base64.StdEncoding.EncodeToString([]byte(auth))

	client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Basic "+encAuth).
		SetBody(`{"name":"ns-unit-01","description":"NameSpace for General Testing"}`).
		Post("http://localhost:31323/tumblebug/ns")

	/**
	** Ladybug Env Setup
	**/
	flag.Parse()

	/**
	** Function Patch for Testing
	**/
	monkey.Patch(sshrun.SSHRun, func(sshInfo sshrun.SSHInfo, cmd string) (string, error) {
		//fmt.Printf("ssh cmd : %s\n", cmd)

		hostName := "ladybugnode"
		if cmd == "/bin/hostname" {
			return hostName, nil
		}
		if strings.Contains(cmd, "cd /tmp;./bootstrap.sh") {
			return "kubectl set on hold", nil
		}
		if strings.Contains(cmd, "cd /tmp;./k8s-init.sh") {
			return "Your Kubernetes control-plane has initialized successfully", nil
		}
		if strings.Contains(cmd, "sudo kubectl drain") {
			return fmt.Sprintf("node/%s drained", hostName), nil
		}
		if strings.Contains(cmd, "sudo kubectl delete node") {
			return "deleted", nil
		}
		if strings.Contains(cmd, "sudo") {
			return "This node has joined the cluster, This node has joined the cluster", nil
		}
		return cmd + " success", nil
	})

	monkey.Patch(sshrun.SSHCopy, func(sshInfo sshrun.SSHInfo, sourcePath string, remotePath string) error {
		return nil
	})

	monkey.Patch(service.GetCSPName, func(providerName string) (lb_conf.CSP, error) {
		return "mock", nil
	})

	monkey.Patch(service.GetVmImageId, func(csp lb_conf.CSP, configName string) (string, error) {
		return "mock-vmimage-01", nil
	})
}

func TearDownForRest() {
	cmd := exec.Command("./stop.sh")
	cmd.Dir = "../backend"
	cmd.Run()

	cbstore.InitData()

	os.Stdout = holdStdout
}
