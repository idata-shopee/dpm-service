package main

import (
	"github.com/idata-shopee/gopcp"
	"github.com/idata-shopee/gopcp_rpc"
	"github.com/idata-shopee/gopcp_stream"
	"log"
	"sync"
	"time"
)

const DPMConfPath = "/data/config.json"

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
	if err != nil {
		panic(err)
	}

	naConf := NAConf{}
	err = ReadJson(dpmConf.NAConfPath, &naConf)
	if err != nil {
		panic(err)
	}

	workerConf := WorkerConf{}
	err = ReadJson(dpmConf.WorkerConfPath, &workerConf)
	if err != nil {
		panic(err)
	}

	// deploy NAs
	for _, na := range naConf.NAs {
		err = DeployNAProcess(na, dpmConf, naConf)
		if err != nil {
			panic(err)
		}
	}

	// deploy worker to connect each NA
	// principle: one worker one machine, one worker connect to all NAs.
	for _, worker := range workerConf.Workers {
		for _, machine := range workerConf.Machines {
			err = DeployWorkerProcess(worker, machine, dpmConf, workerConf, naConf)

			if err != nil {
				panic(err)
			}
		}
	}

	// connects NAs
	for _, na := range naConf.NAs {
		MaintainConnectionWithNA(na.Host, na.Port)
	}

	// blocking forever
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func MaintainConnectionWithNA(NAHost string, NAPort int) {
	log.Printf("try to connect to NA %s:%d\n", NAHost, NAPort)
	_, err := gopcp_rpc.GetPCPRPCClient(NAHost, NAPort, generateSandbox, func(err error) {
		log.Printf("connection to NA %s:%d broken, error is %v\n", NAHost, NAPort, err)
		time.Sleep(2 * time.Second)
		MaintainConnectionWithNA(NAHost, NAPort)
	})

	if err != nil {
		log.Printf("fail to connect to NA %s:%d\n", NAHost, NAPort)
		time.Sleep(2 * time.Second)
		MaintainConnectionWithNA(NAHost, NAPort)
	} else {
		log.Printf("connected to NA %s:%d\n", NAHost, NAPort)
	}
}

// define sandbox for dpm service
func generateSandbox(*gopcp_stream.StreamServer) *gopcp.Sandbox {
	return gopcp.GetSandbox(map[string]*gopcp.BoxFunc{
		// define service type
		"getServiceType": gopcp.ToSandboxFun(func(args []interface{}, attachment interface{}, pcpServer *gopcp.PcpServer) (interface{}, error) {
			return "dpm", nil
		}),
	})
}
