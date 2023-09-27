package adapters

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"its.ac.id/base-go/bootstrap/config"
	"its.ac.id/base-go/pkg/session"
)

type Cookie struct {
}

// TODO: Encrypt cookie data
func NewCookie() *Cookie {
	return &Cookie{}
}

func (c *Cookie) Get(ctx *gin.Context, sessionId string) (*session.Data, error) {
	raw, err := ctx.Cookie(sessionId)
	if err != nil && err == http.ErrNoCookie {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	b64, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(b64), &data); err != nil {
		return nil, err
	}
	session := session.NewData(ctx, sessionId, data, c)

	return &session, nil
}

func (c *Cookie) Save(ctx *gin.Context, sessionId string, data map[string]interface{}) error {
	cfg := do.MustInvoke[config.Config](do.DefaultInjector).Session()
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	b64 := base64.StdEncoding.EncodeToString([]byte(json))
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(sessionId, string(b64), cfg.Lifetime, cfg.CookiePath, cfg.Domain, cfg.Secure, true)
	return nil
}

func (c *Cookie) Delete(ctx *gin.Context, sessionId string) error {
	cfg := do.MustInvoke[config.Config](do.DefaultInjector).Session()
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(sessionId, "", -1, cfg.CookiePath, cfg.Domain, cfg.Secure, true)
	return nil
}
