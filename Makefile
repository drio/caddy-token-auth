caddy: xcaddy
	./xcaddy build --with github.com/TuftsUniversity/caddy-token-auth

run: caddy Caddyfile
	./caddy run ./Caddyfile
