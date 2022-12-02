.PHONY: release
release:
	goreleaser release $(RELEASE_ARGS)
.PHONY: release-snapshot
release-snapshot:
	@make release RELEASE_ARGS="--snapshot"