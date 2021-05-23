# needys-api-need
API micro-service to manage Need objects

## curl GET request example
curl http://localhost:8010/load

## curl PUT request example - 1
curl -X PUT -H "Content-Type: application/json" -d '[{"id": 1, "firstname","albert": "lastname","einstein"}]' http://localhost:8010/register

## curl PUT request example - 2
curl -X PUT -H "Content-Type: application/json" -d '[{"id": 1, "firstname","albert": "lastname","einstein"}, {"id": 2, "firstname","isaac": "lastname","newton"}, {"id": 3, "firstname","marie": "lastname","curie"}]' http://localhost:8010/register

## curl POST request example
curl -X POST -H "Content-Type: application/json" -d '{"firstname","aurélien": "lastname","barrault"}' http://localhost:8010/add

## curl DELETE request example
curl -X DELETE -H "Content-Type: application/json" -d '{"id": 4}' http://localhost:8010/delete
