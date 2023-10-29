.PHONY: static

static:
	go build --ldflags '-extldflags "-static -lwebp -lm -lpthread"' -tags netgo