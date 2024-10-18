package Model

import (
	"agent/Execute"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type ProgramsDB struct {
	dbName string
}

func NewProgramsDB() (*ProgramsDB, error) {
	db := &ProgramsDB{dbName: "Programs"}
	err := db.CreateTable()
	if err != nil {
		return nil, err
	}
	return db, nil
}

type ProgramsRecord struct {
	ID        int       `json:"id"`
	AgentUUID string    `json:"agent_uuid"`
	FileName  string    `json:"file_name"`
	CreateAt  time.Time `json:"create_at"`
	UpdateAt  time.Time `json:"update_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (a *ProgramsDB) CreateTable() error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStmt := `
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER PRIMARY KEY AUTOINCREMENT,       -- 내부 ID, 자동 증가
			AgentUUID VARCHAR(255),
			FileName VARCHAR(255),
			createAt DATETIME DEFAULT CURRENT_TIMESTAMP, -- 레코드 생성 시간
			updateAt DATETIME DEFAULT CURRENT_TIMESTAMP,  -- 레코드 업데이트 시간
		    deletedAt DATETIME DEFAULT CURRENT_TIMESTAMP	-- 제거된 시간
		);
	`
	sqlStmt = fmt.Sprintf(sqlStmt, a.dbName)

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	sqlModifyTrigger := fmt.Sprintf(`
		CREATE TRIGGER IF NOT EXISTS update_ModificationTime
		AFTER UPDATE ON %s
		FOR EACH ROW
		BEGIN
			UPDATE %s SET
				updateAt = CURRENT_TIMESTAMP
			WHERE id = NEW.id;
		END;
	`, a.dbName, a.dbName)

	_, err = db.Exec(sqlModifyTrigger)
	if err != nil {
		return err
	}

	return nil
}

func GetApplicationList_fileBase() ([]string, error) {
	var shell Execute.ICommandExecutor
	shell = Execute.NewShell()

	output, err := shell.Execute("apt list --installed | grep -v \"automatic\" | awk '{print $1}' | head -n 300")
	if err != nil {
		return nil, err
	}

	var fileNameList []string
	for _, fileName := range strings.Split(output, "\n") {
		if fileName == "" {
			continue
		}
		fileNameList = append(fileNameList, fileName)
	}

	return fileNameList, nil
}

// InsertRecord: 레코드 삽입
func (a *ProgramsDB) InsertRecord(agentUUID string, fileName string) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(`INSERT INTO %s (AgentUUID, FileName) VALUES (?, ?)`, a.dbName)
	_, err = db.Exec(query, agentUUID, fileName)
	if err != nil {
		return err
	}

	return nil
}

// SelectAllRecords: 모든 레코드를 선택
func (a *ProgramsDB) SelectAllRecords() ([]ProgramsRecord, error) {
	db, err := getDBPtr()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf(`SELECT id, AgentUUID, FileName, createAt, updateAt, deletedAt FROM %s`, a.dbName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []ProgramsRecord
	for rows.Next() {
		var record ProgramsRecord
		err = rows.Scan(&record.ID, &record.AgentUUID, &record.FileName, &record.CreateAt, &record.UpdateAt, &record.DeletedAt)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

// SelectRecordsByUUID: 특정 AgentUUID에 따른 레코드를 선택
func (a *ProgramsDB) SelectRecordsByUUID(agentUUID string) ([]ProgramsRecord, error) {
	db, err := getDBPtr() // 데이터베이스 연결
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// 쿼리문 작성 (AgentUUID에 따라 필터링)
	query := fmt.Sprintf(`SELECT id, AgentUUID, FileName, createAt, updateAt, deletedAt FROM %s WHERE AgentUUID = ?`, a.dbName)
	rows, err := db.Query(query, agentUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []ProgramsRecord
	for rows.Next() {
		var record ProgramsRecord
		// 각 필드의 값을 스캔하여 구조체에 저장
		err = rows.Scan(&record.ID, &record.AgentUUID, &record.FileName, &record.CreateAt, &record.UpdateAt, &record.DeletedAt)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

// UpdateRecordByID: ID를 기준으로 레코드 업데이트
func (a *ProgramsDB) UpdateRecordByID(id int, newFileName string) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(`UPDATE %s SET FileName = ?, updateAt = ? WHERE id = ?`, a.dbName)
	_, err = db.Exec(query, newFileName, time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}

// DeleteRecordByID: ID를 기준으로 레코드 삭제
func (a *ProgramsDB) DeleteRecordByID(id int) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, a.dbName)
	_, err = db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (a *ProgramsDB) DeleteAllRecords() error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(`DELETE FROM %s`, a.dbName)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// ToJSON: ProgramNameDB 레코드 리스트를 JSON 바이트로 변환
func (s *ProgramsDB) ToJSON(data []ProgramsRecord) ([]byte, error) {
	// JSON 마샬링하여 []byte로 반환
	return json.Marshal(data)
}

// FromJSON: JSON 바이트를 ProgramNameDB 레코드 리스트로 변환
func (s *ProgramsDB) FromJSON(data []byte) ([]ProgramsRecord, error) {
	var result []ProgramsRecord
	// JSON 언마샬링하여 구조체로 변환
	err := json.Unmarshal(data, &result)
	return result, err
}
