# grapes [![Go Report Card](https://goreportcard.com/badge/github.com/yaronsumel/grapes)](https://goreportcard.com/report/github.com/yaronsumel/grapes) [![Build Status](https://travis-ci.org/yaronsumel/grapes.svg?branch=master)](https://travis-ci.org/yaronsumel/grapes) [![Build status](https://ci.appveyor.com/api/projects/status/fnepp81rdi8prawn/branch/master?svg=true)](https://ci.appveyor.com/project/yaronsumel/grapes/branch/master) [![GoDoc](https://godoc.org/github.com/yaronsumel/grapes?status.svg)](https://godoc.org/github.com/yaronsumel/grapes)

grapes is lightweight tool designed to distribute commands over ssh with ease.

### Update (25/04/2019)
 
Handshake validation is now in place in order to fix `CVE-2017-3204`, The validation will use the built-in fingerprint list `~/.ssh/known_hosts` as default. 
 
In order to add your ssh server fingerprint to `known_hosts` run the following:

    $ ssh-keyscan -H YOURHOST.COM >> ~/.ssh/known_hosts

### Installation ###

  Run (golang v1.10+ required):

    $ export GO111MODULE=on; go get -u github.com/yaronsumel/grapes

### Usage ###

 Example:

    $ grapes -c config.yml -i ~/.ssh/id_rsa -s prod -cmd whats_up --async

* use the --help flag for full usage output.

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
 
> ##### Written and Maintained by [@YaronSumel](https://twitter.com/yaronsumel) #####
