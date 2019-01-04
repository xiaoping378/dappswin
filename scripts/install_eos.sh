#! /bin/bash

set -e
set -u
set -x


wget https://github.com/EOSIO/eos/releases/download/v1.5.2/eosio_1.5.2-1-ubuntu-18.04_amd64.deb

dpkg -i eosio_1.5.2-1-ubuntu-18.04_amd64.deb

