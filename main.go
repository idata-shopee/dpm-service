package main

import (
	"encoding/json"
	"fmt"
	"github.com/idata-shopee/gopcp"
	"github.com/idata-shopee/gopcp_rpc"
	"github.com/idata-shopee/gopcp_stream"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

const DATA = "/data"
const DPMConfPath = DATA + "/config.json"

type DPMConf struct {
	RemoteRoot     string
	OnlineType     string
	NAConfPath     string
	WorkerConfPath string
	TargetDir      string
	SrcDir         string
}

type NA struct {
	Host string
	Port int
}

type NAConf struct {
	NADeployCnfPath  string
	NAMachineCnfPath string
	NAs              []NA
}

type Worker struct {
	ServiceType string
}

type Machine struct {
	Host string
}

type WorkerConf struct {
	WorkerDeployCnfPath  string
	WorkerMachineCnfPath string
	Workers              []Worker
	Machines             []Machine
}

type WorkerProcessBox struct {
	ServiceType          string
	WorkerDeployCnfPath  string
	WorkerMachineCnfPath string
	Host                 string // where to start worker process
	NAs                  []NA
}

func getNAsStr(NAs []NA) string {
	var texts []string
	for _, na := range NAs {
		texts = append(texts, fmt.Sprintf("%s:%d", na.Host, na.Port))
	}

	return strings.Join(texts, ";")
}

func deployWorkerProcess(dpmConf DPMConf, workerProcessBox WorkerProcessBox) error {
	var project = workerProcessBox.ServiceType

	cmd := exec.Command(
		"ideploy",
		"--onlineType", dpmConf.OnlineType,
		"--config", workerProcessBox.WorkerDeployCnfPath,
		"--machineConfig", workerProcessBox.WorkerMachineCnfPath,
		"--host", workerProcessBox.Host,
		"--deployDir", path.Join(dpmConf.TargetDir, project),
		"--srcDir", path.Join(dpmConf.SrcDir, project),
		"--project", project,
		"--NAs", getNAsStr(workerProcessBox.NAs),
		"--remoteDir", dpmConf.RemoteRoot+"/"+project,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func deployNAProcess(dpmConf DPMConf, NADeployCnfPath string, NAMachineCnfPath string, NAHost string, port int) error {
	var project = fmt.Sprintf("na_%s_%d", NAHost, port)

	cmd := exec.Command(
		"ideploy",
		"--onlineType", dpmConf.OnlineType,
		"--config", NADeployCnfPath,
		"--machineConfig", NAMachineCnfPath,
		"--host", NAHost,
		"--deployDir", path.Join(dpmConf.TargetDir, project),
		"--srcDir", path.Join(dpmConf.SrcDir, "na"),
		"--project", project,
		"--remoteDir", dpmConf.RemoteRoot+"/"+project,
		"--NAPort", strconv.Itoa(port),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// 1. dpm will try to connect to all NA nodes
//  (1) read NA config
//  (2) try to connect to NA.
//
// 2. update mc
//

func main() {
	var err error
	// read configs
	dpmConf := DPMConf{}
	err = ReadJson(DPMConfPath, &dpmConf)

	naConf := NAConf{}
	err = ReadJson(dpmConf.NAConfPath, &naConf)
	if err != nil {
		panic(err)
	}

	workerConf := WorkerConf{}
	err = ReadJson(dpmConf.WorkerConfPath, &workerConf)

	// deploy NAs
	for _, na := range naConf.NAs {
		err = deployNAProcess(dpmConf, naConf.NADeployCnfPath, naConf.NAMachineCnfPath, na.Host, na.Port)
		if err != nil {
			panic(err)
		}
	}

	// deploy worker to connect each NA
	// principle: one worker one machine, one worker connect to all NAs.
	for _, worker := range workerConf.Workers {
		for _, machine := range workerConf.Machines {
			err = deployWorkerProcess(dpmConf, WorkerProcessBox{
				ServiceType:          worker.ServiceType,
				WorkerDeployCnfPath:  workerConf.WorkerDeployCnfPath,
				WorkerMachineCnfPath: workerConf.WorkerMachineCnfPath,
				Host:                 machine.Host,
				NAs:                  naConf.NAs,
			})

			if err != nil {
				panic(err)
			}
		}
	}

	// connects NAs
	for _, na := range naConf.NAs {
		maintainConnectionWithNA(na.Host, na.Port)
	}

	// blocking forever
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func maintainConnectionWithNA(NAHost string, NAPort int) {
	log.Printf("try to connect to NA %s:%d\n", NAHost, NAPort)
	_, err := gopcp_rpc.GetPCPRPCClient(NAHost, NAPort, generateSandbox, func(err error) {
		log.Printf("connection to NA %s:%d broken, error is %v\n", NAHost, NAPort, err)
		time.Sleep(2 * time.Second)
		maintainConnectionWithNA(NAHost, NAPort)
	})

	if err != nil {
		log.Printf("fail to connect to NA %s:%d\n", NAHost, NAPort)
		time.Sleep(2 * time.Second)
		maintainConnectionWithNA(NAHost, NAPort)
	} else {
		log.Printf("connected to NA %s:%d\n", NAHost, NAPort)
	}
}

func generateSandbox(*gopcp_stream.StreamServer) *gopcp.Sandbox {
	return gopcp.GetSandbox(map[string]*gopcp.BoxFunc{
		// define service type
		"getServiceType": gopcp.ToSandboxFun(func(args []interface{}, attachment interface{}, pcpServer *gopcp.PcpServer) (interface{}, error) {
			return "dpm", nil
		}),
	})
}

func ReadJson(filePath string, f interface{}) error {
	source, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(source), f)
}
