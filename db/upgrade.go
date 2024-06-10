package db

import (
	_ "github.com/mattn/go-sqlite3"
)

// 数据库更新列表，列表中数据不能删除，如果有重复update，则将老记录留空，最后添加新记录
// 如果有添加，就往最后添加，不要往前面或中间添加
func DbUpdateSqlList(version string) []string {
	return []string{
		//`update config SET version = ` + version + ` where id = 1`,
		`ALTER TABLE platform ADD COLUMN "hide_name" INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE rom ADD COLUMN "sim_setting" TEXT NOT NULL DEFAULT ''`,
		``,
		`CREATE UNIQUE INDEX "idx_name" ON "config" ("name");`,
		`INSERT INTO config ("name", "data", "desc") VALUES ('GameMultiOpen', '0', '是否允许模拟器游戏多开')`,
	}
}
