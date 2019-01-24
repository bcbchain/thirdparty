echo get goimports gometalinter ...

pushd ..
cd > cwd.txt
set /p cwd=<cwd.txt
del /q cwd.txt
set GOPATH=%cwd%

git clone https://github.com/alecthomas/gometalinter.git  ./src/github.com/alecthomas/gometalinter
git clone https://github.com/alecthomas/gocyclo.git       ./src/github.com/alecthomas/gocyclo
git clone https://github.com/golang/tools.git             ./src/golang.org/x/tools
git clone https://github.com/tsenart/deadcode.git         ./src/github.com/tsenart/deadcode
git clone https://github.com/golang/lint.git              ./src/golang.org/x/lint
git clone https://github.com/jgautheron/goconst.git       ./src/github.com/jgautheron/goconst
git clone https://github.com/securego/gosec.git           ./src/github.com/securego/gosec
git clone https://github.com/nbutton23/zxcvbn-go.git      ./src/github.com/nbutton23/zxcvbn-go
git clone https://github.com/ryanuber/go-glob.git         ./src/github.com/ryanuber/go-glob
git clone https://github.com/mvdan/interfacer.git         ./src/mvdan.cc/interfacer
git clone https://github.com/mvdan/lint.git               ./src/mvdan.cc/lint
git clone https://github.com/gordonklaus/ineffassign.git  ./src/github.com/gordonklaus/ineffassign
git clone https://github.com/kisielk/errcheck.git         ./src/github.com/kisielk/errcheck
git clone https://github.com/kisielk/gotool.git           ./src/github.com/kisielk/gotool
git clone https://github.com/mdempsky/maligned.git        ./src/github.com/mdempsky/maligned
git clone https://github.com/mdempsky/unconvert.git       ./src/github.com/mdempsky/unconvert
git clone https://github.com/opennota/check.git           ./src/github.com/opennota/check
git clone https://github.com/dominikh/go-tools.git        ./src/honnef.co/go/tools

popd
