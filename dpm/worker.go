package dpm

import (
	"log"
	"os"
	"os/exec"
	"path"
)

type Worker struct {
	ServiceType      string
	Name             string
	DcyTplPath       string
	DcyTplConfigPath string
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

func (workerConf *WorkerConf) GetWorkers(only string) []Worker {
	// deploy workers
	var workers []Worker
	if only != "" { // filter by only
		for _, worker := range workerConf.Workers {
			if worker.ServiceType == only {
				workers = append(workers, worker)
				break
			}
		}
	} else {
		workers = workerConf.Workers
	}
	return workers
}

// deploy worker to target machine
func DeployWorkerProcess(worker Worker, machine Machine, dpmConf DPMConf, workerConf WorkerConf, naConf NAConf) error {
	project := worker.Name
	if project == "" {
		project = worker.ServiceType
	}

	dcyTplPath, dcyTplConfigPath := worker.DcyTplPath, worker.DcyTplConfigPath
	if dcyTplPath == "" {
		dcyTplPath = "./dcy.tpl"
	}
	if dcyTplConfigPath == "" {
		dcyTplConfigPath = "./dcy-cnf.json"
	}

	cmd := exec.Command(
		"ideploy",
		"--onlineType", dpmConf.OnlineType,
		"--config", workerConf.WorkerDeployCnfPath,
		"--machineConfig", workerConf.WorkerMachineCnfPath,
		"--host", machine.Host,
		"--deployDir", path.Join(dpmConf.TargetDir, worker.ServiceType),
		"--srcDir", path.Join(dpmConf.SrcDir, worker.ServiceType),
		"--project", project,
		"--NAs", getNAsStr(naConf.NAs),
		"--dcyTplPath", dcyTplPath,
		"--dcyTplConfigPath", dcyTplConfigPath,
		"--remoteDir", dpmConf.RemoteRoot+"/"+worker.ServiceType,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func DeployWorkers(dpmConf DPMConf, workerConf WorkerConf, naConf NAConf) error {
	// deploy workers
	var workers = workerConf.GetWorkers(dpmConf.Only)

	// deploy worker to each machine
	for _, worker := range workers {
		for _, machine := range workerConf.Machines {
			log.Printf("start deploy %s to %s\n", worker.ServiceType, machine.Host)
			err := DeployWorkerProcess(worker, machine, dpmConf, workerConf, naConf)

			if err != nil {
				return err
			}
		}
	}

	return nil
}
