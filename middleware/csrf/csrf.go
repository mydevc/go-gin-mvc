package csrf

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dchest/uniuri"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go-gin-mvc/utils"
	"io"
	"strings"
)

const (
	csrfSecret = "csrfSecret"
	csrfSalt   = "csrfSalt"
	csrfToken  = "csrfToken"
)

var IgnoreAction []string

var defaultIgnoreMethods = []string{"GET", "HEAD", "OPTIONS"}

var defaultErrorFunc = func(c *gin.Context) {
	panic(errors.New("CSRF token mismatch"))
}

var defaultTokenGetter = func(c *gin.Context) string {
	r := c.Request

	if t := r.FormValue("_csrf"); len(t) > 0 {
		return t
	} else if t := r.URL.Query().Get("_csrf"); len(t) > 0 {
		return t
	} else if t := r.Header.Get("X-CSRF-TOKEN"); len(t) > 0 {
		return t
	} else if t := r.Header.Get("X-XSRF-TOKEN"); len(t) > 0 {
		return t
	}

	return ""
}

// Options stores configurations for a CSRF middleware.
type Options struct {
	Secret        string
	IgnoreMethods []string
	ErrorFunc     gin.HandlerFunc
	TokenGetter   func(c *gin.Context) string
}

func init()  {
	fs := utils.CsrfExcept.Section("fuction").Keys()
	for _, f := range fs {
		IgnoreAction = append(IgnoreAction,f.String())
	}
}


func tokenize(secret, salt string) string {
	h := sha1.New()
	io.WriteString(h, salt+"-"+secret)
	hash := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return hash
}

func inArray(arr []string, value string) bool {
	inarr := false

	for _, v := range arr {
		if v == value {
			inarr = true
			break
		}
	}

	return inarr
}

// Middleware validates CSRF token.
func Middleware(options Options) gin.HandlerFunc {
	ignoreMethods := options.IgnoreMethods
	errorFunc := options.ErrorFunc
	tokenGetter := options.TokenGetter

	if ignoreMethods == nil {
		ignoreMethods = defaultIgnoreMethods
	}

	if errorFunc == nil {
		errorFunc = defaultErrorFunc
	}

	if tokenGetter == nil {
		tokenGetter = defaultTokenGetter
	}

	return func(c *gin.Context) {

		fn := c.HandlerName();
		fn = fn[strings.LastIndex(fn,"/"):]

		for _, action := range IgnoreAction {
			if(strings.Contains(fn,action)){
				fmt.Println(action)
				c.Next()
				return
			}
		}

		session := sessions.Default(c)
		c.Set(csrfSecret, options.Secret)

		if inArray(ignoreMethods, c.Request.Method) {
			c.Next()
			return
		}

		salt, ok := session.Get(csrfSalt).(string)

		if !ok || len(salt) == 0 {
			errorFunc(c)
			return
		}

		token := tokenGetter(c)

		if tokenize(options.Secret, salt) != token {
			errorFunc(c)
			return
		}

		c.Next()
	}
}

// GetToken returns a CSRF token.
func GetToken(c *gin.Context) string {
	session := sessions.Default(c)
	secret := c.MustGet(csrfSecret).(string)

	if t, ok := c.Get(csrfToken); ok {
		return t.(string)
	}

	salt, ok := session.Get(csrfSalt).(string)
	if !ok {
		salt = uniuri.New()
		session.Set(csrfSalt, salt)
		session.Save()
	}
	token := tokenize(secret, salt)
	c.Set(csrfToken, token)

	return token
}
