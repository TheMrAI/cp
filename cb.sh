#!/usr/bin/bash

# give it a .cpp file and it will compile it and place the result in the same folder
# without any extensions

filename=${1}
clang-format --style=file -i "${filename}"
g++ -std=c++20 -O3 -Wall "${filename}" -o "${filename%.*}"
# clang-tidy "${filename}"