dev:
	xcaddy run

build: xcaddy
	xcaddy build --with github.com/TuftsUniversity/caddy-token-auth

xcaddy:
	go install github.com/caddyserver/xcaddy/cmd/xcaddy@latest

run: caddy Caddyfile
	./caddy run ./Caddyfile

clean:
	rm -f caddy
