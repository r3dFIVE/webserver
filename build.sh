#!/bin/bash

go build .
sudo chown webserver:webserver ./webserver
