package Core

import (
	"fmt"
	"os"
	"strings"
)

type ReplaceModule struct {
}

//func (rm *ReplaceModule) Execute()

func (rm *ReplaceModule) replaceDomainWithEnv(url string) string {
	serverIP := os.Getenv("SERVER_IP")

	// #{http://server/ipinfo}을 "http://SERVER_IP/ipinfo"로 대체
	replacedURL := strings.Replace(url, "#{http://server/ipinfo}", fmt.Sprintf("http://%s/ipinfo", serverIP), -1)
	return replacedURL
}

func (rm *ReplaceModule) replaceagentUUID(str string, uuid string) string {
	replaceStr := strings.Replace(str, "#{agentUUID}", uuid, -1)
	return replaceStr
}

func (rm *ReplaceModule) replacePlaceholder(command string) string {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return command
	}
	return strings.ReplaceAll(command, `#{C:\Path\To\agent.exe}`, exePath)
}
