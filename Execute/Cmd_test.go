package Execute_test

import (
	"agent/Execute"
	"testing"
)

func TestCmdExecute(t *testing.T) {
	// Cmd 객체 생성
	cmdExecutor := Execute.NewCmd()

	// 테스트할 명령어
	command := "echo Hello, World"

	// 명령어 실행
	output, err := cmdExecutor.Execute(command)
	if err != nil {
		t.Fatalf("Command execution failed with error: %v", err)
	}

	// 기대하는 출력값 (Windows의 echo 명령어 결과)
	expectedOutput := "Hello, World\r\n"

	// 출력이 기대하는 값과 같은지 확인
	if output != expectedOutput {
		t.Errorf("Expected output: %v, got: %v", expectedOutput, output)
	}
}

func TestCmdExecuteInvalidCommand(t *testing.T) {
	// Cmd 객체 생성
	cmdExecutor := Execute.NewCmd()

	// 잘못된 명령어 테스트
	command := "invalid_command"

	// 명령어 실행
	_, err := cmdExecutor.Execute(command)
	if err == nil {
		t.Fatalf("Expected error for invalid command but got none")
	}
}
