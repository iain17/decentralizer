#!/bin/bash -x

git remote add $1 git@github.com:ucirello/$1.git
git subtree -q push --prefix=$1 $1 master 
