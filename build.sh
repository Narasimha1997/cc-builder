#!/bin/bash

export GOPATH=$(pwd)
#echo "$GOPATH"

package=github.com/Narasimha1997/ccbuilder

go get ${package}
go install ${package}

#mv cchelper bin/