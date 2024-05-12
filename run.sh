#!/bin/bash


if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <filename>"
    exit 1
fi

docker run -v "$(pwd)"/"$1":/root/"$1" yadro-computer-club "$1"