package context

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-macaron/session"
	"github.com/philtyl/anonymoe/pkg/template"
	log "gopkg.in/clog.v1"
	"gopkg.in/macaron.v1"
)

type Context struct {
	*macaron.Context
	Flash   *session.Flash
	Session session.Store
}

// HTML responses template with given status.
func (c *Context) HTML(status int, name string) {
	log.Trace("Template: %s", name)
	c.Context.HTML(status, name)
}

// Success responses template with status http.StatusOK.
func (c *Context) Success(name string) {
	c.HTML(http.StatusOK, name)
}

// JSONSuccess responses JSON with status http.StatusOK.
func (c *Context) JSONSuccess(data interface{}) {
	c.JSON(http.StatusOK, data)
}

// RawRedirect simply calls underlying Redirect method with no escape.
func (c *Context) RawRedirect(location string, status ...int) {
	c.Context.Redirect(location, status...)
}

// Redirect responses redirection wtih given location and status.
// It escapes special characters in the location string.
func (c *Context) Redirect(location string, status ...int) {
	c.Context.Redirect(template.EscapePound(location), status...)
}

// Handle handles and logs error by given status.
func (c *Context) Handle(status int, title string, err error) {
	switch status {
	case http.StatusNotFound:
		c.Data["Title"] = "Page Not Found"
	case http.StatusInternalServerError:
		c.Data["Title"] = "Internal Server Error"
		log.Fatal(2, "%s: %v", title, err)
	}
	c.HTML(status, fmt.Sprintf("status/%d", status))
}

// NotFound renders the 404 page.
func (c *Context) NotFound() {
	c.Handle(http.StatusNotFound, "", nil)
}

// ServerError renders the 500 page.
func (c *Context) ServerError(title string, err error) {
	c.Handle(http.StatusInternalServerError, title, err)
}

// NotFoundOrServerError use error check function to determine if the error
// is about not found. It responses with 404 status code for not found error,
// or error context description for logging purpose of 500 server error.
func (c *Context) NotFoundOrServerError(title string, errck func(error) bool, err error) {
	if errck(err) {
		c.NotFound()
		return
	}
	c.ServerError(title, err)
}

func (c *Context) HandleText(status int, title string) {
	c.PlainText(status, []byte(title))
}

func (c *Context) ServeEmbeddedContent(name string, contentType string, r io.ReadSeeker) {
	c.Resp.Header().Set("Content-Type", contentType)
	c.Resp.Header().Set("Expires", "0")
	http.ServeContent(c.Resp, c.Req.Request, name, time.Now(), r)
}

// Contexter initializes a classic context for a request.
func Contexter() macaron.Handler {
	return func(ctx *macaron.Context, sess session.Store, f *session.Flash) {
		c := &Context{
			Context: ctx,
			Flash:   f,
			Session: sess,
		}
		log.Trace("Session ID: %s", sess.ID())
		ctx.Map(c)
	}
}
