# simple_transaction_app
Backend for simple transaction application 

Pre-requisites:
Download make(to run make file) and docker in local machine

Steps:
1. Start docker and run below command
`make start`

It will start all docker containers required for app. Example DB and APP.

2. Create tables in DB. For this, run the migration file
`make run-migration`

3. Start the application now
`go run main.go`

If any issues, run go get ., go mod tidy, go mod vendor

4. Test the APIs with postman collection. Import it from here.
`https://api.postman.com/collections/24956101-29fcc38b-95a8-4d85-8806-df9ce9f20d8d?access_key=PMAT-01J4K61PBJN11CJ3MDVKYWNNF7`

Note:
- Change the configurations from config.yml if needed
- Application is designed in such a way thaty it can be extended to any database not only limited to PGSQL. You can spin up MySQL in docker compose and change configs in config.yml and run the app to try it out.
- Make file and docker containers are deployable to cloud if needed to share with team to test
- All components use interfaces rather than structs, hence very flexible to change any component in app without breaking entire app
- wire is used for dependency management. If you are changing dependency, run `make wire` after this

Thank you for trying it out.

