# saaa-api

## How to install google cloud platform.
1. Download GCloud 
    - for mac -> https://cloud.google.com/sdk/docs/quickstart-macos
    - for window -> https://cloud.google.com/sdk/docs/quickstart-windows
2. gcloud auth application-default login
3. gcloud components install kubectl
4. gcloud config set project saaa
5. gcloud container clusters get-credentials saaa --zone asia-southeast1-a

## How to install mysql on docker.
1. docker pull mysql
2. docker run --name=saaa-local -e MYSQL_ROOT_PASSWORD=P@ssw0rd -e MYSQL_DATABASE=saaa -p 3306:3306 -d mysql

## Customize configuration
Google Cloud Platform [Reference](https://cloud.google.com/sdk/docs/quickstart).
Docker Mysql [Reference](https://hub.docker.com/_/mysql).