package Network

import (
	"fmt"
	"github.com/HTTPs-omma/HTTPsBAS-HSProtocol/HSProtocol"
	"net"
)

func sendPacketByTcp(hs *HSProtocol.HS, conn net.Conn) (*HSProtocol.HS, error) {
	HSMgr := HSProtocol.NewHSProtocolManager()
	data, err := HSMgr.ToBytes(hs)
	if err != nil {
		fmt.Println("Error serializing HS object:", err)
		return nil, err
	}

	_, err = conn.Write(data) // 에러
	if err != nil {
		fmt.Println("Error sending data to server:", err)
		return nil, err
	}

	msg := make([]byte, 1024*1024)
	conn.Read(msg)

	HSMgr = HSProtocol.NewHSProtocolManager()

	ack, err := HSMgr.Parsing(msg)
	if err != nil {
		fmt.Println("Error Parsing ack data", err)
		return nil, err
	}

	return ack, nil
}
