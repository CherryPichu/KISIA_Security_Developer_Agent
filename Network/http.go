package Network

import (
	"bytes"
	"fmt"
	"github.com/HTTPs-omma/HTTPsBAS-HSProtocol/HSProtocol"
	"io"
	"net/http"
	"os"
)

//https://agent

func sendPacketByHttp(hs *HSProtocol.HS) (*HSProtocol.HS, error) {
	// HS 객체를 직렬화 (예: ToBytes 함수 사용)
	HSMgr := HSProtocol.NewHSProtocolManager()
	data, err := HSMgr.ToBytes(hs)
	if err != nil {
		fmt.Println("Error serializing HS  object:", err)
		return nil, err
	}

	// HTTP POST 요청 생성
	url := "http://" + os.Getenv("SERVER_IP") + ":" + os.Getenv("HTTP_PORT") + "/api/checkInstReq" // 실제 서버 URL로 변경
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error : ", err)
	}

	ack, err := HSMgr.Parsing(body)
	if err != nil {
		fmt.Println("Error : ", err)
	}

	return ack, nil
}
