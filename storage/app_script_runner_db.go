package storage

import (
	"database/sql"
	"github.com/BabySid/gobase"
	"path/filepath"
	"sid-desktop/common"
	"sync"
)

type AppScriptRunnerDB struct {
	sqlite gobase.SQLite
	dbName string

	appScriptTbl string
}

var (
	appScriptRunnerDB   *AppScriptRunnerDB
	appScriptRunnerOnce sync.Once
)

func GetAppScriptRunnerDB() *AppScriptRunnerDB {
	appScriptRunnerOnce.Do(func() {
		appScriptRunnerDB = &AppScriptRunnerDB{
			dbName:       "sid_app_script_runner.db",
			appScriptTbl: "app_script",
		}
	})
	return appScriptRunnerDB
}

func (db *AppScriptRunnerDB) Open(rootPath string) error {
	return db.sqlite.Open(filepath.Join(rootPath, db.dbName))
}

func (db *AppScriptRunnerDB) Close() {
	_ = db.sqlite.Close()
}

func (db *AppScriptRunnerDB) Init() error {
	if err := db.createAppScript(); err != nil {
		return err
	}
	return nil
}

func (db *AppScriptRunnerDB) createAppScript() error {
	sqlStrArray := []string{
		"drop table if exists " + db.appScriptTbl + ";",
		"CREATE TABLE IF NOT EXISTS " + db.appScriptTbl + `(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(256) NOT NULL,
			cont TEXT NOT NULL,
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

func (db *AppScriptRunnerDB) NeedInit() (bool, error) {
	tables, err := db.sqlite.ListTables()
	if err != nil {
		return false, err
	}

	if gobase.ContainsString(tables, db.appScriptTbl) < 0 {
		return true, nil
	}

	return false, nil
}

func (db *AppScriptRunnerDB) LoadScriptFiles() (*common.ScriptFileList, error) {
	files := common.NewScriptFileList()

	err := db.sqlite.Query(func(rows *sql.Rows) error {
		for rows.Next() {
			var script common.ScriptFile
			err := rows.Scan(&script.ID, &script.Name, &script.Cont, &script.CreateTime, &script.AccessTime)
			if err != nil {
				return err
			}
			files.Append(script)
		}
		return nil
	}, "select id, name, cont, create_time, access_time from "+db.appScriptTbl+" order by create_time desc")

	return files, err
}

func (db *AppScriptRunnerDB) AddScriptFile(file common.ScriptFile) error {
	_, _, err := db.sqlite.Exec("insert into "+db.appScriptTbl+
		" (name, cont, create_time, access_time) values(?, ?, ?, ?)",
		file.Name, file.Cont, file.CreateTime, file.AccessTime)
	return err
}

func (db *AppScriptRunnerDB) AddScriptFileList(files *common.ScriptFileList) error {
	scripts := files.GetScriptFiles()

	tx, err := db.sqlite.Begin()
	if err != nil {
		return err
	}

	for _, file := range scripts {
		_, _, err := tx.Exec("insert into "+db.appScriptTbl+
			" (name, cont, create_time, access_time) values(?, ?, ?, ?)",
			file.Name, file.Cont, file.CreateTime, file.AccessTime)

		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	_ = tx.Commit()
	return nil
}

func (db *AppScriptRunnerDB) DelScriptFile(file common.ScriptFile) error {
	_, _, err := db.sqlite.Exec("delete from "+db.appScriptTbl+" where id = ?", file.ID)
	return err
}

func (db *AppScriptRunnerDB) UpdateScriptFile(file common.ScriptFile) error {
	_, _, err := db.sqlite.Exec("update "+db.appScriptTbl+" set name = ?, cont = ?, access_time = ? where id = ?",
		file.Name, file.Cont, file.AccessTime, file.ID)
	return err
}
