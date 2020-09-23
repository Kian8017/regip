.PHONY: build
build: # Using go install to make it global for me
	esc -o esc_autogen.go -pkg regip -prefix "web/" web/
	cd regip/ ; go install -race
.PHONY: fmt
fmt:
	go fmt ./...
	go generate ./...
	go test ./...

.PHONY: install
install:
	esc -o esc_autogen.go -pkg regip -prefix "web/" web/
	cd regip/ ; go install -race
