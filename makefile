all:
	@echo building...
	@go build
	@mkdir -p bin
	@mv estuary bin/
