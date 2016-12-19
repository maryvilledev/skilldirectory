echo "Testing..."
go test $(glide novendor)|| { echo "Tests failed" ; exit 1; }

echo "Buildling..."
go build

running=$(docker inspect -f {{.State.Running}} cassandra_container)
echo $running
if $running; then
  echo "Already Running"
else
  docker-compose build
  docker-compose up -d
  sleep 20
  
fi
docker exec -it cassandra_container bash usr/bin/cqlsh -f /data/skilldirectoryschema.cql

echo "Running Skill Directory..."
./skilldirectory
