# caddy-token-auth

Caddy module to implement token authentication middlware for
the Caddy [server](https://github.com/caddyserver/caddy).

## Why

If you use something like [SAML](https://www.cloudflare.com/learning/access-management/what-is-saml/)
to authenticate your users against your services you know it has some limitations. One of them is
that building a [service provider](https://www.cloudflare.com/learning/access-management/what-is-saml/)
can be a laborious task. It really depends on your infrastructure and how automated it is. In my
environment (I won't get into details)
the process of adding a service provider to the
[Identity provider](https://www.cloudflare.com/learning/access-management/what-is-saml/) is... complicated.
Also, the authentication you grant to the users last as long as the SAML session
is active.

The realization I had was: ok, creating a SP is painful, let's do it only once and reuse the authentication
mechanism. Concretely, once the user successfully completes a SAML flow and we have established trust, we
can issue a token. Other services can then check for that token to validate access.

That last part is where this project comes in. This Caddy module implements that token authentication by
extracting the token from the session (cookie) and checking against a service to see if the token is valid.
If it is then the middleware stack continues, probably proxying traffic to the final service.

## dev

Make sure you have [xcaddy](https://github.com/caddyserver/xcaddy) installed.

`$ make`

That will build caddy with this module so you can test locally

Use `make build` if you want to generate a caddy binary with this module on it.
