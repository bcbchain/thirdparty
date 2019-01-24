echo go install goimports gometalinter ...

pushd ..
cd > cwd.txt
set /p cwd=<cwd.txt
del /q cwd.txt
set GOPATH=%cwd%

go install golang.org/x/tools/cmd/goimports
go install github.com/alecthomas/gometalinter
go install github.com/alecthomas/gocyclo
go install golang.org/x/tools/go/analysis/cmd/vet
go install golang.org/x/tools/cmd/gotype
go install github.com/tsenart/deadcode
go install golang.org/x/lint/golint
go install github.com/jgautheron/goconst/cmd/goconst
go install github.com/securego/gosec/cmd/gosec
go install mvdan.cc/interfacer
go install github.com/gordonklaus/ineffassign
go install github.com/kisielk/errcheck
go install github.com/mdempsky/maligned
go install github.com/mdempsky/unconvert
go install github.com/opennota/check/cmd/varcheck
go install github.com/opennota/check/cmd/structcheck
go install honnef.co/go/tools/cmd/megacheck

popd
