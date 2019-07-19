echo go install goimports gometalinter ...

pushd ..
cd > cwd.txt
set /p cwd=<cwd.txt
del /q cwd.txt
set GOPATH=%cwd%

go install golang.org/x/tools/cmd/goimports
go install github.com/golangci/golangci-lint/cmd/golangci-lint

popd
