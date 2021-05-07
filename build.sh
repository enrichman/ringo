#!/bin/bash

set -e

rm -rf build

gox -output 'build/{{.OS}}_{{.Arch}}/ringo' \
    --os "linux darwin windows" \
    -arch "386 amd64" \
    -osarch '!darwin/386'

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
