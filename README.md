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
   0.1.0

COMMANDS:
   login			Login as AppAdmin
   login:info			Print login info
   log				Disply logs for an app
   servercode:list		List versions of server code
   servercode:deploy		Deploy a server code
   servercode:get		Get specified server code
   servercode:invoke		Invoke an entry point of server code
   servercode:activate		Activate a version
   servercode:delete		Delete a version of server code
   servercode:hook-attach	Attach a hook config to current or specified server code
   servercode:hook-get		Get hook the config of current or specified server code
   servercode:hook-delete	Delete the hook config of current specified server code
   servercode:list-executions	List executions for 7 days before
   bucket:list			List buckets
   bucket:acl			Show a bucket ACL
   user:login			Login as a user
   user:create			Create a user
   object:create		Create an object in application scope
   object:read			Read the object in application scope
   object:replace		Replate the object in application scope with a new one
   object:delete		Delete the object in application scope
   server			WebSocket echo server for testing
   help, h			Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --app-id 		AppID
   --app-key 		AppKey
   --client-id 		ClientID
   --client-secret 	ClientSecret
   --site 		us,jp,cn,sg
   --endpoint-url 	Site URL
   --verbose		Verbosely
   --profile 'default'	Profile name for ~/.kii/config
   --curl		Print curl command saving body as a tmp file if body exists
   --help, -h		show help
   --version, -v	print the version
   
```

