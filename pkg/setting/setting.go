package setting

import (
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
	// App settings
	AppName   string
	AppVer    string
	AppURL    string
	AppDomain string
	AppPath   string
	ProdMode  bool

	// Server settings
	StaticRootPath string
	HTTPAddr       string
	HTTPPort       string
	Protocol       string

	// Database settings
	DatabaseType string
	DatabasePath string

	// Global setting objects
	Cfg         *ini.File
	CfgFilePath string

	// Session settings
	SessionConfig session.Options

	// Mail settings
	MailPort        string
	PrivateAccounts []string
)

// execPath returns the executable path.
func execPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Abs(file)
}

func init() {
	AppVer = string(bindata.MustAsset("conf/VERSION"))
	var err error
	if AppPath, err = execPath(); err != nil {
		log.Fatal(2, "Fail to get app path: %v\n", err)
	}
	AppPath = strings.Replace(AppPath, "\\", "/", -1)
}

func WorkDir() string {
	i := strings.LastIndex(AppPath, "/")
	if i == -1 {
		return AppPath
	}
	return AppPath[:i]
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
	CfgFilePath = path.Join(InstallPath, "app.ini")
	Cfg, err = ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
	}, CfgFilePath)
	if err != nil {
		log.Fatal(2, "Fail to parse 'app.ini': %v", err)
	}

	if com.IsFile(CfgFilePath) {
		if err = Cfg.Append(CfgFilePath); err != nil {
			log.Fatal(2, "Fail to load custom config '%s': %v", CfgFilePath, err)
		}
	} else {
		log.Fatal(2, "Install config '%s' not found, please install server", CfgFilePath)
	}
	Cfg.NameMapper = ini.AllCapsUnderscore

	AppName = Cfg.Section("").Key("APP_NAME").MustString("Anonymoe")
	ProdMode = Cfg.Section("").Key("PRODUCTION_MODE").MustBool(true)

	serverSec := Cfg.Section("server")
	AppURL = serverSec.Key("ROOT_URL").MustString("http://localhost:3000")
	Protocol = serverSec.Key("PROTOCOL").String()
	AppDomain = serverSec.Key("DOMAIN").MustString("localhost")
	HTTPAddr = serverSec.Key("HTTP_ADDR").MustString("0.0.0.0")
	HTTPPort = serverSec.Key("HTTP_PORT").MustString("3000")

	dbSec := Cfg.Section("database")
	DatabaseType = dbSec.Key("TYPE").MustString("sqlite3")
	DatabasePath = path.Join(InstallPath, dbSec.Key("PATH").MustString("anonymail.db"))

	mailSec := Cfg.Section("mail")
	MailPort = mailSec.Key("PORT").MustString("1025")
	PrivateAccounts = strings.Split(mailSec.Key("PRIVATE_ACCOUNTS").MustString("hostmaster"), ",")

	return err
}

func IsPrivateAccount(account string) bool {
	for _, a := range PrivateAccounts {
		if account == a {
			return true
		}
	}
	return false
}
