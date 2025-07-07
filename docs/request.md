# Request

Quitar el semver

## Health

    curl -X GET localhost:9090/v1/health

## Register

    curl -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "email": "usuario@example.com",
        "password": "123123123",
        "role": "user"
    }' \
    http://localhost:9090/v1/register

## Authenticate

    curl -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "email": "diego@example.com",
            "password": "123123123"
            }' \
        http://localhost:9090/v1/authenticate


curl -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "email": "diego@diego.com",
            "password": "123123123"
            }' \
        http://localhost:9090/v1/authenticate  NO DA MARAÑA TREVOR


curl -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "email": "usuario@example.com",
            "password": "123123123"
            }' \
        http://localhost:9090/v1/authenticate


## Enrichment (Input)

    curl -X POST \
        -H "Content-Type: application/json" \
        -d @cloudtrail_sample.json \
        -w "%{http_code}\n" \
        http://localhost:9090/v1/enrichment/input | jq


## Enrichment (Get)

    curl -X GET \
        -H "Content-Type: application/json" \
        http://localhost:9090/v1/enrichment/get


root@pho3nix:/home/diegoall/Projects/cloudtrail-enrichment-api-golang# curl -X POST     -H "Content-Type: application/json"     -d '{
        "email": "usuario@example.com",
        "password": "123123123",
        "role": "user"
    }'     http://localhost:9090/v1/register
{"error":true,"message":"error al verificar email: error al obtener usuario por email: pq: column \"uuid\" does not exist"}


Ya registra pero no autentica.


