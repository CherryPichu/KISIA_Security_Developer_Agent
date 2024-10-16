package Test

import (
	"agent/Network"
	"fmt"
	"github.com/HTTPs-omma/HTTPsBAS-HSProtocol/HSProtocol"
	"net"
	"os"
	"testing"
)

func TestNetworkManager_connectTCP(t *testing.T) {
	type fields struct {
		protocol HSProtocol.PROTOCOL
		conn     *net.TCPConn
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "Test 연결 테스트"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("SERVER_IP", "uskawjdu.iptime.org")

			ng, err := Network.NewNetworkManager()
			err = ng.ChangeProtocol(HSProtocol.TCP)
			if err != nil {
				fmt.Println(err)
				return
			}

			uuid, err := HSProtocol.HexStringToByteArray("09a4e53c7a1c4b4e9a519f36df29d8a2")
			if err != nil {
				fmt.Println(err)
				return
			}

			hsItem := &HSProtocol.HS{
				ProtocolID:     1,
				HealthStatus:   0,
				Command:        HSProtocol.FETCH_INSTRUCTION,
				Identification: 12345,
				Checksum:       6789,
				TotalLength:    50,
				//UUID:           [16]byte{0xc3, 0xcb, 0x84, 0x23, 0x34, 0x16, 0x49, 0x76, 0x94, 0x56, 0x9d, 0x75, 0x9a, 0x8a, 0x13, 0xe7},
				UUID: uuid,
				Data: []byte{},
			}

			ack, err := ng.SendPacket(hsItem)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(ack.Command)

		})
	}
}
