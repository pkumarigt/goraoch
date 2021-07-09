### Clone the repo
`git clone git@github.com:engineyard/goroach.git`

### Run application locally
`docker-compose up -d`

### Validate application
`curl http://127.0.0.1:8880/quote/`

### Commit app container to image
Run docker ps and note the id of app container
`docker commit <id> <dockerhubusername>/goroach`
example: docker commit a42a74dfe65f sergeyabrahamyandf/goroach

### Push to dockerhub
`docker push  <dockerhubusername>/goroach`
example: docker push sergeyabrahamyandf/goroach

### Create new database in EYK
```
Navigate to https://eyk.ey.io/app/databases/add-database
Cluster: <select your cluster>
Database Name: tesdb
Database Username: goroachuser
Engine: postgres
Engine Version: 9.6
Database Storage: 20
Instance Size: Small

Once database is ready note the Host and Password

```

### Install eyk cli
Follow instructions from https://support.cloud.engineyard.com/hc/en-us/articles/360057913834-Download-the-Kontainers-CLI-Tool

### Login to eyk
Navigate to https://eyk.ey.io/app/clusters, copy the CLI login command and run the command in your terminal
example: eyk ssologin https://eyk.lab-two.ey-dedicated-internal.ey.io

### Create app in eyk
`eyk apps:create goroach --no-remote`

### Set environment variables for the app
`eyk config:set PORT=8880 SERVICE_PORT=8880 DB_SERVER=svc-rdb-testdb20210528092111872700000002.cyerxl1scpna.ap-southeast-2.rds.amazonaws.com DB_PORT=5432 DB_USER=goroachuser DB_DATABASE=testdb DB_PASSWORD=kdsJRcNCXMyGxMYStb7QwcJhFQ9J9PmYP9N0YBJv DEFAULT_PAGE_SIZE=20 -a goroach`

### Deploy the app using docker image created earlier
`eyk builds:create sergeyabrahamyandf/goroach:latest -a goroach --procfile='web: /main'`

### Test the app
run `eyk info` and note the url. Access https://url/quote
example: https://goroach.lab-two.ey-dedicated-internal.ey.io/quote/
