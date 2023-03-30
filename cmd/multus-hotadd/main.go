// Copyright (c) 2022 Tomofumi Hayashi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This is a "Multi-plugin".The delegate concept referred from CNI project
// It reads other plugin netconf, and then invoke them, e.g.
// flannel or sriov plugin.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	cnitypes "github.com/containernetworking/cni/pkg/types"
	"gopkg.in/k8snetworkplumbingwg/multus-cni.v3/pkg/server/api"
)

func main() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	mode := ""
	flag.StringVar(&mode, "mode", "", "mode")
	pod := ""
	flag.StringVar(&pod, "pod", "", "pod namespace/name")
	cniFile := ""
	flag.StringVar(&cniFile, "file", "", "cni config")
	containerID := ""
	flag.StringVar(&containerID, "containerid", "", "container id")
	netns := ""
	flag.StringVar(&netns, "netns", "", "netns")
	ifname := ""
	flag.StringVar(&ifname, "ifname", "", "ifname")
	podUID := ""
	flag.StringVar(&podUID, "poduid", "", "poduid")
	ips := ""
	flag.StringVar(&ips, "ips", "", "ips, comma separated")

	flag.Parse()

	if mode == "" || pod == "" || cniFile == "" || containerID == "" || netns == "" || ifname == "" || podUID == "" {
		fmt.Fprintf(os.Stderr, "option is not enough")
		os.Exit(1)
	}

	podNamespaceName := strings.Split(pod, "/")
	if len(podNamespaceName) != 2 {
		fmt.Fprintf(os.Stderr, "pod should be <podnamespace>/<pod>: %s", pod)
		os.Exit(1)
	}
	podNamespace := podNamespaceName[0]
	podName := podNamespaceName[1]

	data, err := ioutil.ReadFile(cniFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err at readfile: %v", err)
		os.Exit(2)
	}

	var interfaceAttributes *api.DelegateInterfaceAttributes
	if ips != "" {
		podIPs := strings.Split(ips, ",")
		interfaceAttributes = &api.DelegateInterfaceAttributes{
			IPRequest: podIPs,
		}
	}

	// From here, multus-daemon delegate API is invoked to add/del given CNI config
	req := api.CreateDelegateRequest(mode, containerID, netns, ifname, podNamespace, podName, podUID, data, interfaceAttributes)
	body, err := api.DoCNI(api.GetAPIEndpoint(api.MultusDelegateAPIEndpoint), req, api.SocketPath("/run/multus"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error at DoCNI(): %v", err)
	}

	resp := &api.Response{}
	if len(body) != 0 {
		if err = json.Unmarshal(body, resp); err != nil {
			fmt.Fprintf(os.Stderr, "Unmarshal failed: %v : %s", err, string(body))
			os.Exit(3)
		}
	}

	if resp.Result != nil {
		err = cnitypes.PrintResult(resp.Result, "1.0.0")
		if err != nil {
			fmt.Fprintf(os.Stderr, "printResult: %v", err)
		}
	}
}
