package Execute

import (
	"agent/Model"
	"fmt"
	"github.com/HTTPs-omma/HTTPsBAS-HSProtocol/HSProtocol"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os/exec"
	"strings"
)

type PowerShell struct {
	Name        string
	IsAvailable bool
}

func NewPowerShell() *PowerShell {
	return &PowerShell{
		Name:        "PowerShell",
		IsAvailable: true,
	}
}

func (p *PowerShell) Execute(command string) (string, error) {
	fmt.Println("run Poershell Code : " + command)
	cmd := exec.Command("powershell", "-Command", command)

	// 표준 출력 및 표준 에러를 함께 캡처
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("command execution failed: %s, error: %w", string(output), err)
	}
	output, err = decodeCP949_poweshl(output)
	if err != nil {
		return string(output), fmt.Errorf("command execution failed: %s, error: %w", string(output), err)
	}

	agdb, err := Model.NewAgentStatusDB()
	reds, err := agdb.SelectAllRecords()
	if reds[0].Protocol == HSProtocol.TCP {
		const maxLogLength = 5000
		if len(output) > maxLogLength {
			return string(output[:maxLogLength]) + "(... 중간 생략)", nil
		}
	}

	return string(output), nil
}
func decodeCP949_poweshl(input []byte) ([]byte, error) {
	reader := transform.NewReader(strings.NewReader(string(input)), korean.EUCKR.NewDecoder())
	decoded, err := ioutil.ReadAll(reader)
	if err != nil {
		return []byte{}, err
	}
	return decoded, nil
}

///*
//파워셀을 빠르게 사용해보자 : https://stackoverflow.com/questions/65331558/how-to-call-powershell-from-go-faster
//*/
//func (p *PowerShell) execute_osExec(commands []string) ([]string, error) {
//	cmd := exec.Command("powershell", "-nologo", "-noprofile")
//	stdin, err := cmd.StdinPipe()
//	if err != nil {
//		log.Fatal(err)
//		p.isAvailable = false
//		return nil, err
//	}
//	p.isAvailable = true
//
//	go func() {
//		defer stdin.Close()
//		for _, command := range commands {
//			fmt.Fprintf(stdin, "%s\n", command)
//		}
//	}()
//
//	// 이렇게 안하면 속도가 너무 느림
//	out, err := cmd.CombinedOutput() //비동기적으로 동작한 함수가 실행된 코루틴의 작업이 완료될 때 까지 대기한다.
//	if err != nil {
//		log.Fatal(err)
//	}
//	log.Println(string(out))
//
//	return nil, nil
//}
