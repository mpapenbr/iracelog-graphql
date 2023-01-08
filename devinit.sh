go install github.com/spf13/cobra-cli@latest
go install github.com/goreleaser/goreleaser@latest
go install github.com/caarlos0/svu@latest
go install github.com/99designs/gqlgen@latest
go install github.com/cweill/gotests/gotests@v1.6.0

if [ -f setuplinks.sh ]; then
    . ./setuplinks.sh
fi