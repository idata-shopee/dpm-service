package main

import (
	"github.com/lock-free/dpm_service/dpm"
	"github.com/lock-free/gopcp"
	"github.com/lock-free/gopcp_stream"
	"github.com/lock-free/obrero"
	"time"
)

const DPMConfPath = "/data/config.json"

func main() {
	var err error

	// read configs
	dpmConf := dpm.DPMConf{}
	err = obrero.ReadJson(DPMConfPath, &dpmConf)
	if err != nil {
		panic(err)
	}

	naConf := dpm.NAConf{}
	err = obrero.ReadJson(dpmConf.NAConfPath, &naConf)
	if err != nil {
		panic(err)
	}

	workerConf := dpm.WorkerConf{}
	err = obrero.ReadJson(dpmConf.WorkerConfPath, &workerConf)
	if err != nil {
		panic(err)
	}

	// deploy NAs
	err = dpm.DeployNAs(dpmConf, naConf)
	if err != nil {
		panic(err)
	}

	// deploy workers
	err = dpm.DeployWorkers(dpmConf, workerConf, naConf)
	if err != nil {
		panic(err)
	}

	// dpm also is special service which will be a client to NAs.
	obrero.StartWorkerWithNAs(func(*gopcp_stream.StreamServer) *gopcp.Sandbox {
		return gopcp.GetSandbox(map[string]*gopcp.BoxFunc{
			// define service type
			"getServiceType": gopcp.ToSandboxFun(func(args []interface{}, attachment interface{}, pcpServer *gopcp.PcpServer) (interface{}, error) {
				return "dpm", nil
			}),
		})
	}, obrero.WorkerStartConf{
		PoolSize:            2,
		Duration:            20 * time.Second,
		RetryDuration:       20 * time.Second,
		NAGetClientMaxRetry: 3,
	}, naConf.NAs)

	obrero.RunForever()
}
