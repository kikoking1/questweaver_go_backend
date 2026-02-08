# Mongo
## Linux Installation
https://www.mongodb.com/docs/v7.0/tutorial/install-mongodb-on-ubuntu/

## Start Mongo
sudo systemctl start mongod

## Check Mongo Running
sudo systemctl status mongod

## Reload Mongo
sudo systemctl restart mongod

## Verify Started Correctly
sudo systemctl status mongod

## Auto Start Mongo on System Restart
sudo systemctl enable mongod

## Stop Mongo 
sudo systemctl stop mongod

## Mongo Logs location
/var/log/mongodb/mongod.log

## Connect to Mongo
mongosh