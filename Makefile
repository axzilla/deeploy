.PHONY: app landing

app:
	cd internal/app && make

web:
	cd internal/web && make

dev-app:
	cd internal/app && make dev

dev-web:
	cd internal/web && make dev
