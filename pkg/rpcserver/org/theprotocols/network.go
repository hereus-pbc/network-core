package rpc_org_theprotocols

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/hereus-pbc/network-core/pkg/httpserver/helpers"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
	"github.com/hereus-pbc/network-core/pkg/types"
)

func NetworkRpc() *helpers.RpcFunctionHandlerNoArgumentsNoSession {
	return &helpers.RpcFunctionHandlerNoArgumentsNoSession{
		Handler: func(kernel interfaces.Kernel) (interface{}, error) {
			return Network(kernel)
		},
	}
}

func Network(kernel interfaces.Kernel) (types.Network, error) {
	helpUsername, err := kernel.ReadConfigString("help_username")
	if err != nil {
		return types.Network{}, fmt.Errorf("failed to read help_username config: %v", err)
	}
	var subscriptionPlans struct {
		Plans []types.SubscriptionPlan `bson:"config_value"`
	}
	err = kernel.ReadConfig("subscription_plans", &subscriptionPlans)
	if err != nil {
		return types.Network{}, fmt.Errorf("failed to read subscription_plans config: %v", err)
	}
	osName := runtime.GOOS
	osVersion := ""
	if runtime.GOOS == "darwin" {
		osName = "macOS"
		data, err := os.ReadFile("/System/Library/CoreServices/SystemVersion.plist")
		if err != nil {
			return types.Network{}, fmt.Errorf("failed to read macOS version file: %v", err)
		}
		osVersionLineNumber := -1
		for lineNo, line := range strings.Split(string(data), "\n") {
			if strings.Contains(line, "<key>ProductVersion</key>") {
				osVersionLineNumber = lineNo + 1
				break
			}
		}
		if osVersionLineNumber != -1 && osVersionLineNumber < len(strings.Split(string(data), "\n")) {
			osVersion = strings.Trim(strings.Split(strings.Split(string(data), "\n")[osVersionLineNumber], "<string>")[1], "</string>")
		}
	} else if runtime.GOOS == "linux" {
		data, _ := os.ReadFile("/etc/os-release")
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "NAME=") {
				osName = strings.Trim(line[5:], `"`)
			} else if strings.HasPrefix(line, "VERSION_ID=") {
				osVersion = strings.Trim(line[11:], `"`)
			}
		}
	}
	domain, err := kernel.ReadConfigString("domain")
	if err != nil {
		return types.Network{}, fmt.Errorf("failed to read domain config: %v", err)
	}
	return types.Network{
		HelpUsername:      helpUsername,
		SubscriptionPlans: subscriptionPlans.Plans,
		OS: types.NetworkOS{
			Arch:    runtime.GOARCH,
			Family:  runtime.GOOS,
			Name:    osName,
			Version: osVersion,
		},
		Rules: types.NetworkRules{
			NewAccountsAllowed: false,
		},
		Software: types.NetworkSoftware{
			Name:    "Network",
			Version: kernel.GetSoftwareVersion(),
			Build:   kernel.GetSoftwareBuild(),
			Channel: kernel.GetSoftwareVersionChannel(),
		},
		TheProtocolsVersion: 4.0,
		AuthorizationUrl:    fmt.Sprintf("https://%s/flow/authorize", domain),
	}, nil
}
