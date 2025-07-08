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
        -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6InVzdWFyaW9AZXhhbXBsZS5jb20iLCJ0b2tlbiI6IiIsInRva2VuX2hhc2giOiIiLCJleHBpcnkiOiIyMDI1LTA3LTA5VDE3OjU4OjA1LjAxMDY4NjgxMVoiLCJyb2xlIjoidXNlciIsImlzcyI6Imczbm90eXBlIiwic3ViIjoidXN1YXJpb0BleGFtcGxlLmNvbSIsImF1ZCI6WyJtaXMtdXN1YXJpb3MiXSwiZXhwIjoxNzUyMDgzODg1LCJuYmYiOjE3NTE5OTc0ODUsImlhdCI6MTc1MTk5NzQ4NSwiY3JlYXRlZF9hdCI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIiwidXBkYXRlZF9hdCI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIn0.kzLoR4uIzgtGorZjwg3RmdHsc0YC7RtMYeS808NkyUg" \
        -d @cloudtrail_sample.json \
        -w "\nHTTP Code: %{http_code}\n" \
        http://localhost:9090/v1/enrichment | jq


## Enrichment (Get)

    curl -X GET \
        -H "Content-Type: application/json" \
        http://localhost:9090/v1/enrichment

    curl -X GET \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6InVzdWFyaW9AZXhhbXBsZS5jb20iLCJ0b2tlbiI6IiIsInRva2VuX2hhc2giOiIiLCJleHBpcnkiOiIyMDI1LTA3LTA5VDE3OjU4OjA1LjAxMDY4NjgxMVoiLCJyb2xlIjoidXNlciIsImlzcyI6Imczbm90eXBlIiwic3ViIjoidXN1YXJpb0BleGFtcGxlLmNvbSIsImF1ZCI6WyJtaXMtdXN1YXJpb3MiXSwiZXhwIjoxNzUyMDgzODg1LCJuYmYiOjE3NTE5OTc0ODUsImlhdCI6MTc1MTk5NzQ4NSwiY3JlYXRlZF9hdCI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIiwidXBkYXRlZF9hdCI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIn0.kzLoR4uIzgtGorZjwg3RmdHsc0YC7RtMYeS808NkyUg" \
        -w "\nHTTP Code: %{http_code}\n" \
        http://localhost:9090/v1/enrichment | jq





