package storage

import (
	"database/sql"
	"path/filepath"
	"sid-desktop/base"
	"sid-desktop/desktop/common"
	"strings"
	"sync"
)

type AppFavoritesDB struct {
	sqlite base.SQLite
	dbName string

	appFavoritesTbl string
}

var (
	appFavoritesDB   *AppFavoritesDB
	appFavoritesOnce sync.Once
)

func GetAppFavoritesDB() *AppFavoritesDB {
	appFavoritesOnce.Do(func() {
		appFavoritesDB = &AppFavoritesDB{
			dbName:          "sid_app_favorites.db",
			appFavoritesTbl: "app_favorites",
		}
	})
	return appFavoritesDB
}

func (db *AppFavoritesDB) Open(rootPath string) error {
	return db.sqlite.Open(filepath.Join(rootPath, db.dbName))
}

func (db *AppFavoritesDB) Close() {
	_ = db.sqlite.Close()
}

func (db *AppFavoritesDB) Init() error {
	if err := db.createAppFavorites(); err != nil {
		return err
	}
	return nil
}

func (db *AppFavoritesDB) createAppFavorites() error {
	sqlStrArray := []string{
		"drop table if exists " + db.appFavoritesTbl + ";",
		"CREATE TABLE IF NOT EXISTS " + db.appFavoritesTbl + `(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(256) NOT NULL,
			url VARCHAR(512) NOT NULL,
			tags VARCHAR(256) NOT NULL,
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

func (db *AppFavoritesDB) NeedInit() (bool, error) {
	tables, err := db.sqlite.ListTables()
	if err != nil {
		return false, err
	}

	if base.ContainsString(tables, db.appFavoritesTbl) < 0 {
		return true, nil
	}

	return false, nil
}

func (db *AppFavoritesDB) LoadFavorites() (*common.FavoritesList, error) {
	favors := common.NewFavoritesList()

	err := db.sqlite.Query(func(rows *sql.Rows) error {
		for rows.Next() {
			var fav common.Favorites
			tagStr := ""
			err := rows.Scan(&fav.ID, &fav.Name, &fav.Url, &tagStr, &fav.CreateTime, &fav.AccessTime)
			if err != nil {
				return err
			}
			fav.Tags = strings.Split(tagStr, common.FavorTagSep)
			favors.Append(fav)
		}
		return nil
	}, "select id, name, url, tags, create_time, access_time from "+db.appFavoritesTbl+" order by access_time desc")

	return favors, err
}

func (db *AppFavoritesDB) AddFavorites(favor common.Favorites) error {
	_, _, err := db.sqlite.Exec("insert into "+db.appFavoritesTbl+
		" (name, url, tags, create_time, access_time) values(?, ?, ?, ?, ?)",
		favor.Name, favor.Url, strings.Join(favor.Tags, common.FavorTagSep), favor.CreateTime, favor.AccessTime)
	return err
}

func (db *AppFavoritesDB) AddFavoritesList(favors *common.FavoritesList) error {
	favs := favors.GetFavorites()

	tx, err := db.sqlite.Begin()
	if err != nil {
		return err
	}

	for _, favor := range favs {
		_, _, err := tx.Exec("insert into "+db.appFavoritesTbl+
			" (name, url, tags, create_time, access_time) values(?, ?, ?, ?, ?)",
			favor.Name, favor.Url, strings.Join(favor.Tags, common.FavorTagSep), favor.CreateTime, favor.AccessTime)

		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	_ = tx.Commit()
	return nil
}

func (db *AppFavoritesDB) RmFavorites(favor common.Favorites) error {
	_, _, err := db.sqlite.Exec("delete from "+db.appFavoritesTbl+" where id = ?", favor.ID)
	return err
}

func (db *AppFavoritesDB) UpdateFavorites(favor common.Favorites) error {
	_, _, err := db.sqlite.Exec("update "+db.appFavoritesTbl+" set name = ?, url = ?, tags = ?, access_time = ? where id = ?",
		favor.Name, favor.Url, strings.Join(favor.Tags, common.FavorTagSep), favor.AccessTime, favor.ID)
	return err
}
