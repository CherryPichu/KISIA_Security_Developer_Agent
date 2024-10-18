package Network

import (
	"agent/Extension"
	"agent/Model"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/HTTPs-omma/HTTPsBAS-HSProtocol/HSProtocol"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var NgMgr *NetworkManager
var once sync.Once

type NetworkManager struct {
	protocol HSProtocol.PROTOCOL
	conn     *net.TCPConn
}

func NewNetworkManager() (*NetworkManager, error) {
	once.Do(func() {
		NgMgr = &NetworkManager{protocol: HSProtocol.HTTP}
	})

	agdb, err := Model.NewAgentStatusDB()
	if err != nil {
		return nil, errors.New("Network Manager 생성 실패")
	}

	records, err := agdb.SelectAllRecords()
	if err != nil {
		return nil, err
	}
	rd := records[0]
	NgMgr.protocol = rd.Protocol

	if NgMgr.protocol == HSProtocol.TCP {
		if err = NgMgr.connectTCP(); err != nil {
			return nil, err
		}
	}

	return NgMgr, nil
}

func (ng *NetworkManager) connectTCP() error {
	fmt.Println("Connection 시도")
	var conn net.Conn
	conn, err := net.Dial("tcp", os.Getenv("SERVER_IP")+":8080")
	if err != nil {
		return err
	}
	fmt.Println(os.Getenv("SERVER_IP") + ":8080")
	tcpConn := conn.(*net.TCPConn)

	message := ""
	// ping
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("TCP 연결 실패:", err)
		return err
	}
	ng.conn = tcpConn
	ng.conn.SetKeepAlive(true)
	ng.conn.SetKeepAlivePeriod(60 * time.Second)

	return nil
}

func (ng *NetworkManager) ChangeProtocol(ptl HSProtocol.PROTOCOL) error {
	agdb, err := Model.NewAgentStatusDB()
	if err != nil {
		return err
	}
	statusRecords, err := agdb.SelectAllRecords()
	fmt.Println(statusRecords)
	record := statusRecords[0]
	if record.Protocol == ptl {
		fmt.Println("Protocol is already set.")
		return nil
	}
	record.Protocol = ptl
	ng.protocol = ptl

	if ptl == HSProtocol.TCP {

		if err = ng.connectTCP(); err != nil {
			return err
		}
	}

	if err = agdb.UpdateRecord(&record); err != nil {
		return err
	}
	return nil
}

func (ng *NetworkManager) SendPacket(hs *HSProtocol.HS) (*HSProtocol.HS, error) {
	buffer := []byte{}
	hs.ProtocolID = ng.protocol
	switch ng.protocol {
	case HSProtocol.TCP:
		reader := bufio.NewReader(ng.conn)
		reader.Discard(reader.Buffered()) // 남은 버퍼를 버림
		_, err := ng.conn.Write(buffer)
		if err != nil {
			fmt.Printf("Error writing to connection: %v\n", err)
			// 재연결
			err = ng.connectTCP()
			if err != nil {
				return nil, err
			}
		}
		//fmt.Println("SendPacketByTCP")

		return sendPacketByTcp(hs, ng.conn)
	case HSProtocol.HTTP:
		return sendPacketByHttp(hs)
	}

	return nil, fmt.Errorf("protocol not support")
}

func (ng *NetworkManager) SendAck(hs *HSProtocol.HS) error {
	hsItem := &HSProtocol.HS{
		ProtocolID:     hs.ProtocolID,   //
		HealthStatus:   HSProtocol.WAIT, //
		Command:        HSProtocol.UPDATE_AGENT_STATUS,
		Identification: 12345, // 아직 구현 안함
		Checksum:       0,     // 자동으로 채워줌
		TotalLength:    0,     // 자동으로 채워줌
		UUID:           hs.UUID,
		Data:           []byte{},
	}

	_, err := ng.SendPacket(hsItem)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (ng *NetworkManager) SendErrorAck(hs *HSProtocol.HS) error {
	hsItem := &HSProtocol.HS{
		ProtocolID:     hs.ProtocolID,   //
		HealthStatus:   HSProtocol.WAIT, //
		Command:        HSProtocol.UPDATE_AGENT_STATUS,
		Identification: 12345, // 아직 구현 안함
		Checksum:       0,     // 자동으로 채워줌
		TotalLength:    0,     // 자동으로 채워줌
		UUID:           hs.UUID,
		Data:           []byte{},
	}

	_, err := ng.SendPacket(hsItem)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

type OperationLogDocument struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	AgentUUID   string             `bson:"agentUUID"`
	ProcedureID string             `bson:"procedureID"`
	MessageUUID string             `bson:"messageUUID"`
	ConductAt   time.Time          `bson:"conductAt"`
	ExitCode    int                `bson:"exitCode"`
	Log         string             `bson:"log"`
	Command     string             `bson:"command"` // Command 필드로 변경
}

const (
	EXIT_SUCCESS = 1
	EXIT_Unknown = 0
	EXIT_FAIL    = -1
)

type InstructionData struct {
	ID               string `yaml:"id"`
	MITREID          string `yaml:"MITRE_ID"`
	Description      string `yaml:"Description"`
	Escalation       bool   `yaml:"Escalation"`
	Tool             string `yaml:"tool"`
	RequisiteCommand string `yaml:"requisite_command"`
	Command          string `yaml:"command"`
	Cleanup          string `yaml:"cleanup"`
}

func (ng *NetworkManager) SendLogData(hs *HSProtocol.HS, cmdLog string, command string, PID string, messageUUID string, resultCode int) error {
	logD, err := json.Marshal(&OperationLogDocument{
		ID:          primitive.ObjectID{},
		Command:     command,
		AgentUUID:   HSProtocol.ByteArrayToHexString(hs.UUID),
		ProcedureID: PID,
		MessageUUID: messageUUID,
		ConductAt:   time.Now(),
		Log:         cmdLog,
		ExitCode:    resultCode,
	})
	if err != nil {
		fmt.Println("Error : ", err)
	}

	hsItem := &HSProtocol.HS{
		ProtocolID:     hs.ProtocolID,
		HealthStatus:   0,
		Command:        HSProtocol.SEND_PROCEDURE_LOG,
		Identification: 12345,
		Checksum:       6789,
		TotalLength:    50,
		UUID:           hs.UUID,
		Data:           logD,
	}
	ack, err := ng.SendPacket(hsItem)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if ack.Command == HSProtocol.ERROR_ACK {
		return fmt.Errorf("Error Ack!")
	}

	return nil
}

// 데이터를 수집하고 보내는 과정
func (ng *NetworkManager) SendApplicationInfo() error {
	sysuil, err := Extension.NewSysutils()
	if err != nil {
		return err
	}
	strUuid := sysuil.GetUniqueID()

	fmt.Println("application 정보 가져오는 중...")
	//applist := Model.GetApplicationList() // 현재 사용 안함

	fileNames, err := Model.GetApplicationList_fileBase()
	if err != nil {
		return fmt.Errorf("GetApplication File Names error : %v", err)
	}

	prgdb, err := Model.NewProgramsDB()
	if err != nil {
		return fmt.Errorf("NewProgramsDB error: %v", err)
	}

	err = prgdb.DeleteAllRecords()
	if err != nil {
		return fmt.Errorf("prgdb.DeleteAllRecords error: %v", err)
	}

	for _, fileName := range fileNames {
		prgdb.InsertRecord(strUuid, fileName)
	}

	fmt.Printf("완료 len(data) : ")
	fmt.Println(len(fileNames))

	fileList, err := prgdb.SelectAllRecords()
	if err != nil {
		return err
	}

	uuid, err := HSProtocol.HexStringToByteArray(sysuil.GetUniqueID())
	if err != nil {
		return err
	}

	if len(fileNames) < 10 {
		var chunk []Model.ProgramsRecord
		chunk = fileList

		bapplist, err := prgdb.ToJSON(chunk)
		if err != nil {
			return err
		}

		hsItem := &HSProtocol.HS{
			ProtocolID:     HSProtocol.UNKNOWN, // 자동으로 채워줌
			HealthStatus:   HSProtocol.NEW,
			Command:        HSProtocol.SEND_AGENT_APP_INFO,
			Identification: 12345, // 아직 구현 안함
			Checksum:       0,     // 자동으로 채워줌
			TotalLength:    0,     // 자동으로 채워줌
			UUID:           uuid,
			Data:           bapplist,
		}
		ack, err := ng.SendPacket(hsItem)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if ack.Command == HSProtocol.ERROR_ACK {
			return fmt.Errorf("Application 정보 송신 실패")
		}

		return nil
	} else {
		for i := 0; i < 10; i++ {
			offset := len(fileNames) / 10
			var chunk []Model.ProgramsRecord
			chunk = fileList[offset*i : offset*(i+1)]

			bapplist, err := prgdb.ToJSON(chunk)
			if err != nil {
				return err
			}

			hsItem := &HSProtocol.HS{
				ProtocolID:     HSProtocol.UNKNOWN, // 자동으로 채워줌
				HealthStatus:   HSProtocol.NEW,
				Command:        HSProtocol.SEND_AGENT_APP_INFO,
				Identification: 12345, // 아직 구현 안함
				Checksum:       0,     // 자동으로 채워줌
				TotalLength:    0,     // 자동으로 채워줌
				UUID:           uuid,
				Data:           bapplist,
			}
			ack, err := ng.SendPacket(hsItem)
			if err != nil {
				fmt.Println(err)
				return err
			}
			if ack.Command == HSProtocol.ERROR_ACK {
				return fmt.Errorf("Application 정보 송신 실패")
			}
		}
	}
	return nil
}

// ByChatgpt
func (ng *NetworkManager) getIPv4AndMAC() (string, string, error) {
	// 시스템의 모든 네트워크 인터페이스 가져오기
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("네트워크 인터페이스를 가져오는 데 실패했습니다: %v", err)
	}

	// 이더넷 및 Wi-Fi 인터페이스만 출력
	for _, i := range interfaces {
		// 네트워크 인터페이스가 활성화되어 있고, 루프백이 아닌 경우 필터링
		if i.Flags&net.FlagUp != 0 && i.Flags&net.FlagLoopback == 0 && !strings.Contains(i.Name, "VMware") {
			fmt.Printf("인터페이스 이름: %s\n", i.Name)

			// MAC 주소 출력 (만약 존재한다면)
			if i.HardwareAddr != nil {
				fmt.Printf("MAC 주소: %s\n", i.HardwareAddr)
			} else {
				return "", "", fmt.Errorf("MAC 주소 없음")
			}

			// 네트워크 인터페이스에 연결된 주소 가져오기
			addrs, err := i.Addrs()
			if err != nil {
				return "", "", fmt.Errorf("인터페이스 %s의 주소를 가져오는 데 실패했습니다: %v", i.Name, err)

			}

			for _, addr := range addrs {
				fmt.Printf("주소: %v\n", addr)
			}
			fmt.Println("-----------------------------")

			return addrs[1].String(), i.HardwareAddr.String(), nil
		}
	}
	return "", "", fmt.Errorf("활성화된 인터페이스 카드가 없음!")
}

func (ng *NetworkManager) SendSystemInfo() error {
	sysdb := Model.NewSystemInfoDB()
	sysuil, err := Extension.NewSysutils()
	if err != nil {
		return err
	}

	var ip, mac string
	if ip, mac, err = ng.getIPv4AndMAC(); err != nil {
		ip = ""
		mac = ""
	}

	sysInfo := &Model.DsystemInfoDB{
		0,
		sysuil.GetUniqueID(),
		sysuil.GetHostName(),
		sysuil.GetOsName(),
		sysuil.GetOsVersion(),
		sysuil.GetFamily(),
		sysuil.GetArchitecture(),
		sysuil.GetKernelVersion(),
		sysuil.GetBootTime(),
		ip,
		mac,
		time.Now(),
		time.Now(),
	}
	fmt.Printf("해당 PC UUID : ")
	fmt.Println(sysuil.GetUniqueID())

	err = sysdb.InsertRecord(sysInfo)
	if err != nil {
		return err
	}
	uuid, err := HSProtocol.HexStringToByteArray(sysuil.GetUniqueID())
	if err != nil {
		return fmt.Errorf("HexStringToByteArray error : %v", err)
	}

	bsysinfo, err := json.Marshal(sysInfo)
	hsItem := &HSProtocol.HS{
		ProtocolID:     HSProtocol.HTTP, //
		HealthStatus:   HSProtocol.NEW,  //
		Command:        HSProtocol.SEND_AGENT_SYS_INFO,
		Identification: 12345, // 아직 구현 안함
		Checksum:       0,     // 자동으로 채워줌
		TotalLength:    0,     // 자동으로 채워줌
		UUID:           uuid,
		Data:           bsysinfo,
	}

	ack, err := ng.SendPacket(hsItem)
	if err != nil {
		return err
	}
	if ack.Command == HSProtocol.ERROR_ACK {
		return fmt.Errorf("sysinfo 정보 송신 실패")
	}

	return nil
}
