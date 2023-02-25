HOST=$(shell hostname)

ifeq ($(HOST), air)
include .env.dev
export $(shell sed 's/=.*//' .env.dev)
endif

dev:
	xcaddy run

test-env:
	echo "cookiename=$$COOKIE_NAME"

build: xcaddy
	xcaddy build --with github.com/TuftsUniversity/caddy-token-auth

xcaddy:
	go install github.com/caddyserver/xcaddy/cmd/xcaddy@latest

run: caddy Caddyfile
	./caddy run ./Caddyfile

clean:
	rm -f caddy
