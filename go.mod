module github.com/s1061123/multus-delegate-sample

go 1.18

replace gopkg.in/k8snetworkplumbingwg/multus-cni.v3 => gopkg.in/s1061123/multus-cni.v3 v3.0.0-20221012071351-54e57c230bdf

require (
	github.com/containernetworking/cni v1.1.2
	gopkg.in/k8snetworkplumbingwg/multus-cni.v3 v3.9.2
)

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)
