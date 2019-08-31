package models

import (
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/microcosm-cc/bluemonday"
	"github.com/philtyl/anonymoe/pkg/setting"
	"xorm.io/core"
)

// Engine represents a XORM engine or session.
type Engine interface {
	Get(interface{}) (bool, error)
	Insert(...interface{}) (int64, error)
	Count(interface{}) (int64, error)
}

var (
	x       *xorm.Engine
	tables  []interface{}
	policy  *bluemonday.Policy
	cleaner *Cleaner

	DbCfg struct {
		Type, Path string
	}
)

func init() {
	tables = append(tables, new(Attachment), new(EmbeddedFile), new(Mail), new(MailRecipient), new(User))
	policy = bluemonday.UGCPolicy()
	cleaner = &Cleaner{
		tags: map[string]*regexp.Regexp{
			"div": regexp.MustCompile(`(?s)<div(.*)</div>`),
		},
	}
}

func LoadConfigs() {
	DbCfg.Type = setting.Config.DatabaseType
	DbCfg.Path = setting.Config.DatabasePath
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

type Cleaner struct {
	tags map[string]*regexp.Regexp
}

func (c *Cleaner) clean(s string) string {
	for search, regex := range c.tags {
		for regex.MatchString(s) {
			s = regex.ReplaceAllString(s, fmt.Sprintf("<MATCHED%s$1</MATCHED%s>", search, search))
		}
		openTags := regexp.MustCompile(fmt.Sprintf(`(?s)<%s.*>`, search))
		closeTags := regexp.MustCompile(fmt.Sprintf(`<\%s>`, search))
		undoMatcheTags := regexp.MustCompile(fmt.Sprintf(`(?s)<MATCHED%s(.*)</MATCHED%s>`, search))
		s = openTags.ReplaceAllString(s, "")
		s = closeTags.ReplaceAllString(s, "")
		for undoMatcheTags.MatchString(s) {
			s = undoMatcheTags.ReplaceAllString(s, fmt.Sprintf("<%s$1</%s>", search, search))
		}
	}
	return s
}
