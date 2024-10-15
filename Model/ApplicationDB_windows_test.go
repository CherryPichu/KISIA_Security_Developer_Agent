package Model

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApplicationDB_createTable(t1 *testing.T) {
	type fields struct {
		dbName string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{name: "Test case 1"},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t *testing.T) {
			appdb := NewApplicationDB()
			err := appdb.CreateTable()
			if err != nil {
				t1.Fatalf("Error creating table: %v", err)
			}

			// ========== 검증 =============
			dbPtr, err := getDBPtr()
			if err != nil {
				t1.Fatalf(appdb.dbName + " : DB 포인터를 가져올 수 없습니다. getDBPtr() 함수 오류\n" + err.Error())
			}

			query := fmt.Sprintf("select * from sqlite_master where name = '%s'", appdb.dbName)

			dsys := &sqlite_master{}

			rst := dbPtr.QueryRow(query).Scan(&dsys.Type, &dsys.name, &dsys.tbl_name, &dsys.rootpage, &dsys.sql)
			if rst != nil {
				t1.Fatalf(appdb.dbName + " : 생성된 테이블이 존재하지 않습니다.")
			}
			assert.Equal(t1, dsys.name, appdb.dbName)
		})
	}
}

func TestApplicationDB_CRUD(t1 *testing.T) {
	type fields struct {
		dbName string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{name: "Test case 1"},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t *testing.T) {
			// ============ Create 테스트 ================
			appdb := NewApplicationDB()
			appdb.CreateTable()

			wind32 := GetApplicationList()[0]
			data := DapplicationDB{}
			data.Name = wind32.Name
			data.Description = wind32.Description
			data.Version = wind32.Version
			data.Vendor = wind32.Vendor
			data.InstallDate2 = wind32.InstallDate2
			data.InstallLocation = wind32.InstallLocation
			data.InstallSource = wind32.InstallSource
			data.Language = wind32.Language
			data.PackageCode = wind32.PackageCode
			data.PackageName = wind32.PackageName
			data.RegCompany = wind32.RegCompany
			data.RegOwner = wind32.RegOwner
			data.URLInfoAbout = wind32.URLInfoAbout

			err := appdb.InsertRecord(data)
			if err != nil {
				t.Fatalf(appdb.dbName + " : insert 에러\n" + err.Error())
			}

			// ========== Create 검증 =============
			dbPtr, err := getDBPtr()
			if err != nil {
				t.Fatalf(appdb.dbName + " : DB 포인터를 가져올 수 없습니다. getDBPtr() 함수 오류\n" + err.Error())
			}

			query := fmt.Sprintf("select * from %s", appdb.dbName)
			data2 := &DapplicationDB{}
			row := dbPtr.QueryRow(query)
			if err != nil {
				t.Fatalf(appdb.dbName + " : select Qeury 오류\n" + err.Error())
			}

			err = row.Scan(&data2.ID, &data2.Name, &data2.Version, &data2.Language, &data2.Vendor,
				&data2.InstallDate2, &data2.InstallLocation, &data2.InstallSource, &data2.PackageName,
				&data2.PackageCode, &data2.RegCompany, &data2.RegOwner, &data2.URLInfoAbout, &data2.Description,
				&data2.isDeleted, &data2.CreateAt, &data2.UpdateAt, &data2.deletedAt)

			if err != nil {
				t.Fatalf(appdb.dbName + " : DapplicationDB 내용이 실제 DB 컬럼 내용과 일치하지 않습니다.\n" + err.Error())
			}

			assert.Equal(t, data2.URLInfoAbout, data.URLInfoAbout)
			assert.Equal(t, data2.Name, data.Name)
			assert.Equal(t, data2.RegCompany, data.RegCompany)
			assert.Equal(t, data2.RegOwner, data.RegOwner)
			assert.Equal(t, data2.Language, data.Language)
			assert.Equal(t, data2.InstallSource, data.InstallSource)
			assert.Equal(t, data2.InstallDate2, data.InstallDate2)
			assert.Equal(t, data2.Description, data.Description)
			assert.Equal(t, data2.Version, data.Version)

			// ============ Read 테스트 ================
			data3, err := appdb.SelectAllRecords()
			data4 := data3[0]
			if err != nil {
				t1.Fatalf("select 테스트 에러\n" + err.Error())
			}

			assert.Equal(t, data2.URLInfoAbout, data4.URLInfoAbout)
			assert.Equal(t, data2.Name, data4.Name)
			assert.Equal(t, data2.RegCompany, data4.RegCompany)
			assert.Equal(t, data2.RegOwner, data4.RegOwner)
			assert.Equal(t, data2.Language, data4.Language)
			assert.Equal(t, data2.InstallSource, data4.InstallSource)
			assert.Equal(t, data2.InstallDate2, data4.InstallDate2)
			assert.Equal(t, data2.Description, data4.Description)
			assert.Equal(t, data2.Version, data4.Version)

			// ============ Delete 테스트 ================
			data6, err := appdb.SelectByPackageCode(data.PackageCode)
			if err != nil {
				t1.Fatalf(": selectByPackageCode 실패\n" + err.Error())
			}
			assert.NotEqual(t, data6.PackageCode, "-1")

			err = appdb.DeleteByPackageCode(data.PackageCode)
			if err != nil {
				t1.Fatalf("삭제 실패\n" + err.Error())
			}

			data7, err := appdb.SelectByPackageCode(data.PackageCode)
			if err != nil {
				t1.Fatalf(": selectByPackageCode 실패\n" + err.Error())
			}
			assert.Equal(t, data7.PackageCode, "-1")
		})
	}
}

//func TestApplicationDB_CreateAll(t1 *testing.T) {
//	type fields struct {
//		dbName string
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		{name: "Test case 1"},
//	}
//	for _, tt := range tests {
//		t1.Run(tt.name, func(t *testing.T) {
//			// ============ Create 테스트 ================
//			appdb := NewApplicationDB()
//			appdb.createTable()
//
//			wind32 := getApplicationList()
//			for _, wind := range wind32 {
//				data := DapplicationDB{}
//				data.Name = wind.Name
//				data.Description = wind.Description
//				data.Version = wind.Version
//				data.Vendor = wind.Vendor
//				data.InstallDate2 = wind.InstallDate2
//				data.InstallLocation = wind.InstallLocation
//				data.InstallSource = wind.InstallSource
//				data.Language = wind.Language
//				data.PackageCode = wind.PackageCode
//				data.PackageName = wind.PackageName
//				data.RegCompany = wind.RegCompany
//				data.RegOwner = wind.RegOwner
//				data.URLInfoAbout = wind.URLInfoAbout
//
//				err := appdb.insertRecord(data)
//				if err != nil {
//					t.Fatalf(appdb.dbName + " : insert 에러\n" + err.Error())
//				}
//			}
//		})
//	}
//}
