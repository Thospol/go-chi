#!/bin/bash
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}|------------------------------------------------------------------------|${NC}"
echo -e "${GREEN}|========== Welcome to SAAA! Have a good day  :)  ==========|${NC}"
echo -e "${GREEN}|------------------------------------------------------------------------|${NC}"
sleep 0.5

echo
echo "### Please Input namespaces for deployment ###"
echo "1) dev"
echo "2) prod"
echo

while :; do
	echo -n "Input: "
	read INPUT_STRING

	case $INPUT_STRING in
	dev)
		echo -e "${GREEN}<SAAA>${NC} 1. set project"
		gcloud config set project saaa
		echo -e "${GREEN}<SAAA>${NC} 2. docker build"
		gcloud build -f Dockerfile.dev -t asia.gcr.io/saaa/saaa-api:develop .
        echo -e "${GREEN}<SAAA>${NC} 3. docker push"
		docker push asia.gcr.io/saaa/saaa-api:develop
        echo -e "${GREEN}<SAAA>${NC} 4. deploy to cloud run"
		gcloud run deploy saaa-api-dev --image asia.gcr.io/saaa/saaa-api:develop --region asia-southeast1 --platform managed --allow-unauthenticated
		break
		;;

	prod)
		echo -e "${GREEN}<SAAA>${NC} 1. set project"
		gcloud config set project saaa
		echo -e "${GREEN}<SAAA>${NC} 2. docker build"
		gcloud build -t asia.gcr.io/saaa/saaa-api:master .
        echo -e "${GREEN}<SAAA>${NC} 3. docker push"
		docker push asia.gcr.io/saaa/saaa-api:master
        echo -e "${GREEN}<SAAA>${NC} 4. deploy to cloud run"
		gcloud run deploy saaa-api --image asia.gcr.io/saaa/saaa-api:master --region asia-southeast1 --platform managed --allow-unauthenticated
		break
		break
		;;
	*)
		echo -e "${RED}Input invalid Please try again.${NC}"
		break
		;;
	esac
done
echo "Finish."
sleep 0.5
exit 0
