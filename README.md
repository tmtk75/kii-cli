# README

# Installation

```
go get github.com/tmtk75/kii-cli
```

# Getting Started

Create config file in `~/.kii/config` as follow.

    [default]
    app_id = aaaaa...
    app_key = bbbbbbbbbbbbbbbbb...
    client_id = xxxxxxxxxxxxxx...
    client_secret = yyyyyyyyyyyyyyyyyyyyyyyyy...
    site = us

First, login to get your app admin access token.
If you encounter an error message, please make sure your credentials in `default` section of `~/.kii/config`.

    kii-cli login

OK, you're ready to run kii-cli. Try `log` subcommand.

    kii-cli log

You might not see any outputs, in the case, type next.

    touch foobar
    kii-cli servercode:deploy foobar

TBD


# Usage

```
NAME:
   kii-cli - KiiCloud command line interface

USAGE:
   kii-cli [global options] command [command options] [arguments...]

VERSION:
   0.0.4

COMMANDS:
   login                Login as AppAdmin
   login:info           Print login info
   log                  Disply logs for an app
   servercode:list      List versions of server code
   servercode:deploy    Deploy a server code
   servercode:get       Get specified server code
   servercode:invoke    Invoke an entry point of server code
   servercode:activate  Activate a version
   servercode:delete    Delete an entry point of server code
   bucket:list          List buckets
   bucket:acl           Show a bucket ACL
   users:create         Create user
   server               WebSocket echo server for testing
   help, h              Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --app-id             AppID
   --app-key            AppKey
   --client-id          ClientID
   --client-secret      ClientSecret
   --site               us,jp,cn,sg
   --endpoint-url       Site URL
   --verbose            Verbosely
   --profile 'default'  Profile name for ~/.kii/config
   --version, -v        print the version
   --help, -h           show help
```
