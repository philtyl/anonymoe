package models

import (
	"fmt"
	"os"
	"path"

	"anonymoe/pkg/setting"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/core"
)

// Engine represents a XORM engine or session.
type Engine interface {
	Get(interface{}) (bool, error)
	Insert(...interface{}) (int64, error)
	Count(interface{}) (int64, error)
}

var (
	x      *xorm.Engine
	tables []interface{}

	DbCfg struct {
		Type, Path string
	}
)

func init() {
	tables = append(tables, new(User), new(Mail), new(MailRecipient))
}

func LoadConfigs() {
	DbCfg.Type = setting.DatabaseType
	DbCfg.Path = setting.DatabasePath
}

func getEngine() (*xorm.Engine, error) {
	LoadConfigs()
	if err := os.MkdirAll(path.Dir(DbCfg.Path), os.ModePerm); err != nil {
		return nil, fmt.Errorf("create directories: %v", err)
	}
	return xorm.NewEngine(DbCfg.Type, "file:"+DbCfg.Path+"?cache=shared&mode=rwc")
}

func SetEngine() (err error) {
	LoadConfigs()
	x, err = getEngine()
	if err != nil {
		return fmt.Errorf("connect to database: %v", err)
	}

	x.SetMapper(core.GonicMapper{})
	x.ShowSQL(true)
	return nil
}

func NewEngine() (err error) {
	if err = SetEngine(); err != nil {
		return err
	}

	if err = x.StoreEngine("InnoDB").Sync2(tables...); err != nil {
		return fmt.Errorf("sync structs to database tables: %v\n", err)
	}

	return nil
}
