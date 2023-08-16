#!/bin/sh

trap 'sudo networksetup -setsocksfirewallproxystate Wi-fi off' 2

sudo networksetup -setsocksfirewallproxy Wi-fi localhost 9050
sudo networksetup -setsocksfirewallproxystate Wi-fi on

tor
