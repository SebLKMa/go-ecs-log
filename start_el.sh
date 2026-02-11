echo "es01 starting..."
docker run -d --name es01 --net elastic -p 9200:9200 -it -m 1GB docker.elastic.co/elasticsearch/elasticsearch:9.2.4
sleep 8
echo "es01 started"

docker logs es01