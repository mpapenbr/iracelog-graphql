#
# Arguments: all arguments are passed to svu
# Tasks
# - verify tests are running
# - verify goreleaser with "--snapshot --skip-publish --rm-dist"
# - replace cmd/version via sed
# 
FULL_TAG=$(svu $@)
NO_PREFIX_TAG=$(svu $@ --strip-prefix | head -1)
echo $FULL_TAG
echo $NO_PREFIX_TAG

# updating the version in graph/version.go
go test ./... -p 1 && \
goreleaser --snapshot --skip-publish --clean && \
sed -i -E "s/(Version\W+=\W*)\"(.*?)\"/\1\"$NO_PREFIX_TAG\"/" graph/version.go && \
git add graph/version.go && git commit -m "chore: Release $FULL_TAG" && \
git tag -a $FULL_TAG -m "Release $FULL_TAG" && \
git push origin main --tags
