package dpm

import (
	"fmt"
	"github.com/lock-free/obrero"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type NAConf struct {
	NADeployCnfPath  string
	NAMachineCnfPath string
	NAs              []obrero.NA
}

func getNAsStr(NAs []obrero.NA) string {
	var texts []string
	for _, na := range NAs {
		texts = append(texts, fmt.Sprintf("%s:%d", na.Host, na.Port))
	}

	return strings.Join(texts, ";")
}

func DeployNAProcess(na obrero.NA, dpmConf DPMConf, naConf NAConf) error {
	var project = fmt.Sprintf("na_%s_%d", na.Host, na.Port)

	cmd := exec.Command(
		"ideploy",
		"--onlineType", dpmConf.OnlineType,
		"--config", naConf.NADeployCnfPath,
		"--machineConfig", naConf.NAMachineCnfPath,
		"--host", na.Host,
		"--deployDir", path.Join(dpmConf.TargetDir, project),
		"--srcDir", path.Join(dpmConf.SrcDir, "na_service"),
		"--project", project,
		"--remoteDir", dpmConf.RemoteRoot+"/"+project,
		"--NAPort", strconv.Itoa(na.Port),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func DeployNAs(dpmConf DPMConf, naConf NAConf) error {
	// deploy NAs
	if dpmConf.Only == "" || dpmConf.Only == "na" {
		for _, na := range naConf.NAs {
			err := DeployNAProcess(na, dpmConf, naConf)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
