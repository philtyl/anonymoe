package setting

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-macaron/session"
	"github.com/philtyl/anonymoe/pkg/bindata"
	"github.com/unknwon/com"
	log "gopkg.in/clog.v1"
	"gopkg.in/ini.v1"
)

var (
	Config Properties
	cfg    *ini.File

	// Session settings
	SessionConfig session.Options
)

type Properties struct {
	// App settings
	AppName   string `json:"app_name"`
	AppVer    string `json:"app_ver"`
	AppURL    string `json:"app_url"`
	AppDomain string `json:"app_domain"`
	AppPath   string `json:"app_path"`
	ProdMode  bool   `json:"prod_mode"`

	// Server settings
	StaticRootPath string `json:"static_root_path"`
	HTTPAddr       string `json:"http_addr"`
	HTTPPort       string `json:"http_port"`
	Protocol       string `json:"protocol"`

	// Database settings
	DatabaseType string `json:"database_type"`
	DatabasePath string `json:"database_path"`

	// Global setting objects
	CfgFilePath string `json:"cfg_file_path"`

	// Mail settings
	MailPort        string   `json:"mail_port"`
	PrivateAccounts []string `json:"private_accounts"`
}

// execPath returns the executable path.
func execPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Abs(file)
}

func init() {
	Config.AppVer = string(bindata.MustAsset("conf/VERSION"))
	var err error
	if Config.AppPath, err = execPath(); err != nil {
		log.Fatal(2, "Fail to resolve [AppPath]: %v\n", err)
	}
	Config.AppPath = strings.Replace(Config.AppPath, "\\", "/", -1)
}

func WorkDir() string {
	i := strings.LastIndex(Config.AppPath, "/")
	if i == -1 {
		return Config.AppPath
	}
	return Config.AppPath[:i]
}

func InstallDir() (dir string) {
	dir = os.Getenv("ANONY_CONFIG")
	if len(dir) == 0 {
		dir = path.Join(WorkDir(), "anonymoe-data")
	}
	return dir
}

func NewContext() (err error) {
	InstallPath := InstallDir()
	Config.CfgFilePath = path.Join(InstallPath, "app.ini")
	cfg, err = ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
	}, bindata.MustAsset("conf/app.ini"))
	if err != nil {
		log.Fatal(2, "Fail to parse default config file [conf/app.ini]: %v", err)
	}

	if com.IsFile(Config.CfgFilePath) {
		if err = cfg.Append(Config.CfgFilePath); err != nil {
			log.Fatal(2, "Fail to load custom config [%s]: %v", Config.CfgFilePath, err)
		}
	} else {
		log.Fatal(0, "Install config [%s] not found, please install server", Config.CfgFilePath)
	}
	cfg.NameMapper = ini.AllCapsUnderscore

	Config.AppName = cfg.Section("").Key("APP_NAME").MustString("Anonymoe")
	Config.ProdMode = cfg.Section("").Key("PRODUCTION_MODE").MustBool(true)

	serverSec := cfg.Section("server")
	Config.AppURL = serverSec.Key("ROOT_URL").MustString("http://localhost:3000")
	Config.Protocol = serverSec.Key("PROTOCOL").String()
	Config.AppDomain = serverSec.Key("DOMAIN").MustString("localhost")
	Config.HTTPAddr = serverSec.Key("HTTP_ADDR").MustString("0.0.0.0")
	Config.HTTPPort = serverSec.Key("HTTP_PORT").MustString("3000")

	dbSec := cfg.Section("database")
	Config.DatabaseType = dbSec.Key("TYPE").MustString("sqlite3")
	Config.DatabasePath = path.Join(InstallPath, dbSec.Key("PATH").MustString("anonymail.db"))

	mailSec := cfg.Section("mail")
	Config.MailPort = mailSec.Key("PORT").MustString("1025")
	Config.PrivateAccounts = strings.Split(mailSec.Key("PRIVATE_ACCOUNTS").MustString("hostmaster"), ",")

	return
}

func IsPrivateAccount(account string) bool {
	for _, a := range Config.PrivateAccounts {
		if account == a {
			return true
		}
	}
	return false
}

func Info() string {
	prettyJSON, err := json.MarshalIndent(Config, "", "    ")
	if err != nil {
		log.Warn("Failed to generate json: %v", err)
		return fmt.Sprintf("Failed to generate json: %v", err)
	}
	return fmt.Sprintf("%s\n", string(prettyJSON))
}
