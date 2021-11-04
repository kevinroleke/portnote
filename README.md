# PortNote

Generate encrypted pastebins or file shares, using IPFS. No third party server required. 

## Notes

- In order to access your PortNotes from public IPFS gateways, your local IPFS node must be connected to the SWARM (may not be the case if you're behind a NAT). 

## Installation

```shell
$ snap install ipfs
$ ipfs init
$ ipfs daemon &
$ 
$ git clone git@github.com:kevinroleke/portnote.git
$ cd portnote
$ go build
```

## Usage

```shell
$ ./portnote --password abc123 --input test.png
$ ./portnote --password 904359043df --eof
TEST PASTE
123
124566
EOF
$ ./portnote --password chr112fjdskjDKFJ39 --input /etc/shadow --daemon localhost:5002
```
