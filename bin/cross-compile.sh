#!/usr/bin/env bash
#Usage: $0 <version>
set -e

# Prepare tmp folder
rm -Rf gen
mkdir -p gen

# Insert version
version="${VERSION:-v0.0.1-latest}.$(git rev-parse --short HEAD).$(date -u +"%Y%m%d%H%M%S")"
echo $version >> "gen/version"

# Generate license notices
deps="github.com/goeuro/ingress-generator-kit $(go list -f '{{ join .Deps "\n"}}' ./... | grep -v 'goeuro/ingress-generator-kit')"
out="gen/LICENSES"
echo -e "OPEN SOURCE LICENSES\n" > $out

for dep in $deps; do
  if [ -d "$GOPATH/src/$dep" ]; then
    notices=$(ls -d $GOPATH/src/$dep/* 2>/dev/null | grep -i -e "license" -e "licence" -e "copying" -e "notice" || echo)
    if [ ! -z "$notices" ]; then
      echo -e "$dep\n\n" >> $out
      for notice in $notices; do
        echo "Adding license: $notice"
        cat $notice >> $out
      done
      echo -e "\n\n" >> $out
    fi
  fi
done

# Compile bindata
go-bindata -o core/bindata.go -nometadata -pkg core gen/

# Cross compile
gox \
  -osarch="darwin/amd64 linux/amd64" \
  -ldflags="-s -w" \
  -output="gen/{{.Dir}}_{{.OS}}_{{.Arch}}"
