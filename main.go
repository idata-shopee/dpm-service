package main

import (
	"fmt"
	"github.com/idata-shopee/gopcp"
	"github.com/idata-shopee/gopcp_rpc"
	"github.com/idata-shopee/gopcp_service"
	"github.com/idata-shopee/gopcp_stream"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"
)

const workerDeployCnfPath = "/Users/yuer/workspaceforme/projects/insight-in-one/goProjects/src/github.com/idata-shopee/dpm-service/example/local/worker/deploy-cnf.json"
const workerCodeDir = "/Users/yuer/workspaceforme/projects/insight-in-one/goProjects/src/github.com/idata-shopee/dpm-service/example/local/code"
const onlineType = "staging"

type WorkerProcessBox struct {
	ServiceType string
	Index       int
	Host        string // where to start worker process
	NAHost      string // connect to which NA
	NAPort      int
}

func deployWorkerProcess(workerProcessBox WorkerProcessBox) error {
	var project = workerProcessBox.ServiceType + "-" + strconv.Itoa(workerProcessBox.Index)

	cmd := exec.Command("../../../../../thirdparty/flexdeploy/bin/ideploy",
		"--onlineType", onlineType,
		"--config", workerDeployCnfPath,
		"--host", workerProcessBox.Host,
		"--deployDir", path.Join(workerCodeDir, project),
		"--project", project,
		"--publishOnly",

		"--NAHost", workerProcessBox.NAHost,
		"--NAPort", strconv.Itoa(workerProcessBox.NAPort),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// TODO get worker hosts from deploy-cnf.json

func StartTcpServer(port int, mcCallTimeout time.Duration) error {
	return gopcp_service.StartTcpServer(port, func(streamServer *gopcp_stream.StreamServer) *gopcp.Sandbox {
		return gopcp.GetSandbox(map[string]*gopcp.BoxFunc{
			"startWorkerProcess": gopcp.ToSandboxFun(func(args []interface{}, attachment interface{}, pcpServer *gopcp.PcpServer) (interface{}, error) {
				return nil, deployWorkerProcess(WorkerProcessBox{
					ServiceType: "user",
					Index:       0,
					Host:        "yuer@local",
					NAHost:      "test.com",
					NAPort:      9988,
				})
			}),
		})
	}, func() *gopcp_rpc.ConnectionEvent {
		return &gopcp_rpc.ConnectionEvent{
			// on close of connection
			func(err error) {
			},
			// new connection
			func(pcpConnectionHandler *gopcp_rpc.PCPConnectionHandler) {
			},
		}
	})
}

func main() {
	StartTcpServer(7655, 2*time.Minute)
}
