#ls
#echo "/*" >> tmp/fake.templ
#date +%F" "%T%n"*/" >> tmp/fake.templ

OPEN="/*"
CLOSE="*/"
DATE=$(date +%F" "%T)
echo $OPEN$DATE$CLOSE >> tmp/fake.templ
#echo "add to fake templ"

# build go binary, valja povremeno da se uradi radi testiranja
#go build -o bin src/main.go

cd utils_dev
# npx esbuild ../src/react/*.jsx --outdir=../assets/assignments/ --minify --bundle --platform=node --global-name=bundle
npx esbuild ../src/ext/react/*.* --outdir=../assets/assignments/ --minify --bundle --global-name=bundle
cd ..