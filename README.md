# grapes [![Go Report Card](https://goreportcard.com/badge/github.com/yaronsumel/grapes)](https://goreportcard.com/report/github.com/yaronsumel/grapes) [![Build Status](https://travis-ci.org/yaronsumel/grapes.svg?branch=master)](https://travis-ci.org/yaronsumel/grapes) [![Build status](https://ci.appveyor.com/api/projects/status/fnepp81rdi8prawn/branch/master?svg=true)](https://ci.appveyor.com/project/yaronsumel/grapes/branch/master) [![GoDoc](https://godoc.org/github.com/yaronsumel/grapes?status.svg)](https://godoc.org/github.com/yaronsumel/grapes)

grapes is lightweight tool designed to distribute commands over ssh with ease.

### Update (25/04/2019)
 
Handshake validation is now in place in order to fix `CVE-2017-3204`, The validation will use the built-in fingerprint list `~/.ssh/known_hosts` as default. 
 
In order to add your ssh server fingerprint to `known_hosts` run the following:

    $ ssh-keyscan -H YOURHOST.COM >> ~/.ssh/known_hosts

### Run with docker 

```bash
docker run -v ~/.grapes.yml:/root/.grapes.yml -v ~/.ssh/id_rsa:/root/.ssh/id_rsa -it docker.pkg.github.com/yaronsumel/grapes/grapes:0.3.0
```

### Installation ###

  Run (Go v1.26+ required):

    $ go install github.com/yaronsumel/grapes@latest

### Usage ###

 Example:

    $ grapes -c config.yml -i ~/.ssh/id_rsa -s prod -cmd whats_up --async

* use the --help flag for full usage output.

Each command result includes its captured stdout and stderr under `std.out` and
`std.err`. An `error` field describes execution or session failures, and
`exit_status` records the remote exit status when available (including zero for
success). Stderr output alone does not mean that a command failed.

Grapes finishes reporting all selected servers before exiting, including in
async mode. It exits with status 1 if any connection, command, or output write
fails, and status 0 when all succeed. Commands continue after a failed command;
each result is reported separately. Failure summaries are written to stderr.

### Config ###

config structure (YAML):

 ```
version: 1
servers:
  prod :
      - name : "prod server #1"
        host : "prod.example.com:22"
        user : "ubuntu"
  staging :
      - name : "staging server #1"
        host : "staging.example.com:22"
        user : "ubuntu"
      - name : "staging server #2"
        host : "staging.example.com:23"
        user : "ubuntu"
commands:
  whats_up :
      - "ls -al /tmp"
      - "date"
  date :
      - "date"
 ```

### Tests ###

    $ go test -race ./...

The SSH tests use a local test server and do not require access to a public server.
 
> ##### Written and Maintained by [@YaronSumel](https://twitter.com/yaronsumel) #####
