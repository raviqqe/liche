#!/bin/sh

set -ex

bundler install
bundler exec cucumber PATH=$PWD:$PATH
