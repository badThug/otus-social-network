#
# docker hints
#
preview:
	docker-compose up

docker_network:
	docker network create otus-social-network

force_recreate:
	docker-compose up --build --force-recreate $(n)

#
# Service users
#
users_db_init:
	docker exec -i -t osn__users_mysql-master sh -c "mysql -uroot -pmysql dbase < /init_sql/initdb.sql"
	docker exec -i -t osn__users_mysql-master sh -c "mysql -uroot -pmysql dbase < /init_sql/mock_users.sql"

users_db_clear:
	@echo '------master------'
	sudo rm -rf service-users/database/master/backup/*
	touch service-users/database/master/backup/.gitkeep
	sudo rm -rf service-users/database/master/data/*
	touch service-users/database/master/data/.gitkeep
	sudo rm -rf service-users/database/master/log/*
	touch service-users/database/master/log/.gitkeep

users_db_fix_rights:
	sudo chmod 777 -R service-users/database/master
	sudo chmod 0444  service-users/database/master/conf.d/master.cnf

users_backend_fix_rights:
	sudo chmod 777 -R service-users/backend/logs

#
# Service chat
#
chat_db_init: chat_db_init_single_db

chat_db_init_single_db:
	docker exec -i -t osn__chat_mongo  sh -c "mongo < /init/default_collections_init"

# chat db sharded version
chat_db_init_sharded_db:
	docker exec -i -t osn__chat_mongocfg sh -c "mongo < /init/cfg"
	docker exec -i -t osn__chat_mongo_shard1  sh -c "mongo < /init/shard1"
	docker exec -i -t osn__chat_mongo_shard2  sh -c "mongo < /init/shard2"
	@echo "------ wait 10-15 sec till mongos will get update from config ------"
	sleep 15
	docker exec -i -t osn__chat_mongos  sh -c "mongo < /init/mongos"
	docker exec -i -t osn__chat_mongos  sh -c "mongo < /init/default_collections_init"
	docker exec -i -t osn__chat_mongos  sh -c "mongo < /init/sharding_init"

chat_db_clear:
	@echo '------shard1------'
	sudo rm -rf service-chat/database/shard1/data/*
	touch service-chat/database/shard1/data/.gitkeep
	sudo rm -rf service-chat/database/shard1/logs/*
	touch service-chat/database/shard1/logs/.gitkeep
	@echo '------shard2------'
	sudo rm -rf service-chat/database/shard2/data/*
	touch service-chat/database/shard1/data/.gitkeep
	sudo rm -rf service-chat/database/shard2/logs/*
	touch service-chat/database/shard1/logs/.gitkeep
	@echo '-------cfg-------'
	sudo rm -rf service-chat/database/cfg/data/*
	touch service-chat/database/cfg/data/.gitkeep
	sudo rm -rf service-chat/database/cfg/logs/*
	touch service-chat/database/cfg/logs/.gitkeep

chat_db_fix_rights:
	sudo chmod 777 -R service-chat/database/shard1/data
	sudo chmod 777 -R service-chat/database/shard1/logs
	sudo chmod 777 -R service-chat/database/shard2/data
	sudo chmod 777 -R service-chat/database/shard2/logs
	sudo chmod 777 -R service-chat/database/cfg/data
	sudo chmod 777 -R service-chat/database/cfg/logs

chat_backend_fix_rights:
	sudo chmod 777 -R service-chat/backend/logs

#
# overall services
#
build-dbs:
	docker-compose up osn__users_mysql-master osn__chat_mongos

#
# RabbitMQ
#
rabbit-clear:
	sudo rm -rf rabbitmq/data/*
	touch rabbitmq/data/.gitkeep
	sudo rm -rf rabbitmq/logs/*
	touch rabbitmq/logs/.gitkeep
	sudo rm -rf rabbitmq/etc/*
	touch rabbitmq/etc/.gitkeep

rabbit-fix-rights:
	sudo chmod 777 -R rabbitmq/data
	sudo chmod 777 -R rabbitmq/logs
	sudo chmod 777 -R rabbitmq/etc

db_init: users_db_init chat_db_init
fix_rights: users_db_fix_rights users_backend_fix_rights chat_db_fix_rights chat_backend_fix_rights rabbit-fix-rights
fresh_run: docker_network fix_rights


