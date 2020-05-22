#!/bin/bash
git clone https://github.com/gcp/libogg src/libogg
git clone https://gitlab.xiph.org/xiph/opus.git src/opus
git clone https://gitlab.xiph.org/xiph/opusfile.git src/opus-file
git clone https://gitlab.xiph.org/xiph/libopusenc.git src/opus-libopusenc

rm -rf src/libogg/.git
rm -rf src/opus/.git
rm -rf src/opus-file/.git
rm -rf src/opus-libopusenc/.git