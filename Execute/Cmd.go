package Execute

import (
	"bytes"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os/exec"
	"strings"
)

//type Cmd struct {
//	IsAvailable bool
//	Name        string
//}

//	func NewCmd() *Cmd {
//		return &Cmd{
//			Name:        "Cmd",
//			IsAvailable: true,
//		}
//	}

type Cmd struct {
	shortName string
	path      string
	execArgs  []string
}

func NewCmd() *Cmd {
	shell := &Cmd{
		shortName: "cmd",
		path:      "cmd.exe",
		execArgs:  []string{"/C"},
	}
	return shell
}

func (c *Cmd) Execute(command string) (string, error) {
	cmd := exec.Command(c.path) // * 제거
	commandLineComponents := append([]string{c.path}, c.execArgs...)
	commandLineComponents = append(commandLineComponents, command)

	//cmd.SysProcAttr.CmdLine = strings.Join(commandLineComponents, " ")
	//cmd.SysProcAttr.CmdLine = commandLineComponents
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// 명령어 실행
	err := cmd.Start()
	if err != nil {
		return "", err
	}

	// 명령어 실행 완료 대기
	err = cmd.Wait()
	if err != nil {
		return "", err
	}

	// using Chat gpt
	//if exitError, ok /*:= err.(*exec.ExitError); ok {
	//	// 프로세스 종료 코드가 0이 아니면 에러로 간주
	//	if exitCode := exitError.ExitCode(); exitCode != 0 {
	//		// 표준 출력과 에러 출력 결합하여 반환
	//		combinedOutput := stdoutBuf.String() + stderrBuf.String()
	//		output, _ := decodeCP9492(combinedOutput)
	//		return output, fmt.Errorf("command exited with code %d: %s", exitCode, output)
	//	}
	//}

	// 표준 출력과 에러 출력 결합
	combinedOutput := stdoutBuf.String() + stderrBuf.String()
	output, _ := decodeCP9492(combinedOutput)
	if err != nil {
		// 명령어가 실패했을 때 결합된 출력과 에러 반환
		return output, err
	}
	// 명령어가 성공했을 때 결합된 출력 반환
	return output, nil
}

//func (c *Cmd) Execute(command string) (string, error) {
//	// 임시 배치 파일 생성
//	batchFileName := "executeCommand.bat"
//	batchFile, err := os.Create(batchFileName)
//	if err != nil {
//		fmt.Println("Error creating batch file:", err)
//		return "", err
//	}
//	defer os.Remove(batchFileName) // 실행 후 배치 파일 삭제
//
//	// 배치 파일에 명령어 작성
//	_, err = batchFile.WriteString(command)
//	if err != nil {
//		fmt.Println("Error writing to batch file:", err)
//		return "", err
//	}
//
//	// 배치 파일 닫기
//	err = batchFile.Close()
//	if err != nil {
//		fmt.Println("Error closing batch file:", err)
//		return "", err
//	}
//
//	// 배치 파일 실행
//	cmd := exec.Command("cmd", "/C", batchFileName)
//	output, err := cmd.CombinedOutput()
//	decodedOutput, _ := decodeCP949(output)
//
//	if err != nil {
//		return decodedOutput, err
//	}
//
//	return decodedOutput, nil
//}

func decodeCP949(input []byte) (string, error) {
	reader := transform.NewReader(strings.NewReader(string(input)), korean.EUCKR.NewDecoder())
	decoded, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func decodeCP9492(input string) (string, error) {
	reader := transform.NewReader(strings.NewReader(input), korean.EUCKR.NewDecoder())
	decoded, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
