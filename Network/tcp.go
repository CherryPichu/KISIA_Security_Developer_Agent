package Network

import (
	"bytes"
	"fmt"
	"github.com/HTTPs-omma/HTTPsBAS-HSProtocol/HSProtocol"
	"net"
)

func sendPacketByTcp(hs HSProtocol.HS, conn net.Conn) (*HSProtocol.HS, error) {
	fmt.Println("Connected to TCP server")

	// HS 객체를 직렬화 (예: ToBytes 함수 사용)
	HSMgr := HSProtocol.NewHSProtocolManager()
	data, err := HSMgr.ToBytes(&hs)
	if err != nil {
		fmt.Println("Error serializing HS object:", err)
		return nil, err
	}

	// 서버로 데이터 전송 (Payload 요청)
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Error sending data to server:", err)
		return nil, err
	}
	fmt.Println("Data sent to server")

	// 데이터 응답 (PayLoad 받아옴)
	msg := make([]byte, 65536)
	conn.Read(msg)
	msg = bytes.ReplaceAll(msg, []byte{0x00}, []byte{})

	// HS 객체를 직렬화 (예: ToBytes 함수 사용)
	HSMgr = HSProtocol.NewHSProtocolManager()
	data, err = HSMgr.ToBytes(&hs)

	ack, err := HSMgr.Parsing(data)
	if err != nil {
		fmt.Println("Error Parsing ack data", err)
		return nil, err
	}

	return ack, nil
}
