package template

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/philtyl/anonymoe/pkg/setting"
	"github.com/philtyl/anonymoe/pkg/tool"
)

func NewFuncMap() []template.FuncMap {
	return []template.FuncMap{map[string]interface{}{
		"AppName": func() string {
			return setting.Config.AppName
		},
		"AppURL": func() string {
			return setting.Config.AppURL
		},
		"AppVer": func() string {
			return setting.Config.AppVer
		},
		"AppDomain": func() string {
			return setting.Config.AppDomain
		},
		"LoadTimes": func(startTime time.Time) string {
			return fmt.Sprint(time.Since(startTime).Nanoseconds()/1e6) + "ms"
		},
		"DateFmtLong": func(t time.Time) string {
			return t.Format(time.RFC1123Z)
		},
		"DateFmtShort": func(t time.Time) string {
			return t.Format("Jan 02, 2006")
		},
		"MD5": func(str string) string {
			m := md5.New()
			m.Write([]byte(str))
			return hex.EncodeToString(m.Sum(nil))
		},
		"Str2HTML": func(raw string) template.HTML {
			return template.HTML(raw) //TODO Sanatize
		},
		"HumanTimeSince": func(then int64) string {
			return tool.HumanTimeSince(then)
		},
	}}
}

func EscapePound(str string) string {
	return strings.NewReplacer("%", "%25", "#", "%23", " ", "%20", "?", "%3F").Replace(str)
}
