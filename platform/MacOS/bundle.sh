# Script for building bundle
src=../..
app=frequon.app

# Remove old bundle
rm -rf ./$app

# Copy the Info.plist
mkdir -p $app/Contents
cp Info.plist $app/Contents
 
# Build and copy executable
exe=$app/Contents/MacOS/frequon-invaders-2.2
(cd $src; go build -tags='release bundle') 
mkdir $app/Contents/MacOS
cp $src/FrequonInvaders $exe

# Copy dynamic libraries and deal with rpath stuff
frame=$app/Contents/Frameworks
mkdir $frame
for f in libSDL2-2.0.0.dylib libSDL2_ttf-2.0.0.dylib
do
    cp /usr/local/lib/$f $frame/$f
    install_name_tool -id @rpath/$f $frame/$f
    install_name_tool -change /usr/local/lib/$f @rpath/$f $exe
done
install_name_tool -add_rpath @loader_path/../Frameworks $exe

# Copy resources
mkdir $app/Contents/Resources
cp $src/{Roboto-Regular.ttf,unicons.1.0.ttf,Characters.png} $app/Contents/Resources/
rm -f frequon.icns
iconutil -c icns frequon.iconset
cp frequon.icns $app/Contents/Resources/
