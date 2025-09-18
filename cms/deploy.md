# dev
```sh

export versionTag="c"

cd /home/devsz/chaoshe_data/data_backend/
git pull
cd /home/devsz/chaoshe_data/data_backend/cmd/cms_backend

# 执行打包 & docker build
./packer_build.sh ${versionTag}

# 复制 dev 配置 
cp configs/config.yaml.dev configs/config.yaml

# 启动docker -- 改为使用docker-compose
cd ../../../
docker-compose up -d

# 验证功能

# 推送 docker 至 正式服务器
docker tag chaoshe/data/backend:${versionTag}_latest hub.chaosheapi.com:5001/chaoshe/data/backend:${versionTag}_latest
docker push hub.chaosheapi.com:5001/chaoshe/data/backend:${versionTag}_latest
```


# release
```sh

cd /home/www/chaoshe_data/

docker-compose pull
docker-compose up -d
```
