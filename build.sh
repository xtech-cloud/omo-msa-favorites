go install omo.msa.favorite
mkdir _build
mkdir _build/bin
export GOPROXY=https://goproxy.cn/
cp -rf /root/go/bin/omo.msa.favorite _build/bin/
cp -rf conf _build/
cd _build
tar -zcf msa.favorite.tar.gz ./*
mv msa.favorite.tar.gz ../
cd ../
rm -rf _build
