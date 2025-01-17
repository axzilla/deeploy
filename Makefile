.PHONY: app landing

app:
	cd internal/app && make

web:
	cd internal/web && make

dev-app:
	cd internal/app && make dev

dev-web:
	cd internal/web && make dev

debug-cli:
	tmux new-window -n deeploy_debug 'dlv debug --headless --api-version=2 --listen=127.0.0.1:43000 ./cmd/cli; tmux kill-window -t deeploy_debug'

