package storage

import (
	"database/sql"
	"github.com/BabySid/gobase"
	"path/filepath"
	"sid-desktop/common/apps"
	"sync"
)

type AppLauncherDB struct {
	sqlite gobase.SQLite
	dbName string

	appIndexTbl string
}

var (
	appLauncherDB     *AppLauncherDB
	appLauncherDBOnce sync.Once
)

func GetAppLauncherDB() *AppLauncherDB {
	appLauncherDBOnce.Do(func() {
		appLauncherDB = &AppLauncherDB{
			dbName:      "sid_app_launcher.db",
			appIndexTbl: "app_index",
		}
	})
	return appLauncherDB
}

func (db *AppLauncherDB) Open(rootPath string) error {
	return db.sqlite.Open(filepath.Join(rootPath, db.dbName))
}

func (db *AppLauncherDB) Close() {
	_ = db.sqlite.Close()
}

func (db *AppLauncherDB) Init() error {
	if err := db.createAppIndex(); err != nil {
		return err
	}
	return nil
}

func (db *AppLauncherDB) createAppIndex() error {
	sqlStrArray := []string{
		"drop table if exists " + db.appIndexTbl + ";",
		"CREATE TABLE IF NOT EXISTS " + db.appIndexTbl + `(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			app_name VARCHAR(256) NOT NULL,
			app_full_path VARCHAR(256) NOT NULL,
			icon BLOB NULL,
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

func (db *AppLauncherDB) NeedInit() (bool, error) {
	tables, err := db.sqlite.ListTables()
	if err != nil {
		return false, err
	}

	if gobase.ContainsString(tables, db.appIndexTbl) < 0 {
		return true, nil
	}

	total, err := db.sqlite.GetTableRowCount(db.appIndexTbl)
	if err != nil {
		return false, err
	}

	if total == 0 {
		return true, nil
	}

	return false, nil
}

func (db *AppLauncherDB) LoadAppIndex() (*apps.AppList, error) {
	myApps := apps.NewAppList()

	err := db.sqlite.Query(func(rows *sql.Rows) error {
		for rows.Next() {
			var app apps.AppInfo
			err := rows.Scan(&app.AppID, &app.AppName, &app.FullPath, &app.Icon, &app.CreateTime, &app.AccessTime)
			if err != nil {
				return err
			}
			myApps.Append(app)
		}
		return nil
	}, "select id, app_name, app_full_path, icon, create_time, access_time from "+db.appIndexTbl)

	return myApps, err
}

func (db *AppLauncherDB) LoadAppHistory() (*apps.AppList, error) {
	myApps := apps.NewAppList()

	err := db.sqlite.Query(func(rows *sql.Rows) error {
		for rows.Next() {
			var app apps.AppInfo
			err := rows.Scan(&app.AppID, &app.AppName, &app.FullPath, &app.Icon, &app.CreateTime, &app.AccessTime)
			if err != nil {
				return err
			}
			myApps.Append(app)
		}
		return nil
	}, "select id, app_name, app_full_path, icon, create_time, access_time from "+db.appIndexTbl+
		" where access_time > 0 order by access_time desc")

	return myApps, err
}

func (db *AppLauncherDB) AddAppToIndex(apps *apps.AppList) error {
	myApps := apps.GetAppInfo()

	tx, err := db.sqlite.Begin()
	if err != nil {
		return err
	}

	for _, app := range myApps {
		_, _, err := tx.Exec("insert into "+db.appIndexTbl+
			" (app_name, app_full_path, icon, create_time, access_time) values(?, ?, ?, ?, ?)",
			app.AppName, app.FullPath, app.Icon, app.CreateTime, app.AccessTime)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	_ = tx.Commit()
	return nil
}

func (db *AppLauncherDB) UpdateAppInfo(app apps.AppInfo) error {
	_, _, err := db.sqlite.Exec("update "+db.appIndexTbl+" set access_time = ? where id = ?",
		app.AccessTime, app.AppID)
	return err
}
