# grapes [![Go Report Card](https://goreportcard.com/badge/github.com/yaronsumel/grapes)](https://goreportcard.com/report/github.com/yaronsumel/grapes) [![Build Status](https://travis-ci.org/yaronsumel/grapes.svg?branch=master)](https://travis-ci.org/yaronsumel/grapes) [![Build status](https://ci.appveyor.com/api/projects/status/fnepp81rdi8prawn/branch/master?svg=true)](https://ci.appveyor.com/project/yaronsumel/grapes/branch/master) [![GoDoc](https://godoc.org/github.com/yaronsumel/grapes?status.svg)](https://godoc.org/github.com/yaronsumel/grapes)

grapes is lightweight tool designed to distribute commands over ssh with ease.

> ##### Latest Binaries [Windows](https://github.com/yaronsumel/grapes/releases/download/v0.2/win-grapes.zip), [Mac](https://github.com/yaronsumel/grapes/releases/download/v0.2/darwin-grapes.zip), [Linux](https://github.com/yaronsumel/grapes/releases/download/v0.2/linux-grapes.7z) #####
 
### Installation and usage ###

 To install it, run:

     $ go get -u github.com/yaronsumel/grapes

 Usage Example :

     $ grapes -c config.yml -i ~/.ssh/id_rsa -s prod -cmd whats_up

```
$ grapes --help 

//      ____ __________ _____  ___  _____
//     / __  / ___/ __  / __ \/ _ \/ ___/
//    / /_/ / /  / /_/ / /_/ /  __(__  )
//    \__, /_/   \__,_/ .___/\___/____/
//   /____/          /_/ v 0.1.6 // Yaron Sumel [yaronsu@gmail.com]
//

Usage of ./grapes:
  -async
        async - when true, parallel executing over servers
  -c string
        config file - yaml config file
  -cmd string
        command name - name of the command to run
  -i string
        identity file - path to private key
  -s string
        server group - name of the server group
  -y    force yes
```

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
