#!/bin/sh

DESCRIPTION=$(git describe --tags --abbrev=100 --dirty) &&
echo "${DESCRIPTION##*-g}"
