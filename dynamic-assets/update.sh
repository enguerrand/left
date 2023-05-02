#!/bin/bash

echo "1.0.$(git rev-parse --short HEAD)" > dynamic-assets/version.txt
