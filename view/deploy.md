# dev
```sh

export versionTag="c"

cd /home/devsz/chaoshe_data/data_frontend
git pull

npm run build:dev
npm run build:prod-${versionTag}
docker build -t chaoshe/data/frontend:${versionTag}_latest .

cd ../
docker-compose up -d

docker tag chaoshe/data/frontend:${versionTag}_latest hub.chaosheapi.com:5001/chaoshe/data/frontend:${versionTag}_latest
docker push hub.chaosheapi.com:5001/chaoshe/data/frontend:${versionTag}_latest

docker tag chaoshe/data/frontend:${versionTag}_latest hub.yiku.miyouhudong.com/blind_box/data/frontend:${versionTag}_latest
docker push hub.yiku.miyouhudong.com/blind_box/data/frontend:${versionTag}_latest
```


# release
```sh

cd /home/www/chaoshe_data/

docker-compose pull
docker-compose up -d
```
