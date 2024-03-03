#/usr/bin/sh
#Make .env file
if [ ! -f ./.env ]
then
    sh makeEnv.sh
fi

sudo docker-compose build --no-cache
sudo docker-compose up -d 