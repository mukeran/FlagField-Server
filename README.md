# FlagField-Server

Aiming at being a full-functioned CTF platform.

## First Run

### Build by yourself

You have to install build-essentials, go>=1.13, MySQL>=5.7 (8.0 recommended), redis-server (5.0 recommended) first.  

1. `git clone https://github.com/FlagField/FlagField-Server.git && cd FlagField-Server`
2. `make`
3. `./dist/setup` - Using setup tool to validate and configure.
4. `./dist/server` - Run server

### Using Docker

1. `git clone https://github.com/FlagField/FlagField-Server.git && cd FlagField-Server`
2. Edit the docker-compose.yml File to specify the port and volume the MySQL file, redis file, config.json, etc.;
3. `make tools` - (**You don't have to use this command now. See explanation below.**)
4. `./dist/setup --non-validation` - Using setup tool in non-validation mode to generate a "config.json" file (**not supported yet**, please copy config.example.json and edit by yourself. **Check if upload directory is created**);
5. Running `make up` to generate and run the docker image.

## Make Commands

|Command|Description|
|-------|-----------|
|`make`|Test, compile and build|
|`make test`|Test|
|`make build`|Compile and build all (including server, migrator, setup and manager)|
|`make tools`|Compile and build tools (including migrator, setup and manager)|
|`make clean`|Clean output files|
|`make up`|Run docker-compose up --build|
|`make down`|Run docker-compose down|
|`make start`|Run docker-compose start|
|`make stop`|Run docker-compose stop|

## Tools

### migrator

Migrator is a tool to initialize or migrate database tables.  
It has several templates (currently initial only) to do the migration.

### setup

Setup can generate the config file and execute migrator to initialize the database.

### manager

Manager can help manage users, contests, etc. in command line. But now, it is only able to list and add user.