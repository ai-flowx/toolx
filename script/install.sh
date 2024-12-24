#!/bin/bash

# Install python
sudo add-apt-repository ppa:deadsnakes/ppa
sudo apt install python3.13 python3.13-dev pkg-config

# Install package
pip3 install -r requirements.txt
