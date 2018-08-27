#!/bin/bash

if ! which brew >/dev/null; then
  echo "homebrew is not available. Install it from http://brew.sh"
  exit 1
else
  echo "homebrew already installed"
fi

if ! which go >/dev/null; then
  echo "installing go..."
  brew install go
else
  echo "go already installed"
fi

echo "all dependencies installed."
