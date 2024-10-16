package main

import (
	"agent/Core"
	"agent/Execute"
	"agent/Extension"
	"agent/Model"
	"agent/Network"
	"fmt"
	"github.com/HTTPs-omma/HTTPsBAS-HSProtocol/HSProtocol"
	"github.com/joho/godotenv"
	"os"
	"strings"
	"time"
)

const (
	ExecutePayLoad       string = "ExecutePayLoad"
	ExecuteCleanUp       string = "ExecuteCleanUp"
	GetSystemInfo        string = "GetSystemInfo"
	GetApplication       string = "GetApplication"
	StopAgent            string = "StopAgent"
	ChangeProtocolToTCP  string = "ChangeProtocolToTCP"
	ChangeProtocolToHTTP string = "ChangeProtocolToHTTP"
)

func main() {
	if err := loadEnv(); err != nil {
		fmt.Println("(5 초뒤 종료)에러 발생 : " + err.Error())
		time.Sleep(5 * time.Second)
		return
	}

	var err error
	Network.NgMgr, err = Network.NewNetworkManager()
	if err != nil {
		panic(err)
	}

	_, uuid, err := initSysutil()
	if err != nil {
		fmt.Println("(5 초뒤 종료)에러 발생 : " + err.Error())
		time.Sleep(5 * time.Second)
		return
	}

	if err := registerAgent(uuid); err != nil {
		fmt.Println("(5 초뒤 종료)에러 발생 : " + err.Error())
		time.Sleep(5 * time.Second)
		return
	}
	if err := collectInitialInfo(); err != nil {
		fmt.Println("(5 초뒤 종료)에러 발생 : " + err.Error())
		time.Sleep(5 * time.Second)
		return
	}

	// stage 2-3 : 반복 실행
	for {
		time.Sleep(3 * time.Second)
		if err := executeCommand(uuid); err != nil {
			fmt.Println("Error during command execution: ", err)
			continue
		}
	}
}

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("(5 초뒤 종료)에러 발생 : %v", err)
	}
	return nil
}

func initSysutil() (*Extension.Sysutils, [16]byte, error) {
	sysutil, err := Extension.NewSysutils()
	if err != nil {
		return nil, [16]byte{}, fmt.Errorf("시스템 유틸리티 초기화 실패: %v", err)
	}

	uuid, err := HSProtocol.HexStringToByteArray(sysutil.GetUniqueID())
	if err != nil {
		return nil, [16]byte{}, fmt.Errorf("UUID 변환 실패: %v", err)
	}

	return sysutil, uuid, nil
}

func registerAgent(uuid [16]byte) error {
	hsItem := &HSProtocol.HS{
		ProtocolID:     HSProtocol.UNKNOWN,
		HealthStatus:   HSProtocol.NEW,
		Command:        HSProtocol.UPDATE_AGENT_STATUS,
		Identification: 12345,
		Checksum:       0,
		TotalLength:    0,
		UUID:           uuid,
		Data:           []byte{},
	}

	ack, err := Network.NgMgr.SendPacket(hsItem)
	if err != nil {
		return fmt.Errorf("에이전트 등록 실패: %v", err)
	}

	if ack.Command == HSProtocol.ERROR_ACK {
		return fmt.Errorf("에이전트 등록 에러")
	}

	return nil
}

func collectInitialInfo() error {
	if err := Network.NgMgr.SendApplicationInfo(); err != nil {
		return fmt.Errorf("응용 프로그램 정보 전송 실패: %v", err)
	}

	if err := Network.NgMgr.SendSystemInfo(); err != nil {
		return fmt.Errorf("시스템 정보 전송 실패: %v", err)
	}

	return nil
}

func executeCommand(uuid [16]byte) error {
	agsdb, err := Model.NewAgentStatusDB()
	if err != nil {
		return fmt.Errorf("에이전트 상태 DB 오류: %v", err)
	}

	agsRcrd, err := agsdb.SelectAllRecords()
	if err != nil {
		return fmt.Errorf("상태 기록 조회 실패: %v", err)
	}

	protocol := agsRcrd[0].Protocol
	hsItem := &HSProtocol.HS{
		ProtocolID:     protocol,
		HealthStatus:   HSProtocol.RUN,
		Command:        HSProtocol.FETCH_INSTRUCTION,
		Identification: 12345,
		Checksum:       6789,
		TotalLength:    50,
		UUID:           uuid,
		Data:           []byte{},
	}

	fmt.Print("fetch instruction : ")
	ack, err := Network.NgMgr.SendPacket(hsItem)
	if err != nil {
		return fmt.Errorf("패킷 전송 실패: %v", err)
	}

	inst := &Core.ExtendedInstructionData{}
	instD, err := inst.GetInstData(ack.Data)
	if err != nil {
		return fmt.Errorf("명령어 데이터 처리 실패: %v", err)
	}

	if len(ack.Data) < 1 {
		fmt.Println("... NoData Wait")
		return nil
	}
	fmt.Println("... success")

	if err = Core.ChangeStatusToRun(ack); err != nil {
		return err
	}
	instD.Command = replacePlaceholder(instD.Command)
	instD.Cleanup = replacePlaceholder(instD.Cleanup)
	instD.Command = ReplaceDomainWithEnv(instD.Command)
	instD.Cleanup = ReplaceDomainWithEnv(instD.Cleanup)
	instD.Command = ReplaceagentUUID(instD.Command, HSProtocol.ByteArrayToHexString(uuid))
	instD.Cleanup = ReplaceagentUUID(instD.Cleanup, HSProtocol.ByteArrayToHexString(uuid))

	switch instD.AgentAction {
	case ExecutePayLoad:
		if err := runCommand(instD, hsItem); err != nil {
			return fmt.Errorf("명령어 실행 실패: %v", err)
		}
	case ExecuteCleanUp:
		if err := runCleanup(instD, hsItem); err != nil {
			return fmt.Errorf("명령어 실행 실패: %v", err)
		}
	case GetApplication:
		if err := Network.NgMgr.SendApplicationInfo(); err != nil {
			return fmt.Errorf("명령어 실행 실패: %v", err)
		}
	case GetSystemInfo:
		if err := Network.NgMgr.SendSystemInfo(); err != nil {
			return fmt.Errorf("명령어 실행 실패: %v", err)
		}
	case StopAgent:
		if err := Core.ChangeStatusToDeleted(ack); err != nil {
			return fmt.Errorf("명령어 실행 실패: %v", err)
		}
		fmt.Println("잠시후 종료...")
		time.Sleep(5 * time.Second)
		os.Exit(0)
	case ChangeProtocolToTCP:
		if err := Network.NgMgr.ChangeProtocol(HSProtocol.TCP); err != nil {
			return fmt.Errorf("명령어 실행 실패: %v", err)
		}
		fmt.Println("Agent Change Protocol Type by TCP")
	case ChangeProtocolToHTTP:
		if err := Network.NgMgr.ChangeProtocol(HSProtocol.HTTP); err != nil {
			return fmt.Errorf("명령어 실행 실패: %v", err)
		}
		fmt.Println("Agent Change Protocol Type by HTTP")
	default:
		fmt.Println("invalid Action String")
	}

	if err = Core.ChangeStatusToWait(ack); err != nil {
		return err
	}

	return nil
}

func runCommand(instD *Core.ExtendedInstructionData, hsItem *HSProtocol.HS) error {
	if err := Core.ChangeStatusToRun(hsItem); err != nil {
		return fmt.Errorf("상태 변경 실패: %v", err)
	}

	var shell Execute.ICommandExecutor
	if instD.Tool == "cmd" {
		shell = Execute.NewCmd() // 해당 부분 코드를 powershell 도 실핼 수 있게 수정할 것
	} else if instD.Tool == "powershell" {
		shell = Execute.NewPowerShell()
	} else if instD.Tool == "bash" {
		shell = Execute.NewShell()
	} else {
		fmt.Println(" No Shell! Tool 필드 값이 정확힌 지 확인해주세요. ex. cmd, powershell, bash")
		return nil
	}

	fmt.Println("====== Running ===== ")
	cmdLog, err := shell.Execute(instD.Command)
	fmt.Println("====== Execute Log ===== ")
	fmt.Println("Log : " + cmdLog)
	const maxLogLength = 10000
	if len(cmdLog) > maxLogLength {
		cmdLog = cmdLog[:maxLogLength] + "...(출력 생략됨)"
	}

	if err != nil {
		fmt.Println("====== Run Fail ===== ")
		fmt.Println()
		if err := Network.NgMgr.SendLogData(hsItem, err.Error(), instD.Command, instD.ID, instD.MessageUUID, Network.EXIT_FAIL); err != nil {
			return fmt.Errorf("실행 로그 전송 실패: %v", err)
		}
		return fmt.Errorf("명령어 실행 중 에러: %v", err)
	}
	if err := Network.NgMgr.SendLogData(hsItem, cmdLog, instD.Command, instD.ID, instD.MessageUUID, Network.EXIT_SUCCESS); err != nil {
		return fmt.Errorf("실행 로그 전송 실패: %v", err)
	}
	fmt.Println("====== Run success ===== ")
	fmt.Println()
	return nil
}

func runCleanup(instD *Core.ExtendedInstructionData, hsItem *HSProtocol.HS) error {
	if err := Core.ChangeStatusToRun(hsItem); err != nil {
		return fmt.Errorf("상태 변경 실패: %v", err)
	}

	var shell Execute.ICommandExecutor
	if instD.Tool == "cmd" {
		shell = Execute.NewCmd() // 해당 부분 코드를 powershell 도 실핼 수 있게 수정할 것
	} else if instD.Tool == "powershell" {
		shell = Execute.NewPowerShell()
	} else if instD.Tool == "bash" {
		shell = Execute.NewShell()
	} else {
		fmt.Println(" No Shell! Tool 필드 값에 정확한 값을 넣주세요.. ex. cmd, powershell, bash")
		return nil
	}

	fmt.Println("====== Running ===== ")
	cmdLog, err := shell.Execute(instD.Cleanup)
	fmt.Println("====== Execute Log ===== ")
	fmt.Println("Log : " + cmdLog)
	if err != nil {
		fmt.Println("====== Run Fail ===== ")
		if err := Network.NgMgr.SendLogData(hsItem, err.Error(), instD.Command, instD.ID, instD.MessageUUID, Network.EXIT_FAIL); err != nil {
			return fmt.Errorf("실행 로그 전송 실패: %v", err)
		}
		return fmt.Errorf("명령어 실행 중 에러: %v", err)
	}
	if err := Network.NgMgr.SendLogData(hsItem, cmdLog, instD.Command, instD.ID, instD.MessageUUID, Network.EXIT_SUCCESS); err != nil {
		return fmt.Errorf("실행 로그 전송 실패: %v", err)
	}
	fmt.Println("====== Run success ===== ")
	fmt.Println()
	return nil
}

func ReplaceDomainWithEnv(url string) string {
	serverIP := os.Getenv("SERVER_IP")

	// #{http://server/ipinfo}을 "http://SERVER_IP/ipinfo"로 대체
	replacedURL := strings.Replace(url, "#{http://server/ipinfo}", fmt.Sprintf("http://%s/ipinfo", serverIP), -1)
	return replacedURL
}

func ReplaceagentUUID(str string, uuid string) string {
	replaceStr := strings.Replace(str, "#{agentUUID}", uuid, -1)
	return replaceStr
}

func replacePlaceholder(command string) string {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return command
	}
	return strings.ReplaceAll(command, `#{C:\Path\To\agent.exe}`, exePath)
}
