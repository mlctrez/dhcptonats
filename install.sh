#!/bin/bash

cd /opt/dhcptonats

mkdir -p /opt/dhcptonats/logs
ln -sf /opt/dhcptonats/supervisord.include /etc/supervisor/conf.d/dhcptonats.conf

mkdir -p /opt/dhcptonats/gopath
export GOPATH=/opt/dhcptonats/gopath


go get github.com/nats-io/nats
go get github.com/krolaw/dhcp4

go build -o /opt/dhcptonats/dhcptonats main.go


