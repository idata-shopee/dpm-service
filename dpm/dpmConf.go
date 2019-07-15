package dpm

type DPMConf struct {
	RemoteRoot     string // remote root dir
	OnlineType     string
	NAConfPath     string // na config path
	WorkerConfPath string // worker config path
	TargetDir      string // target dir which used to store stages at dpm local
	SrcDir         string // source dir which store source code
	Only           string // only deploy target worker
}
