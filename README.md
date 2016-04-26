# go-cli
Go CLI is a template project, that you can use to build your
CLI Apps.

## Usage

To create a new CLI application, you can do the following steps

    APP=<appname>
    DST=$GOPATH/src/<repo-hoster>/<user>/$APP
    mkdir -p $DST
    go get -u github.com/szuecs/go-cli
    rsync -a --exclude=.git $GOPATH/src/github.com/szuecs/go-cli/ $DST
    cd $DST
    grep -rl go-cli | xargs sed -i "s@go-cli@$APP@g"
    mv cmd/go-cli cmd/$APP


The main package and function of the CLI app is in
cmd/$APP/main.go. It parses flags and merges the configuration to
start the application.

Configuring your client, use the following make target:

    % make config