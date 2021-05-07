#!/bin/bash

set -e

rm -rf build

LATEST_TAG=$(git describe --tags)

gox -output 'build/{{.OS}}_{{.Arch}}/ringo' \
    --os "linux darwin windows" \
    -arch "386 amd64" \
    -osarch '!darwin/386' \
    -ldflags="-X 'main.Version=${LATEST_TAG}'"

for rls in build/{linux,darwin}*; do \
    tar czf build/ringo-$(echo ${rls} | cut -f2 -d/).tgz -C ${rls} ringo; \
done

for rls in build/windows*; do \
    mv ${rls}/ringo.exe build/ringo-$(echo ${rls} | cut -f2 -d/).exe
done

for rls in build/{linux,darwin,windows}*; do \
    rm -rf ${rls}
done

ls -l build/
