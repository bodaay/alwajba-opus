#!/bin/bash
brew install autoconf
brew install automake

current=$PWD
outputFolder="$current/lib/darwin/x64"
outputFolderDynamic="$current/lib/darwin/x64/dylib"
mkdir -p $outputFolder
mkdir -p $outputFolderDynamic

echo "Building libogg library"
cd $current
cd src/libogg
./autogen.sh
./configure --disable-doc --disable-examples
make clean
make
cp ./src/.libs/*.dylib ./src/.libs/*.a $outputFolder/
# make clean
echo "Finished Building libogg library"


echo "Building Opus Core library"
cd $current
cd src/opus
./autogen.sh
./configure --disable-doc --disable-examples
make clean
make
cp ./.libs/*.dylib ./.libs/*.a $outputFolder/
# make clean
echo "Finished Building Opus Core library"

echo "Building Opus File library"
cd $current
cd src/opus-file
./autogen.sh
DEPS_LIBS="-L$outputFolder -lopus -logg" DEPS_CFLAGS="-I$current/src/opus/include -I$current/src/libogg/include" ./configure --disable-http  --disable-doc --disable-examples --disable-http
make clean
make
cp ./.libs/*.dylib ./.libs/*.a $outputFolder/
# make clean
echo "Finished Building Opus File library"


echo "Building Opus libopusenc library"
cd $current
cd src/opus-libopusenc
./autogen.sh
DEPS_LIBS="-L$outputFolder  -lopus -logg lopusfile" DEPS_CFLAGS="-I$current/src/opus/include  -I$current/src/libogg/include" ./configure  --disable-doc --disable-examples
make clean
make
cp ./.libs/*.dylib ./.libs/*.a $outputFolder/
# make clean
echo "Finished Building Opus libopusenc library"


mv $outputFolder/*.dylib $outputFolderDynamic/
