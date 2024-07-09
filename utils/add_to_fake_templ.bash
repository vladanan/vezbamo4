#ls
#echo "/*" >> tmp/fake.templ
#date +%F" "%T%n"*/" >> tmp/fake.templ

OPEN="/*"
CLOSE="*/"
DATE=$(date +%F" "%T)
echo $OPEN$DATE$CLOSE >> tmp/fake.templ
#echo "add to fake templ"

cd utils
# npx esbuild ../src/react/*.jsx --outdir=../assets/assignments/ --minify --bundle --platform=node --global-name=bundle
npx esbuild ../src/react/*.* --outdir=../assets/assignments/ --minify --bundle --global-name=bundle
cd ..