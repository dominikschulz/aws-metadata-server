# aws-metdata-server

aws-metadata-server is a very simple AWS compatible meta-data / user-data
server for local development. It does not support any kind of security and
SHOULD NOT be run anywhere remotely publicly accessible. It's only meant to
be run on a developers machine, preferable in a virtual machine with
limited network scope.

## Usage

* Build with ```go build```
* Add the well known IP 169.254.169.254, e.g. with iproute2 (linux): ```ip route add 169.254.169.254/32 dev eth0```
* Start as root (so it can bind to Port 80, sorry): ```sudo ./aws-metadata-server```
* Set the keys you need to their desired values, e.g. ```curl -i -X POST -d 'us-east-1a' http://169.254.169.254/latest/meta-data/placement/availability-zone```

