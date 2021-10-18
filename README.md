# ghosts
## Target
* A library for editing system hosts file.
* A tool for adding, deleting, and replacing hosts records.
* A tool to solve DNS pollution of GitHub website. Query the real IP address of domain names such as github.com, and refresh the domain name setting of the system hosts file.
## Points
* Currently only windows and linux systems are supported.
* All records in hosts files will be prioritized by domain name.
* The command to refresh the host DNS in the windows system is: ipconfig /flushdns
* The command to refresh the host DNS in the linux system is: systemctl restart NetworkManager.service
* If the host file has been written successfully and the DNS refresh command failed, try refreshing the DNS or restarting the host yourself.
## How use
1. git clone https://github.com/holimon/ghosts
2. cd ghosts
3. go install ./cmd/github-dns/github-dns.go OR go install ./cmd/ghosts-cli/ghosts-cli.go
### Command for github-dns
> Tips: Administrator permissions must be used to modify the record in order to brush into the hosts file.
* github-dns
### Command for ghosts-cli
```
ghosts-cli

NAME:
   ghosts-cli - A tool for adding, deleting, and replacing hosts records.

USAGE:
   ghosts-cli [global options] command [command options] [arguments...]

VERSION:
   v0.0.1

COMMANDS:
   add      Add a record to the hosts file. Receive 2 arguments.
   del      Remove records from the hosts file.
   replace  Replace a record from the hosts file. Receive 2 arguments.
   resolve  Resolve domains and add records to the hosts file. Receive at least 1 argument.
   show     Show all records in the hosts file. No arguments.
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```
```
ghosts-cli del

NAME:
   ghosts-cli del - Remove records from the hosts file.

USAGE:
   ghosts-cli del command [command options] [arguments...]

COMMANDS:
   field  Remove records by Domain or IP. Receive at least 1 argument.
   index  Remove records by index. Receive at least 1 argument.

OPTIONS:
   --help, -h  show help
```