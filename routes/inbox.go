package routes

import (
	"strconv"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/philtyl/anonymoe/models"
	"github.com/philtyl/anonymoe/pkg/context"
	"github.com/philtyl/anonymoe/pkg/setting"
	log "gopkg.in/clog.v1"
)

const (
	INBOX      = "inbox"
	INBOX_NODE = "model/mailinbox"
)

func NewInbox(c *context.Context) {
	c.Redirect(setting.Config.AppURL + "/inbox/" + strings.ToLower(randomdata.SillyName()))
}

func Inbox(c *context.Context) {
	InboxContents(c)
	c.Success(INBOX)
}

func InboxNode(c *context.Context) {
	InboxContents(c)
	c.Success(INBOX_NODE)
}

func InboxContents(c *context.Context) {
	username := strings.ToLower(c.Params(":user"))
	c.Data["User"] = username

	if setting.Config.ProdMode && setting.IsPrivateAccount(username) {
		c.Data["Private"] = true
	} else {
		mail, err := models.GetMail(username)
		log.Trace("Mail for <%s>: %v", username, mail)
		if err == nil {
			c.Data["Empty"] = len(mail) == 0
			c.Data["Mail"] = mail
		} else {
			log.Warn("Error loading mail for <%s>: %v", username, err)
		}
	}
}

func InboxEmbeddedFile(c *context.Context) {
	contentId := c.Params(":cid")
	mailId, err := strconv.ParseInt(c.Params(":mid"), 10, 64)
	if err == nil {
		embeddedFile, has, err := models.GetEmbeddedFile(mailId, contentId)
		if has && err != nil && embeddedFile != nil {
			log.Trace("Serving file [ID:%d, MailID:%d, ContentID:%s]", embeddedFile.Id, embeddedFile.MailId, embeddedFile.ContentId)
			c.ServeEmbeddedContent(embeddedFile.ContentId, embeddedFile.ContentType, strings.NewReader(embeddedFile.Data))
			return
		}
		log.Warn("Unable to serve file [MailID:%d, ContentID:%s, Has:%b, File:%v]: %v", mailId, contentId, has, embeddedFile, err)
	}
	log.Warn("Unable to load embedded file [MailID:%d, ContentID:%s]: %v", mailId, contentId, err)
	c.NotFound()
}

func InboxAttachment(c *context.Context) {
	c.NotFound() //TODO Attachment implementation
}
