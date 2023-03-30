module github.com/s1061123/multus-delegate-sample

go 1.18

require (
	github.com/containernetworking/cni v1.1.2
	gopkg.in/k8snetworkplumbingwg/multus-cni.v3 v3.9.3
)

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace (
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2
	gopkg.in/k8snetworkplumbingwg/multus-cni.v3 => github.com/k8snetworkplumbingwg/multus-cni v0.0.0-20230302144144-c129e72779ce
)
