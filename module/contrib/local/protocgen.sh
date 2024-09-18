#!/usr/bin/env bash

set -eo pipefail

GOPATH=${GOPATH:-$(go env GOPATH)}
if [ -z $GOPATH ]; then
	echo "GOPATH not set!"
	exit 1
fi

if [[ $PATH != *"$GOPATH/bin"* ]]; then
	echo "GOPATH/bin must be added to PATH"
	exit 1
fi

protoc_install_proto_gen_doc() {
  echo "Installing protobuf protoc-gen-doc plugin"
  (go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest 2> /dev/null)
}

echo "Generating gogo proto code"
cd proto
proto_dirs=$(find ./ -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do  
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    if grep "option go_package" $file &> /dev/null ; then
      buf generate --template buf.gen.gogo.yml $file
    fi
  done
done

protoc_install_proto_gen_doc

echo "Generating proto docs"
buf generate --template buf.gen.doc.yml

cd ..

# move proto files to the right places
cp -r github.com/Gravity-Bridge/Gravity-Bridge/module/* ./
rm -rf github.com
