#!/bin/bash

if [ -z "$1" ]; then
    echo "Error: The first argument is null or empty."
    exit 1 # Exit the script with an error code
fi

EL_PASSWORD=$1

# Start below commands only after elasticsearch is totally started
# ERROR: Failed to determine the health of the cluster., with exit code 69
# Check docker logs es01

printf "$1\n$1" | docker exec -i es01 /usr/share/elasticsearch/bin/elasticsearch-reset-password -b -i -u elastic
#docker exec -it es01 /usr/share/elasticsearch/bin/elasticsearch-reset-password -s -u elastic
echo "elasticsearch-reset-password done"

docker cp es01:/usr/share/elasticsearch/config/certs/http_ca.crt .
sleep 1
echo "ca cert copied"

echo "kibana token generating"
echo
EL_TOKEN=$(docker exec -it es01 /usr/share/elasticsearch/bin/elasticsearch-create-enrollment-token -f -s kibana)
echo "$EL_TOKEN"
echo

echo "kib01 starting..."
docker run -d --name kib01 --net elastic -p 5601:5601 docker.elastic.co/kibana/kibana:9.2.4
sleep 10
echo "kib01 started"

docker logs kib01