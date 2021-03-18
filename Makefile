.PHONY: build
build:
	go build .

.PHONY: image
image: build
	docker build -t ghcr.io/ovrclk/ismyaccountfucked .
