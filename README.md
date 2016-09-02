# grapes
grapes is lightweight tool designed to distribute commands over ssh with ease.

 Installation and usage
 ----------------------

 To install it, run:

     go get -u github.com/yaronsumel/grapes

 Usage Example :

     grapes -c config.yml -i ~/.ssh/id_rsa -s prod -command whats_up

config structure:

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
