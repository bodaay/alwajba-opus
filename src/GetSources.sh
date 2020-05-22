#!/bin/bash
git clone https://github.com/gcp/libogg libogg
git clone https://gitlab.xiph.org/xiph/opus.git opus
git clone https://gitlab.xiph.org/xiph/opusfile.git opus-file
git clone https://gitlab.xiph.org/xiph/libopusenc.git opus-libopusenc

rm -rf libogg/.git
rm -rf opus/.git
rm -rf opus-file/.git
rm -rf opus-libopusenc/.git