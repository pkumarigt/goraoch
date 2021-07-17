### Clone the repo
Open the repo in DevSpaces

### Run application locally
`docker-compose up -d`

### Validate application
`curl http://127.0.0.1:8880/quote/`


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
### Login to eyk from terminal ( one time setup )
Request for a non browser based login. Get the client.json contents set it in variable

`gp env EYK_CLIENT='{"username":"nonssouser","ssl_verify":true,"controller":"https://eyk.playground.eyk-central.ey.io","token":"c49a42a2d215df54c7ffc7820d208f0357e35bb5","response_limit":0}'
`eval $(gp env -e)`
`mkdir -p ~/.eyk && echo $EYK_CLIENT > ~/.eyk/client.json`

.gitpod.yml will take care of login next time.

### Setup AWS Account Details
`gp env GP_AWS_ACCESS_ID=xxx`
`gp env GP_AWS_SECRET_KEY=xxx`
`gp env GP_AWS_REGION=us-east-1`

`eval $(gp env -e)`
`mkdir -p ~/.aws`
`echo -e "[default]\naws_access_key_id = $GP_AWS_ACCESS_ID\naws_secret_access_key = $GP_AWS_SECRET_KEY\n" > ~/.aws/credentials`
`echo -e "[default]\nregion = $GP_AWS_REGION\noutput = json\n" > ~/.aws/config`

`aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws`

### Docker build

Set up a goroach repo in ECR

set the ENV according to your ECR seetings

`gp env ECR_REGISTRY=public.ecr.aws/xxxxxxx/goroach`
`gp env ECR_TAG=dev`
`eval $(gp env -e)`

Default is arm64 in devspaces EYK is not arm yet.
`docker build -t goroach . --build-arg TARGETARCH=amd64`
`docker tag goroach "${ECR_REGISTRY}:${ECR_TAG}"`
`docker push "${ECR_REGISTRY}:${ECR_TAG}"`

### Create app in eyk
`eyk apps:create goroachapp --no-remote`

Set the app vieweable in UI. Mostly you email id
`eyk perms:create <user> --app=goroachapp`
### Set environment variables for the app by replacing DB_SERVER, DB_USER and DB_PASSWORD with your own database parameters
`eyk config:set PORT=8880 SERVICE_PORT=8880 DB_SERVER=<Host> DB_PORT=5432 DB_USER=goroachuser DB_DATABASE=testdb DB_PASSWORD=<Password> DEFAULT_PAGE_SIZE=20 -a goroachapp`

### Deploy the app using docker image created earlier
`eyk builds:create ${ECR_REGISTRY}:${ECR_TAG} -a goroachapp --procfile='web: /main'`

### Test the app
run `eyk apps:info --app=goroachapp` and note the URL. Then, access the https://URL/quote in the browser. The app will respond a json data generated using the quotes from the database.

example: https://goroachapp.lab-two.ey-dedicated-internal.ey.io/quote/

eyk ps:console goroachapp-web-85d465fddb-vdjpf web goroachapp /bin/bash --debug=true