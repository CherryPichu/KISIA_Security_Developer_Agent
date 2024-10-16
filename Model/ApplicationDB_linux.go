package Model

import (
	"encoding/json"
	"fmt"
	"time"
)

type ApplicationDB struct {
	dbName string
}

func NewApplicationDB() (metaTable *ApplicationDB, err error) {
	appDB := &ApplicationDB{"Application"}
	err = appDB.CreateTable()
	if err != nil {
		return nil, err
	}
	return appDB, nil
}

type DapplicationDB struct {
	ID              int       `json:"id"`
	AgentUUID       string    `json:"agent_uuid"`
	Name            string    `json:"name"`
	Version         string    `json:"version"`
	Language        string    `json:"language"`
	Vendor          string    `json:"vendor"`
	InstallDate2    string    `json:"install_date"`
	InstallLocation string    `json:"install_location"`
	InstallSource   string    `json:"install_source"`
	PackageName     string    `json:"package_name"`
	PackageCode     string    `json:"package_code"`
	RegCompany      string    `json:"reg_company"`
	RegOwner        string    `json:"reg_owner"`
	URLInfoAbout    string    `json:"url_info_about"`
	Description     string    `json:"description"`
	IsDeleted       bool      `json:"is_deleted"`
	CreateAt        time.Time `json:"create_at"`
	UpdateAt        time.Time `json:"update_at"`
	DeletedAt       time.Time `json:"deleted_at"`
}

func (a *ApplicationDB) CreateTable() error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStmt := `
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER PRIMARY KEY AUTOINCREMENT,       -- 내부 ID, 자동 증가
			AgentUUID VARCHAR(255),
			Name VARCHAR(255),                          -- 제품 이름
			Version VARCHAR(50),                        -- 제품 버전
			Language VARCHAR(10),                       -- 제품의 언어
			Vendor VARCHAR(255),                        -- 제품 공급자
			InstallDate2 VARCHAR(20),                   -- 설치 날짜
			InstallLocation TEXT,                       -- 패키지 설치 위치
			InstallSource TEXT,                         -- 설치 소스 위치
			PackageName VARCHAR(255),                   -- 원래 패키지 이름
			PackageCode VARCHAR(255) UNIQUE NOT NULL,  	-- 패키지 식별자 UUID
			RegCompany VARCHAR(255),                    -- 제품을 사용하는 것으로 등록된 회사 이름
			RegOwner VARCHAR(255),                      -- 제품을 사용하는 것으로 등록된 사용자 이름
			URLInfoAbout TEXT,                          -- 제품에 대한 정보가 제공되는 URL
			Description TEXT,                           -- 제품 설명
		    isDeleted bool DEFAULT FALSE, 				-- apllication 제거 여부를 파악함
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

/*
refer : https://learn.microsoft.com/en-us/previous-versions/windows/desktop/legacy/aa394378(v=vs.85)
class Win32_Product : CIM_Product

	{
	  uint16   AssignmentType;
	  string   Caption;
	  string   Description;
	  string   IdentifyingNumber;
	  string   InstallDate;
	  datetime InstallDate2;
	  string   InstallLocation;
	  sint16   InstallState;
	  string   HelpLink;
	  string   HelpTelephone;
	  string   InstallSource;
	  string   Language;
	  string   LocalPackage;
	  string   Name;
	  string   PackageCache;
	  string   PackageCode;
	  string   PackageName;
	  string   ProductID;
	  string   RegOwner;
	  string   RegCompany;
	  string   SKUNumber;
	  string   Transforms;
	  string   URLInfoAbout;
	  string   URLUpdateInfo;
	  string   Vendor;
	  uint32   WordCount;
	  string   Version;
	};
*/
type Win32_Product struct {
	Name            string // 제품 이름
	Version         string // 제품 버전
	Language        string // 제품의 언어
	Vendor          string // 제품 공급자
	InstallDate2    string // 설치 날짜
	InstallLocation string // 패키지 설치 위치
	InstallSource   string // 설치 소스 위치
	PackageName     string // 원래 패키지 이름
	PackageCode     string // 패키지 식별자
	RegCompany      string // 제품을 사용하는 것으로 등록된 회사 이름
	RegOwner        string // 제품을 사용하는 것으로 등록된 사용자 이름
	URLInfoAbout    string // 제품에 대한 정보가 제공되는 URL
	Description     string // 제품 설명
}

// ID              int    // 내부 ID, 자동 증가
// Name            string // 제품 이름
// Version         string // 제품 버전
// Language        string // 제품의 언어
// Vendor          string // 제품 공급자
// InstallDate2    string // 설치 날짜
// InstallLocation string // 패키지 설치 위치
// InstallSource   string // 설치 소스 위치
// PackageName     string // 원래 패키지 이름
// PackageCode     string // 패키지 식별자
// RegCompany      string // 제품을 사용하는 것으로 등록된 회사 이름
// RegOwner        string // 제품을 사용하는 것으로 등록된 사용자 이름
// URLInfoAbout    string // 제품에 대한 정보가 제공되는 URL
// Description     string // 제품 설명
// isDeleted       bool
// CreateAt        time.Time // 레코드 생성 시간
// UpdateAt        time.Time // 레코드 업데이트 시간
// deletedAt       time.Time
func GetApplicationList() []DapplicationDB {
	applist := []DapplicationDB{}
	return applist
}

const (
	PROGRAMEPATH_x64 string = "C:\\Program Files"
	PROGRAMEPATH_x86 string = "C:\\Program Files (x86)"
	STARTPROGRAGE           = "C:\\ProgramData\\Microsoft\\Windows\\Start Menu\\Programs"
	STARTMENU               = "C:\\ProgramData\\Microsoft\\Windows\\Start Menu"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (a *ApplicationDB) InsertRecord(data DapplicationDB) error {
	// ProductID 가 있는지 확인 후 중복되는 것이 없으면 insert 하기

	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(`INSERT INTO %s ( AgentUUID, Name, Version, Language, Vendor, 
        InstallDate2, InstallLocation, InstallSource, PackageName, PackageCode, RegCompany, 
        RegOwner, URLInfoAbout, Description ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, a.dbName)

	stmt, err := db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(&data.AgentUUID, &data.Name, &data.Version, &data.Language, &data.Vendor,
		&data.InstallDate2, &data.InstallLocation, &data.InstallSource, &data.PackageName,
		&data.PackageCode, &data.RegCompany, &data.RegOwner, &data.URLInfoAbout, &data.Description)

	if err != nil {
		return err
	}

	return nil
}

/*
selectRecords()를 통해 반환된 DsystemInfoDB 객체의 값을 수정한 후,
수정된 객체를 updateRecord 함수의 매개변수로 전달합시오
*/
func (a *ApplicationDB) UpdateByPackageCode(data *DapplicationDB) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(`UPDATE %s SET Name = ?, AgentUUID = ?, Version = ?, Language = ?, Vendor = ?, InstallDate2 = ?, InstallLocation = ?, InstallSource = ?, PackageName = ?, RegCompany = ?, RegOwner = ?, URLInfoAbout = ?, Description = ? WHERE PackageCode = ?`, a.dbName)
	_, err = db.Exec(query, data.Name, data.AgentUUID, data.Version, data.Language, data.Vendor, data.InstallDate2, data.InstallLocation, data.InstallSource, data.PackageName, data.RegCompany, data.RegOwner, data.URLInfoAbout, data.Description, data.PackageCode)
	if err != nil {
		return err
	}

	return nil
}

func (s *ApplicationDB) SelectByPackageCode(packageCode string) (*DapplicationDB, error) {
	db, err := getDBPtr()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf(`SELECT * FROM %s WHERE PackageCode = '%s' LIMIT 1`, s.dbName, packageCode)
	row, err := db.Query(query)
	defer row.Close()
	if err != nil {
		return nil, err
	}

	var data DapplicationDB

	if row.Next() == false {
		return &DapplicationDB{PackageCode: "-1"}, nil
	}
	err = row.Scan(&data.ID, &data.AgentUUID, &data.Name, &data.Version, &data.Language, &data.Vendor,
		&data.InstallDate2, &data.InstallLocation, &data.InstallSource, &data.PackageName,
		&data.PackageCode, &data.RegCompany, &data.RegOwner, &data.URLInfoAbout, &data.Description,
		&data.IsDeleted, &data.CreateAt, &data.UpdateAt, &data.DeletedAt)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (s *ApplicationDB) DeleteByPackageCode(packageCode string) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(`DELETE FROM %s WHERE PackageCode = ?`, s.dbName)
	_, err = db.Exec(query, packageCode)
	if err != nil {
		return err
	}

	return nil
}

func (s *ApplicationDB) DeleteAllRecord() error {
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

func (s *ApplicationDB) SelectAllRecords() ([]DapplicationDB, error) {
	db, err := getDBPtr()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf(`SELECT * FROM %s `, s.dbName)
	row, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	var rows []DapplicationDB

	for row.Next() {
		var data DapplicationDB

		err = row.Scan(&data.ID, &data.AgentUUID, &data.Name, &data.Version, &data.Language, &data.Vendor,
			&data.InstallDate2, &data.InstallLocation, &data.InstallSource, &data.PackageName,
			&data.PackageCode, &data.RegCompany, &data.RegOwner, &data.URLInfoAbout, &data.Description,
			&data.IsDeleted, &data.CreateAt, &data.UpdateAt, &data.DeletedAt)
		if err != nil {
			return nil, err
		}
		rows = append(rows, data)
	}

	return rows, nil
}

// ToJSON: 구조체를 JSON 바이트로 마샬링
func (s *ApplicationDB) ToJSON(data []DapplicationDB) ([]byte, error) {
	return json.Marshal(data)
}

// FromJSON: JSON 바이트를 구조체로 언마샬링
func (s *ApplicationDB) FromJSON(data []byte) ([]DapplicationDB, error) {
	var result []DapplicationDB
	err := json.Unmarshal(data, &result)
	return result, err
}
