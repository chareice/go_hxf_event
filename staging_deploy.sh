#!/bin/sh
. config-default.sh

ssh root@staging.m.hxf.chareice.com << EOF

docker pull $imageName
docker stop $containerName
docker rm $containerName
docker run --link hxf_mongo:mongo \
           --link hxf_redis:redis \
           --name $containerName -d $imageName
EOF
