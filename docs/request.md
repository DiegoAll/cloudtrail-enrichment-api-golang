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


