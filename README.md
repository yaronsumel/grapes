# grapes [![Go Report Card](https://goreportcard.com/badge/github.com/yaronsumel/grapes)](https://goreportcard.com/report/github.com/yaronsumel/grapes) [![Build Status](https://travis-ci.org/yaronsumel/grapes.svg?branch=master)](https://travis-ci.org/yaronsumel/grapes) [![Build status](https://ci.appveyor.com/api/projects/status/fnepp81rdi8prawn/branch/master?svg=true)](https://ci.appveyor.com/project/yaronsumel/grapes/branch/master) [![GoDoc](https://godoc.org/github.com/yaronsumel/grapes?status.svg)](https://godoc.org/github.com/yaronsumel/grapes)

grapes is lightweight tool designed to distribute commands over ssh with ease.

### Installation ###

 grab binary or run (golang required):

     $ go get -u github.com/yaronsumel/grapes

> Os  | Binary | md5sum
> ------------- | ------------- | -------------
> windows  | [download](https://github.com/yaronsumel/grapes/releases/download/v0.2/win-grapes.zip) | b18c8e4f511329e5d4b9a27bd8aa52c7
> linux  | [download](https://github.com/yaronsumel/grapes/releases/download/v0.2/linux-grapes.7z) | 60734d1004c9266f2cc4987453262931
> darwin  | [download](https://github.com/yaronsumel/grapes/releases/download/v0.2/darwin-grapes.zip) | d89264f774f50a39379ec46c1865e286

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
