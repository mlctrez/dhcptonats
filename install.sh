#!/bin/bash

# install dhcptonats

set -x

export TOP=/opt/dhcptonats
export PKG=github.com/mlctrez/dhcptonats

mkdir -p $TOP/logs
mkdir -p $TOP/gopath
export GOPATH=$TOP/gopath
go get $PKG

cp $GOPATH/src/$PKG/supervisord.include $TOP
cp $GOPATH/bin/dhcptonats $TOP
ln -sf $TOP/supervisord.include /etc/supervisor/conf.d/dhcptonats.conf

supervisorctl reread
supervisorctl add dhcptonats