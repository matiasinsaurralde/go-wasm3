#!/bin/sh
wasicc cstring.c -Wl,--export-all,--allow-undefined -o cstring.wasm
