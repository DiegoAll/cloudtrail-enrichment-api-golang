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



    curl -X POST \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6InVzdWFyaW9AZXhhbXBsZS5jb20iLCJ0b2tlbiI6IiIsInRva2VuX2hhc2giOiIiLCJleHBpcnkiOiIyMDI1LTA3LTA5VDAyOjExOjAwLjIwMzUwMTM3WiIsInJvbGUiOiJ1c2VyIiwiaXNzIjoiZzNub3R5cGUiLCJzdWIiOiJ1c3VhcmlvQGV4YW1wbGUuY29tIiwiYXVkIjpbIm1pcy11c3VhcmlvcyJdLCJleHAiOjE3NTIwMjcwNjAsIm5iZiI6MTc1MTk0MDY2MCwiaWF0IjoxNzUxOTQwNjYwLCJjcmVhdGVkX2F0IjoiMDAwMS0wMS0wMVQwMDowMDowMFoiLCJ1cGRhdGVkX2F0IjoiMDAwMS0wMS0wMVQwMDowMDowMFoifQ.g4OMZY1wh4TVhTdA90O6dDvkxmj9qSeXGcEukXvRWhc" \
        -d @cloudtrail_sample.json \
        -w "\nHTTP Code: %{http_code}\n" \
        http://localhost:9090/v1/enrichment/input | jq


## Enrichment (Get)

    curl -X GET \
        -H "Content-Type: application/json" \
        http://localhost:9090/v1/enrichment/get



curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6InVzdWFyaW9AZXhhbXBsZS5jb20iLCJ0b2tlbiI6IiIsInRva2VuX2hhc2giOiIiLCJleHBpcnkiOiIyMDI1LTA3LTA5VDAyOjE1OjE4LjE3MjEwMjU1MVoiLCJyb2xlIjoidXNlciIsImlzcyI6Imczbm90eXBlIiwic3ViIjoidXN1YXJpb0BleGFtcGxlLmNvbSIsImF1ZCI6WyJtaXMtdXN1YXJpb3MiXSwiZXhwIjoxNzUyMDI3MzE4LCJuYmYiOjE3NTE5NDA5MTgsImlhdCI6MTc1MTk0MDkxOCwiY3JlYXRlZF9hdCI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIiwidXBkYXRlZF9hdCI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIn0.yQYTKSwRiiIGDbNO9tTDyT__FKZLqcUz47O9sGZDEzE" \
    -d @cloudtrail_sample.json \
    -w "\nHTTP Code: %{http_code}\n" \
    http://localhost:9090/v1/enrichment/input | jq


