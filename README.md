# README

# Installation

```
go get github.com/tmtk75/kii-cli
```

or

Download binary from <https://github.com/tmtk75/kii-cli/releases>


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

    kii-cli auth login

OK, you're ready to run kii-cli. Try `log` subcommand.

    kii-cli log

You might not see any outputs, in the case, type next.

    echo "function main() {}" > foobar
    kii-cli servercode deploy foobar

Then try `log` again and you can see a line that showed you deployed a servercode.

    kii-cli log
    0001: &{servercode.file.deploy INFO Server Code File deployed 2014-07-01 14:37:27.92 +0000 UTC}

TBD


# Usage

```
NAME:
   kii-cli - KiiCloud command line interface

USAGE:
   kii-cli [global options] command [command options] [arguments...]

VERSION:
   0.1.2

COMMANDS:
   auth         Authentication
   app          Application management
   log          Disply logs for an app
   servercode   Server code management
   user         User management
   bucket       Bucket management
   object       Object management
   dev          Development support
   help, h      Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --app-id             AppID
   --app-key            AppKey
   --client-id          ClientID
   --client-secret      ClientSecret
   --site               us,jp,cn,sg
   --endpoint-url       Site URL
   --verbose            Verbosely
   --profile 'default'  Profile name for ~/.kii/config
   --curl               Print curl command saving body as a tmp file if body exists
   --help, -h           show help
   --version, -v        print the version
```

And supports flat-style subcommands.
```
FLAT=1 kii-cil
COMMANDS:
   auth:login                   Login as AppAdmin
   auth:info                    Print login info
   app:config                   Print config of app
   log                          Disply logs for an app
   servercode:list              List versions of server code
   servercode:deploy            Deploy a server code
   servercode:get               Get specified server code
   ...

```
