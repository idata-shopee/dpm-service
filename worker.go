package main

import (
	"os"
	"os/exec"
	"path"
)

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

// deploy worker to target machine
func DeployWorkerProcess(worker Worker, machine Machine, dpmConf DPMConf, workerConf WorkerConf, naConf NAConf) error {
	var project = worker.ServiceType

	cmd := exec.Command(
		"ideploy",
		"--onlineType", dpmConf.OnlineType,
		"--config", workerConf.WorkerDeployCnfPath,
		"--machineConfig", workerConf.WorkerMachineCnfPath,
		"--host", machine.Host,
		"--deployDir", path.Join(dpmConf.TargetDir, project),
		"--srcDir", path.Join(dpmConf.SrcDir, project),
		"--project", project,
		"--NAs", getNAsStr(naConf.NAs),
		"--remoteDir", dpmConf.RemoteRoot+"/"+project,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}