all:
	goreleaser release --snapshot --rm-dist

.PHONY: $(PLATFORMS)
