templ:
	templ generate --watch --proxy="http://localhost:8090" --open-browser=false -v

server:
	air \
	--build.cmd "go build -o tmp/main ./cmd/web/main.go" \
	--build.bin "tmp/main" \
	--build.delay "100" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true
	--proxy="http://localhost:7331"

tailwind:
	npx tailwindcss -i ./assets/css/input.css -o ./assets/css/output.css --watch

dev:
	make -j3 templ server tailwind
