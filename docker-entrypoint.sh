#!/bin/sh

cd /opt/FlagField || exit 1

while ! mysql -h mysql -uroot -p${MYSQL_ROOT_PASSWORD} -e 'exit'; do
    sleep 1
    echo 'Waiting for mysql...'
done

if ! mysql -h mysql -uroot -p${MYSQL_ROOT_PASSWORD} -e "use ${DATABASE_NAME}"; then
    mysql -h mysql -uroot -p${MYSQL_ROOT_PASSWORD} -e "create database ${DATABASE_NAME}" || exit 1
    ./dist/migrator -template=initial || exit 1
    ./dist/manager user add --username "${SYSTEM_ADMIN_USERNAME}" --password "${SYSTEM_ADMIN_PASSWORD}" --email "${SYSTEM_ADMIN_EMAIL}" --admin || exit 1
    ./dist/manager config set --key "system.setup_time" --val "$(date -u "+%Y-%m-%dT%H:%M:%SZ")" || exit 1
fi

while ! redis-cli -h redis -p 6379 -r 1 ping; do
    sleep 1
    echo 'Waiting for redis...'
done

./dist/server || exit 1
