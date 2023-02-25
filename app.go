package caddy_token_auth

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

// TODO: name should be token_auth or basic_token_auth
func init() {
	caddy.RegisterModule(Middleware{})
	httpcaddyfile.RegisterHandlerDirective("token_auth", parseCaddyfile)
}

// Holds all the module's data
type Middleware struct {
	CookieName string `json:"cookie_name,omitempty"`
	FailureUrl string `json:"failure_url,omitempty"`
	CheckUrl   string `json:"check_url,omitempty"`
	CheckUser  string `json:"check_user,omitempty"`
	CheckPass  string `json:"check_pass,omitempty"`
}

// CaddyModule returns the Caddy module information.
func (Middleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.token_auth",
		New: func() caddy.Module { return new(Middleware) },
	}
}

// Provision implements caddy.Provisioner.
func (m *Middleware) Provision(ctx caddy.Context) error {
	/*
		switch m.Output {
		case "stdout":
			m.w = os.Stdout
		case "stderr":
			m.w = os.Stderr
		default:
			m.w = os.Stdout
			//return fmt.Errorf("an output stream is required")
		}
	*/
	return nil
}

// Validate implements caddy.Validator.
func (m *Middleware) Validate() error {
	// TODO
	/*
		if m.w == nil {
			return fmt.Errorf("no writer")
		}
	*/
	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	//m.w.Write([]byte(r.RemoteAddr))
	//return next.ServeHTTP(w, r)
	//return nil

	caddy.Log().Sugar().Infof("looking for cookie name=%s", m.CookieName)
	cookie, err := r.Cookie(m.CookieName)
	to := m.FailureUrl
	if err != nil {
		msg := fmt.Sprintf("error reading cookie: %s", err)
		caddy.Log().Info(msg)
		http.Redirect(w, r, to, http.StatusPermanentRedirect)
		return nil
	}

	caddy.Log().Info(fmt.Sprintf("token=%s", cookie.Value))
	ok, err := isValidToken(cookie.Value, m)
	if err != nil {
		msg := fmt.Sprintf("error processing token: %s", err)
		caddy.Log().Info(msg)
		http.Redirect(w, r, to, http.StatusPermanentRedirect)
		return nil
	}
	if !ok {
		msg := "token is not valid"
		caddy.Log().Info(msg)
		http.Redirect(w, r, to, http.StatusPermanentRedirect)
		return nil
	}

	caddy.Log().Info("token is valid")
	return next.ServeHTTP(w, r)
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (m *Middleware) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		// token value
		parameter := d.Val()
		// rest of params
		args := d.RemainingArgs()
		switch parameter {
		case "cookie_name":
			if len(args) != 1 {
				return d.Err("invalid cookie_name")
			}
			m.CookieName = args[0]
		case "failure_url":
			if len(args) != 1 {
				return d.Err("invalid failure_url")
			}
			m.FailureUrl = args[0]
		case "check_url":
			if len(args) != 1 {
				return d.Err("invalid check_url")
			}
			m.CheckUrl = args[0]
		case "check_user":
			if len(args) != 1 {
				return d.Err("invalid check_user")
			}
			m.CheckUser = args[0]
		case "check_pass":
			if len(args) != 1 {
				return d.Err("invalid check_pass")
			}
			m.CheckPass = args[0]
		default:
			//d.Err("Unknow cam parameter: " + parameter)
			caddy.Log().Sugar().Info("skipping: %s %v", parameter, args)
		}
	}
	return nil
}

// parseCaddyfile unmarshals tokens from h into a new Middleware.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m Middleware
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return m, err
}

func isValidToken(token string, m Middleware) (bool, error) {
	client := &http.Client{Timeout: time.Second * 2}

	body := bytes.NewBufferString(fmt.Sprintf("token=%s", token))
	req, err := http.NewRequest(http.MethodPost, m.CheckUrl, body)
	if err != nil {
		return false, err
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(m.CheckUser, m.CheckPass)

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	caddy.Log().Info(fmt.Sprintf("check statuscode=%d", resp.StatusCode))
	if resp.StatusCode != 200 {
		return false, nil
	}
	return true, nil
}

// Interface guards
var (
	_ caddy.Provisioner           = (*Middleware)(nil)
	_ caddy.Validator             = (*Middleware)(nil)
	_ caddyhttp.MiddlewareHandler = (*Middleware)(nil)
	_ caddyfile.Unmarshaler       = (*Middleware)(nil)
)
