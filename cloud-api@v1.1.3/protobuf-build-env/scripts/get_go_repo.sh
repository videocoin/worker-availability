#/bin/bash

# This is a helper script that uses go get to download a specific version of a go package
# Not all required dependencies have adpoted Go 1.11 modules so it is not possible to use
# a go.mod to get all the dependencies
# This script just uses go get to get the latest version of a given package and then switches that
# repository to a given tag

# usage [script_name] [git url_to_go_package] [yes/no] [tag/branch to switch to]
#                                                |
#                                                |
#                                                V
#                                             whether or not to explicitly build and install the package after checking out a particular tagj

GIT_TAG=$3
NEEDS_RECOMPILATION_AND_REINSTALL=$2
GIT_REPO_URL=$1

# Download the packages but do not build/install so that first
# the right tag/commit can be checked out
go get -v -d $GIT_REPO_URL

# --quiet to supress some  messages about detached head
git -C "$(go env GOPATH)/src/$GIT_REPO_URL" checkout $GIT_TAG

if [ $NEEDS_RECOMPILATION_AND_REINSTALL == "yes" ]; then
    go install -v $GIT_REPO_URL
fi

