package Core

import (
	"agent/Model"
	"agent/Network"
	"fmt"
	"github.com/HTTPs-omma/HTTPsBAS-HSProtocol/HSProtocol"
)

func ChangeStatusToRun(hs *HSProtocol.HS) error {

	asdb, err := Model.NewAgentStatusDB()
	if err != nil {
		return err
	}
	err = asdb.UpdateStatus(HSProtocol.RUN)
	if err != nil {
		return err
	}

	// 새로운 구조체 생성 및 초기화
	newHS := HSProtocol.HS{
		ProtocolID:     hs.ProtocolID,
		HealthStatus:   HSProtocol.RUN,
		Command:        HSProtocol.UPDATE_AGENT_STATUS,
		Identification: hs.Identification, // 아직 구현 안함
		Checksum:       0,                 // 자동으로 채워줌
		TotalLength:    0,                 // 자동으로 채워줌
		UUID:           hs.UUID,           // 기존 hs의 UUID 사용
		Data:           []byte{},
	}

	// 서버에 새로운 구조체 전송
	ack, err := Network.NgMgr.SendPacket(newHS)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if ack.Command == HSProtocol.ERROR_ACK {
		return fmt.Errorf("Error acknowledgement in ChangStatus")
	}

	return nil
}

func ChangeStatusToWait(hs *HSProtocol.HS) error {
	asdb, err := Model.NewAgentStatusDB()
	if err != nil {
		return err
	}
	err = asdb.UpdateStatus(HSProtocol.WAIT)
	if err != nil {
		return err
	}

	newHS := HSProtocol.HS{
		ProtocolID:     hs.ProtocolID,
		HealthStatus:   HSProtocol.WAIT,
		Command:        HSProtocol.UPDATE_AGENT_STATUS,
		Identification: hs.Identification, // 아직 구현 안함
		Checksum:       0,                 // 자동으로 채워줌
		TotalLength:    0,                 // 자동으로 채워줌
		UUID:           hs.UUID,           // 기존 hs의 UUID 사용
		Data:           []byte{},
	}

	// 서버에 새로운 구조체 전송
	ack, err := Network.NgMgr.SendPacket(newHS)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if ack.Command == HSProtocol.ERROR_ACK {
		return fmt.Errorf("Error acknowledgement in ChangStatus")
	}

	return nil
}

func ChangeStatusToDeleted(hs *HSProtocol.HS) error {
	asdb, err := Model.NewAgentStatusDB()
	if err != nil {
		return err
	}
	err = asdb.UpdateStatus(HSProtocol.DELETED)
	if err != nil {
		return err
	}

	newHS := HSProtocol.HS{
		ProtocolID:     hs.ProtocolID,
		HealthStatus:   HSProtocol.DELETED,
		Command:        HSProtocol.UPDATE_AGENT_STATUS,
		Identification: hs.Identification, // 아직 구현 안함
		Checksum:       0,                 // 자동으로 채워줌
		TotalLength:    0,                 // 자동으로 채워줌
		UUID:           hs.UUID,           // 기존 hs의 UUID 사용
		Data:           []byte{},
	}

	// 서버에 새로운 구조체 전송
	ack, err := Network.NgMgr.SendPacket(newHS)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if ack.Command == HSProtocol.ERROR_ACK {
		return fmt.Errorf("Error acknowledgement in ChangStatus")
	}

	return nil
}
