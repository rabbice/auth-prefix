package authprefix

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(AuthPrefix{})
}

type AuthPrefix struct {
	Prefix string `json:"prefix,omitempty"`
	logger *zap.Logger
}

func (AuthPrefix) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "http.handlers.auth_prefix",
		New: func() caddy.Module { return new(AuthPrefix) },
	}
}

func (p *AuthPrefix) Provision(ctx caddy.Context) error {
	p.logger = ctx.Logger(p)
	return nil
}

func (p *AuthPrefix) Validate() error {
	if p.Prefix == "" {
		p.Prefix = "."
	}
	return nil
}

func (p AuthPrefix) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	for _, part := range strings.Split(r.URL.Path, "/") {
		if strings.HasPrefix(part, p.Prefix) {
			http.Error(w, "Not Found", http.StatusNotFound)
			if p.logger != nil {
				p.logger.Debug(fmt.Sprintf(
					"authorized prefix: %q in %s", part, r.URL.Path))
				}
				return nil
			}
		}
		return next.ServeHTTP(w, r)
	}

	var (
		_ caddy.Provisioner = (*AuthPrefix)(nil)
		_ caddy.Validator = (*AuthPrefix)(nil)
		_ caddyhttp.MiddlewareHandler = (*AuthPrefix)(nil)
	)
