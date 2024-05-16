#ls
echo "/*" >> tmp/fake.templ
date +%F" "%T%n"*/" >> tmp/fake.templ
#echo "add to fake templ"
cd utils
npx esbuild ../src/react/*.jsx --outdir=../assets/assignments/ --minify --bundle --platform=node
cd ..