package storage

import (
	"database/sql"
	"encoding/json"
	"github.com/BabySid/gobase"
	"path/filepath"
	"sid-desktop/desktop/common"
	"sync"
)

type AppDevToolDB struct {
	sqlite gobase.SQLite
	dbName string

	httpClientHistoryTbl string
}

var (
	appDevToolDB   *AppDevToolDB
	appDevToolOnce sync.Once
)

func GetAppDevToolDB() *AppDevToolDB {
	appDevToolOnce.Do(func() {
		appDevToolDB = &AppDevToolDB{
			dbName:               "sid_app_dev_tool.db",
			httpClientHistoryTbl: "http_client_history",
		}
	})
	return appDevToolDB
}

func (db *AppDevToolDB) Open(rootPath string) error {
	return db.sqlite.Open(filepath.Join(rootPath, db.dbName))
}

func (db *AppDevToolDB) Close() {
	_ = db.sqlite.Close()
}

func (db *AppDevToolDB) Init() error {
	if err := db.createAppDevTool(); err != nil {
		return err
	}
	return nil
}

func (db *AppDevToolDB) createAppDevTool() error {
	sqlStrArray := []string{
		"drop table if exists " + db.httpClientHistoryTbl + ";",
		"CREATE TABLE IF NOT EXISTS " + db.httpClientHistoryTbl + `(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			method VARCHAR(16) NOT NULL,
			url VARCHAR(512) NOT NULL,
			headers VARCHAR(1024) NOT NULL,
			body_type VARCHAR(32) NOT NULL,
			body icon NULL,
			create_time INTEGER NOT NULL,
			access_time INTEGER NOT NULL
    	);
		`,
	}

	for _, sqlStr := range sqlStrArray {
		_, _, err := db.sqlite.Exec(sqlStr)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *AppDevToolDB) NeedInit() (bool, error) {
	tables, err := db.sqlite.ListTables()
	if err != nil {
		return false, err
	}

	if gobase.ContainsString(tables, db.httpClientHistoryTbl) < 0 {
		return true, nil
	}

	return false, nil
}

func (db *AppDevToolDB) LoadHttpClientHistory() (*common.HttpRequestList, error) {
	reqs := common.NewHttpRequestList()

	err := db.sqlite.Query(func(rows *sql.Rows) error {
		for rows.Next() {
			var req common.HttpRequest
			headStr := ""
			err := rows.Scan(&req.ID, &req.Method, &req.Url, &headStr, &req.ReqBodyType, &req.ReqBody, &req.CreateTime, &req.AccessTime)
			if err != nil {
				return err
			}
			err = json.Unmarshal([]byte(headStr), &req.ReqHeader)
			if err != nil {
				return err
			}
			reqs.Append(req)
		}
		return nil
	}, "select id, method, url, headers, body_type, body, create_time, access_time from "+db.httpClientHistoryTbl+" order by access_time desc")

	return reqs, err
}

func (db *AppDevToolDB) AddHttpRequest(req *common.HttpRequest) error {
	header, err := json.Marshal(req.ReqHeader)
	if err != nil {
		return err
	}

	id, _, err := db.sqlite.Exec("insert into "+db.httpClientHistoryTbl+
		" (method, url, headers, body_type, body, create_time, access_time) values(?, ?, ?, ?, ?, ?, ?)",
		req.Method, req.Url, string(header), req.ReqBodyType, req.ReqBody, req.CreateTime, req.AccessTime)
	if err != nil {
		return err
	}
	req.ID = id
	return nil
}

func (db *AppDevToolDB) UpdateHttpRequest(req *common.HttpRequest) error {
	header, err := json.Marshal(req.ReqHeader)
	if err != nil {
		return err
	}

	_, _, err = db.sqlite.Exec("update "+db.httpClientHistoryTbl+" set headers = ?, body_type = ?, body = ?, access_time = ? where id = ?",
		string(header), req.ReqBodyType, req.ReqBody, req.AccessTime, req.ID)
	return err
}
