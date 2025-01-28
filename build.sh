#!/bin/bash

go install github.com/akavel/rsrc@latest

rsrc -ico assets/icon.ico -o rsrc.syso

go build -ldflags "-H=windowsgui"

rm rsrc.syso