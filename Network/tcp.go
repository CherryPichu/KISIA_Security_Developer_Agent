package Network

import (
	"fmt"
	"github.com/HTTPs-omma/HTTPsBAS-HSProtocol/HSProtocol"
	"net"
)

func sendPacketByTcp(hs *HSProtocol.HS, conn net.Conn) (*HSProtocol.HS, error) {
	// HS 객체를 직렬화 (예: ToBytes 함수 사용)

	HSMgr := HSProtocol.NewHSProtocolManager()
	data, err := HSMgr.ToBytes(hs)
	if err != nil {
		fmt.Println("Error serializing HS object:", err)
		return nil, err
	}

	// 서버로 데이터 전송 (Payload 요청)
	_, err = conn.Write(data) // 에러
	if err != nil {
		fmt.Println("Error sending data to server:", err)
		return nil, err
	}
	//fmt.Println("Send data to server successfully")
	// 데이터 응답 (PayLoad 받아옴)
	msg := make([]byte, 1024*1024)
	conn.Read(msg)
	//msg = bytes.ReplaceAll(msg, []byte{0x00}, []byte{})

	// HS 객체를 직렬화 (예: ToBytes 함수 사용)
	HSMgr = HSProtocol.NewHSProtocolManager()

	ack, err := HSMgr.Parsing(msg)
	if err != nil {
		fmt.Println("Error Parsing ack data", err)
		return nil, err
	}
	//fmt.Println("Parsing ack data : ", ack.TotalLength)

	return ack, nil
}
