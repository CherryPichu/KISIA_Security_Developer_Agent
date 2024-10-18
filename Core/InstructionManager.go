package Core

import (
	"gopkg.in/yaml.v3"
)

type ExtendedInstructionData struct {
	ID               string   `yaml:"id"`
	MITREID          string   `yaml:"MITRE_ID"`
	Description      string   `yaml:"Description"`
	Tool             string   `yaml:"tool"`
	RequisiteCommand string   `yaml:"requisite_command"`
	Command          string   `yaml:"command"`
	Cleanup          string   `yaml:"cleanup"`
	MessageUUID      string   `yaml:"messageUUID"`
	AgentAction      string   `yaml:"agentAction"`
	Files            []string `yaml:"files"`
}

// ToBytes는 InstructionData 구조체를 YAML 바이트 슬라이스로 변환하는 함수입니다.
func (cd *ExtendedInstructionData) ToBytes() ([]byte, error) {
	// YAML로 직렬화
	data, err := yaml.Marshal(cd)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (cd *ExtendedInstructionData) GetInstData(bdata []byte) (*ExtendedInstructionData, error) {
	// 새로운 InstructionData 인스턴스를 생성
	var instData ExtendedInstructionData
	//fmt.Println("bytes : ", bdata)

	// YAML 데이터를 InstructionData 구조체로 역직렬화
	err := yaml.Unmarshal(bdata, &instData)
	if err != nil {
		return nil, err
	}

	return &instData, nil
}
