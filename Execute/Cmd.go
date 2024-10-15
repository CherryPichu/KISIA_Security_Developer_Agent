package Execute

import (
	"fmt"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Cmd struct {
	IsAvailable bool
	Name        string
}

func NewCmd() *Cmd {
	return &Cmd{
		Name:        "Cmd",
		IsAvailable: true,
	}
}

func (c *Cmd) Execute(command string) (string, error) {
	// 임시 배치 파일 생성
	batchFileName := "executeCommand.bat"
	batchFile, err := os.Create(batchFileName)
	if err != nil {
		fmt.Println("Error creating batch file:", err)
		return "", err
	}
	defer os.Remove(batchFileName) // 실행 후 배치 파일 삭제

	// 배치 파일에 명령어 작성
	_, err = batchFile.WriteString(command)
	if err != nil {
		fmt.Println("Error writing to batch file:", err)
		return "", err
	}

	// 배치 파일 닫기
	err = batchFile.Close()
	if err != nil {
		fmt.Println("Error closing batch file:", err)
		return "", err
	}

	// 배치 파일 실행
	cmd := exec.Command("cmd", "/C", batchFileName)
	output, err := cmd.CombinedOutput()
	decodedOutput, _ := decodeCP949(output)
	if err != nil {
		return decodedOutput, err
	}

	return decodedOutput, nil
}
func decodeCP949(input []byte) (string, error) {
	reader := transform.NewReader(strings.NewReader(string(input)), korean.EUCKR.NewDecoder())
	decoded, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
