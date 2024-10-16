package Model

import (
	"agent/Extension"
	"database/sql"
	"fmt"
	"github.com/HTTPs-omma/HTTPsBAS-HSProtocol/HSProtocol"
	"time"
)

// Protocol 유형을 정의합니다.

type AgentStatusDB struct {
	dbName string
}

type AgentStatusRecord struct {
	ID        int
	UUID      string
	Status    HSProtocol.AGENTSTATUS
	Protocol  HSProtocol.PROTOCOL
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewAgentStatusDB creates a new instance of AgentStatusDB with the default table name.
func NewAgentStatusDB() (*AgentStatusDB, error) {
	db := &AgentStatusDB{dbName: "AgentStatus"}
	err := db.CreateTable()
	if err != nil {
		return nil, err
	}
	isExist, err := db.ExistRecord()
	if isExist == false {
		sysutil, err := Extension.NewSysutils()
		if err != nil {
			return nil, err
		}
		err = db.InsertRecord(&AgentStatusRecord{
			ID:        0,
			UUID:      sysutil.GetUniqueID(),
			Status:    HSProtocol.NEW,
			Protocol:  HSProtocol.HTTP,
			CreatedAt: time.Now(),
		})
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

// CreateTable creates the AgentStatus table if it does not exist.
func (s *AgentStatusDB) CreateTable() error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStmt := `
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			uuid TEXT NOT NULL UNIQUE,
			status int,
			protocol int Default 0,
			createAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			updateAt DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`

	sqlStmt = fmt.Sprintf(sqlStmt, s.dbName)

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	sqlTrigger := fmt.Sprintf(`
		CREATE TRIGGER IF NOT EXISTS update_ModificationTime
		AFTER UPDATE ON %s
		FOR EACH ROW
		BEGIN	
			UPDATE %s SET
				updateAt = CURRENT_TIMESTAMP
			WHERE id = NEW.id;
		END;
	`, s.dbName, s.dbName)

	_, err = db.Exec(sqlTrigger)
	if err != nil {
		return err
	}

	return nil
}

func (s *AgentStatusDB) InsertRecord(data *AgentStatusRecord) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	exists, err := s.ExistRecord()
	if err != nil {
		return err
	}

	if exists {
		return s.InsertRecord(data)
	}

	query := fmt.Sprintf(`INSERT INTO %s (uuid, status, protocol) VALUES (?, ?, ?)`, s.dbName)
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(data.UUID, data.Status, data.Protocol)
	if err != nil {
		return err
	}

	return nil
}

func (s *AgentStatusDB) UpdateRecord(data *AgentStatusRecord) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(`UPDATE %s SET status = ?, protocol = ? WHERE uuid = ?`, s.dbName)
	_, err = db.Exec(query, data.Status, data.Protocol, data.UUID)
	if err != nil {
		return err
	}

	return nil
}

func (s *AgentStatusDB) UpdateStatus(status HSProtocol.AGENTSTATUS) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	// Update the status of the first (and only) record in the table
	query := fmt.Sprintf(`UPDATE %s SET status = ? WHERE id = (SELECT id FROM %s LIMIT 1)`, s.dbName, s.dbName)
	_, err = db.Exec(query, status)
	if err != nil {
		return err
	}

	return nil
}

func (s *AgentStatusDB) DeleteRecord(uuid string) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(`DELETE FROM %s WHERE uuid = ?`, s.dbName)
	_, err = db.Exec(query, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (s *AgentStatusDB) DeleteAllRecord() error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(`DELETE FROM %s`, s.dbName)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (s *AgentStatusDB) ExistRecord() (bool, error) {
	db, err := getDBPtr()
	if err != nil {
		return false, err
	}
	defer db.Close()

	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s)`, s.dbName)
	var exists bool
	err = db.QueryRow(query).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *AgentStatusDB) getAgentStatus() (int, error) {
	db, err := getDBPtr()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var status int
	query := fmt.Sprintf(`SELECT status FROM %s LIMIT 1`, s.dbName)
	err = db.QueryRow(query).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // Return 0 if no record is found
		}
		return 0, err
	}

	return status, nil
}

func (s *AgentStatusDB) getAgentUUID() (string, error) {
	db, err := getDBPtr()
	if err != nil {
		return "", err
	}
	defer db.Close()

	var uuid string
	query := fmt.Sprintf(`SELECT uuid FROM %s LIMIT 1`, s.dbName)
	err = db.QueryRow(query).Scan(&uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // Return empty string if no record is found
		}
		return "", err
	}

	return uuid, nil
}

func (s *AgentStatusDB) SelectAllRecords() ([]AgentStatusRecord, error) {
	db, err := getDBPtr()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf(`SELECT id, uuid, status, protocol, createAt, updateAt FROM %s`, s.dbName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := []AgentStatusRecord{}
	for rows.Next() {
		var record AgentStatusRecord
		err := rows.Scan(&record.ID, &record.UUID, &record.Status, &record.Protocol, &record.CreatedAt, &record.UpdatedAt)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

func (s *AgentStatusDB) SelectRecordByUUID(uuid string) ([]AgentStatusRecord, error) {
	db, err := getDBPtr()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf(`SELECT id, uuid, status, protocol, createAt, updateAt FROM %s WHERE uuid = ?`, s.dbName)
	row := db.QueryRow(query, uuid)

	var records []AgentStatusRecord
	var record AgentStatusRecord
	err = row.Scan(&record.ID, &record.UUID, &record.Status, &record.Protocol, &record.CreatedAt, &record.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil if no record is found
		}
		return nil, err
	}

	return append(records, record), nil
}
