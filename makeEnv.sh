
echo "Setting up .env variables"

read -p "API_PORT [80]: " APP_PORT
read -p "API_SECRET: " API_SECRET
read -p "DB_HOST [db]: " DB_HOST
read -p "DB_PORT [3306]: " DB_PORT
read -p "DB_USER [user]: " DB_USER
read -p "DB_PASS [password]:" DB_PASS
read -p "DB_NAME [main_db]: " DB_NAME
read -p "DB_CHARSET [uft8mb4]: " DB_CHARSET
read -p "DB_LOC [local]: " DB_LOC
read -p "TOKEN_EXPIRE [3600]: " TOKEN_EXPIRE

rm ./.env.temp || true
touch ./.env.temp

echo APP_PORT=$APP_PORT >> ./.env.temp
echo API_SECRET=$API_SECRET >> ./.env.temp
echo DB_HOST=$DB_HOST >> ./.env.temp
echo DB_PORT=$DB_PORT >> ./.env.temp
echo DB_USER=$DB_USER >> ./.env.temp
echo DB_PASS=$DB_PASS >> ./.env.temp
echo DB_NAME=$DB_NAME >> ./.env.temp
echo DB_CHARSET=$DB_CHARSET >> ./.env.temp
echo DB_LOC=$DB_LOC >> ./.env.temp
echo TOKEN_EXPIRE=$TOKEN_EXPIRE >> ./.env.temp

rm ./.env.temp || true