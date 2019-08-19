package setting

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"anonymoe/pkg/bindata"
	"github.com/Unknwon/com"
	"github.com/go-macaron/session"
	log "gopkg.in/clog.v1"
	"gopkg.in/ini.v1"
)

var (
	// App settings
	AppVer      string
	AppName     string
	AppURL      string
	AppDomain   string
	AppPath     string
	AppDataPath string

	// Server settings
	StaticRootPath string
	HTTPAddr       string
	HTTPPort       string
	Protocol       string

	// Global setting objects
	Cfg        *ini.File
	CustomPath string
	CustomConf string

	// Session settings
	SessionConfig session.Options

	// Mail settings
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
	var err error
	if AppPath, err = execPath(); err != nil {
		log.Fatal(2, "Fail to get app path: %v\n", err)
	}
	AppPath = strings.Replace(AppPath, "\\", "/", -1)
}

// WorkDir returns absolute path of work directory.
func WorkDir() (string, error) {
	i := strings.LastIndex(AppPath, "/")
	if i == -1 {
		return AppPath, nil
	}
	return AppPath[:i], nil
}

func NewContext() {
	workDir, err := WorkDir()
	if err != nil {
		log.Fatal(2, "Fail to get work directory: %v", err)
	}

	Cfg, err = ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
	}, bindata.MustAsset("conf/app.ini"))
	if err != nil {
		log.Fatal(2, "Fail to parse 'conf/app.ini': %v", err)
	}

	CustomPath = os.Getenv("ANONY_CUSTOM")
	if len(CustomPath) == 0 {
		CustomPath = workDir + "/custom"
	}

	if len(CustomConf) == 0 {
		CustomConf = CustomPath + "/conf/app.ini"
	}

	if com.IsFile(CustomConf) {
		if err = Cfg.Append(CustomConf); err != nil {
			log.Fatal(2, "Fail to load custom conf '%s': %v", CustomConf, err)
		}
	} else {
		log.Warn("Custom config '%s' not found, ignore this if you're running first time", CustomConf)
	}
	Cfg.NameMapper = ini.AllCapsUnderscore

	homeDir, err := com.HomeDir()
	if err != nil {
		log.Fatal(2, "Fail to get home directory: %v", err)
	}
	homeDir = strings.Replace(homeDir, "\\", "/", -1)

	AppName = Cfg.Section("").Key("APP_NAME").MustString("Anonymoe")

	serverSec := Cfg.Section("server")
	AppURL = serverSec.Key("ROOT_URL").MustString("http://localhost:3000")
	Protocol = serverSec.Key("PROTOCOL").String()
	AppDomain = serverSec.Key("DOMAIN").MustString("localhost")
	HTTPAddr = serverSec.Key("HTTP_ADDR").MustString("0.0.0.0")
	HTTPPort = serverSec.Key("HTTP_PORT").MustString("3000")

	mailSec := Cfg.Section("mail")
	PrivateAccounts = strings.Split(mailSec.Key("PRIVATE_ACCOUNTS").MustString("hostmaster"), ",")
}

func IsPrivateAccount(account string) bool {
	for _, a := range PrivateAccounts {
		if account == a {
			return true
		}
	}
	return false
}
