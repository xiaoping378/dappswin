#! /bin/bash

set -e
set -u
set -x

downloadAndInstallEos()
{
    apt update && apt install -y wget

    wget https://github.com/EOSIO/eos/releases/download/v1.5.2/eosio_1.5.2-1-ubuntu-18.04_amd64.deb

    dpkg -i eosio_1.5.2-1-ubuntu-18.04_amd64.deb || apt install -f -y && dpkg -i eosio_1.5.2-1-ubuntu-18.04_amd64.deb
}

type cleos >/dev/null 2>&1 || downloadAndInstallEos

rm /usr/local/sbin/dappswin && cp -f dappswin /usr/local/sbin/dappswin
mkdir -p /etc/dappswin && cp dappswin.toml /etc/dappswin/

rm /etc/systemd/system/dappswin.service && cp -f ./scripts/systemd/dappswin.service /etc/systemd/system/dappswin.service

systemctl daemon-reload
systemctl start dappswin
systemctl enable dappswin
