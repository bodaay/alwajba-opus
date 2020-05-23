#!/bin/bash
sudo apt-get install git autoconf automake libtool libssl-dev gcc make

current=$PWD
outputFolder="$current/lib/linux/x64"
outputFolderDynamic="$current/lib/linux/x64/so"
mkdir -p $outputFolder
mkdir -p $outputFolderDynamic

echo "Building libogg library"
cd $current
cd src/libogg
./autogen.sh
./configure --disable-doc --disable-examples
make clean
make
cp ./src/.libs/*.so ./src/.libs/*.a $outputFolder/
# make clean
echo "Finished Building libogg library"


echo "Building Opus Core library"
cd $current
cd src/opus
./autogen.sh
./configure --disable-doc --disable-examples
make clean
make
cp ./.libs/*.so ./.libs/*.a $outputFolder/
# make clean
echo "Finished Building Opus Core library"

echo "Building Opus File library"
cd $current
cd src/opus-file
./autogen.sh
DEPS_LIBS="-L$outputFolder -lopus -logg" DEPS_CFLAGS="-I$current/src/opus/include -I$current/src/libogg/include" ./configure --disable-http  --disable-doc --disable-examples --disable-http
make clean
make
cp ./.libs/*.so ./.libs/*.a $outputFolder/
# make clean
echo "Finished Building Opus File library"


echo "Building Opus libopusenc library"
cd $current
cd src/opus-libopusenc
./autogen.sh
DEPS_LIBS="-L$outputFolder -lopus -logg lopusfile" DEPS_CFLAGS="-I$current/src/opus/include  -I$current/src/libogg/include" ./configure  --disable-doc --disable-examples
make clean
make
cp ./.libs/*.so ./.libs/*.a $outputFolder/
# make clean
echo "Finished Building Opus libopusenc library"



mv $outputFolder/*.so $outputFolderDynamic/
