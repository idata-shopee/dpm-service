package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type NA struct {
	Host string
	Port int
}

type NAConf struct {
	NADeployCnfPath  string
	NAMachineCnfPath string
	NAs              []NA
}

func getNAsStr(NAs []NA) string {
	var texts []string
	for _, na := range NAs {
		texts = append(texts, fmt.Sprintf("%s:%d", na.Host, na.Port))
	}

	return strings.Join(texts, ";")
}

func DeployNAProcess(dpmConf DPMConf, naConf NAConf, na NA) error {
	var project = fmt.Sprintf("na_%s_%d", na.Host, na.Port)

	cmd := exec.Command(
		"ideploy",
		"--onlineType", dpmConf.OnlineType,
		"--config", naConf.NADeployCnfPath,
		"--machineConfig", naConf.NAMachineCnfPath,
		"--host", na.Host,
		"--deployDir", path.Join(dpmConf.TargetDir, project),
		"--srcDir", path.Join(dpmConf.SrcDir, "na"),
		"--project", project,
		"--remoteDir", dpmConf.RemoteRoot+"/"+project,
		"--NAPort", strconv.Itoa(na.Port),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
