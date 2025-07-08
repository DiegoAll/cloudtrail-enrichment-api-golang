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
    http://localhost:9090/v1/signup

## Authenticate

    curl -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "email": "usuario@example.com",
            "password": "123123123"
            }' \
        http://localhost:9090/v1/login


## Enrichment (Input)

    curl -X POST \
        -H "Content-Type: application/json" \
        -d @cloudtrail_sample.json \
        -w "%{http_code}\n" \
        http://localhost:9090/v1/enrichment | jq


    curl -X POST \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6InVzdWFyaW9AZXhhbXBsZS5jb20iLCJ0b2tlbiI6IiIsInRva2VuX2hhc2giOiIiLCJleHBpcnkiOiIyMDI1LTA3LTA5VDE4OjI3OjA3LjczMTAzNzI2OFoiLCJyb2xlIjoidXNlciIsImlzcyI6Imczbm90eXBlIiwic3ViIjoidXN1YXJpb0BleGFtcGxlLmNvbSIsImF1ZCI6WyJtaXMtdXN1YXJpb3MiXSwiZXhwIjoxNzUyMDg1NjI3LCJuYmYiOjE3NTE5OTkyMjcsImlhdCI6MTc1MTk5OTIyNywiY3JlYXRlZF9hdCI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIiwidXBkYXRlZF9hdCI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIn0.vEmht_x5sJJDdOgwAPV6-qD2J4V-ceLOBlz33Yy0e8I" \
        -d @cloudtrail_sample.json \
        -w "\nHTTP Code: %{http_code}\n" \
        http://localhost:9090/v1/enrichment | jq


## Enrichment (Get)

    curl -X GET \
        -H "Content-Type: application/json" \
        http://localhost:9090/v1/enrichment

    curl -X GET \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6InVzdWFyaW9AZXhhbXBsZS5jb20iLCJ0b2tlbiI6IiIsInRva2VuX2hhc2giOiIiLCJleHBpcnkiOiIyMDI1LTA3LTA5VDE4OjI3OjA3LjczMTAzNzI2OFoiLCJyb2xlIjoidXNlciIsImlzcyI6Imczbm90eXBlIiwic3ViIjoidXN1YXJpb0BleGFtcGxlLmNvbSIsImF1ZCI6WyJtaXMtdXN1YXJpb3MiXSwiZXhwIjoxNzUyMDg1NjI3LCJuYmYiOjE3NTE5OTkyMjcsImlhdCI6MTc1MTk5OTIyNywiY3JlYXRlZF9hdCI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIiwidXBkYXRlZF9hdCI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIn0.vEmht_x5sJJDdOgwAPV6-qD2J4V-ceLOBlz33Yy0e8I" \
        -w "\nHTTP Code: %{http_code}\n" \
        http://localhost:9090/v1/enrichment | jq





