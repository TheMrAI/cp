#!/usr/bin/bash

# give it a .cpp file and it will compile it and place the result in the same folder
# with .out extension

filename=${1}
clang-format --style=file -i "${filename}"
g++ -std=c++23 -O3 -Wall "${filename}" -o "${filename%.*}.out"
# clang-tidy "${filename}"
