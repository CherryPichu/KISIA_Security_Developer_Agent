package Execute

import (
	"agent/Core"
	"fmt"
	"github.com/HTTPs-omma/HTTPsBAS-HSProtocol/HSProtocol"
)

type ICommandExecutor interface {
	Execute(command string) (string, error)
}

func (shell *ICommandExecutor) RunCommand(instD *Core.ExtendedInstructionData, hsItem HSProtocol.HS) error {
	if err := Core.ChangeStatusToRun(&hsItem); err != nil {
		return fmt.Errorf("상태 변경 실패: %v", err)
	}

	if instD.Tool == "cmd" {
		shell = Execute.NewCmd() // 해당 부분 코드를 powershell 도 실핼 수 있게 수정할 것
	} else if instD.Tool == "powershell" {
		shell = Execute.NewPowerShell()
	} else if instD.Tool == "shell" {
		shell = Execute.NewShell()
	}

	fmt.Println("====== Running ===== ")
	cmdLog, err := shell.Execute(instD.Command)
	fmt.Println("====== Execute Log ===== ")
	fmt.Println("Log : " + cmdLog)
	if err != nil {
		fmt.Println("====== Run Fail ===== ")
		fmt.Println()
		if err := NgMgr.SendLogData(&hsItem, err.Error(), instD.Command, instD.ID, instD.MessageUUID, NgMgr.EXIT_FAIL); err != nil {
			return fmt.Errorf("실행 로그 전송 실패: %v", err)
		}
		return fmt.Errorf("명령어 실행 중 에러: %v", err)
	}
	if err := NgMgr.SendLogData(&hsItem, cmdLog, instD.Command, instD.ID, instD.MessageUUID, NgMgr.EXIT_SUCCESS); err != nil {
		return fmt.Errorf("실행 로그 전송 실패: %v", err)
	}
	fmt.Println("====== Run success ===== ")
	fmt.Println()
	return nil
}
