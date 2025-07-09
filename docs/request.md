# Request

Quitar el semver

## Health

    curl -X GET localhost:9090/v1/health

## Signup

    curl -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "email": "user@example.com",
        "password": "secret"
    }' \
    http://localhost:9090/v1/signup

## Login

    curl -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "email": "user@example.com",
            "password": "secret"
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
        -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJlbWFpbCI6InVzZXJAZXhhbXBsZS5jb20iLCJ0b2tlbiI6IiIsInRva2VuX2hhc2giOiIiLCJleHBpcnkiOiIyMDI1LTA3LTEwVDA0OjE4OjQ4Ljg0OTUwMDg3N1oiLCJyb2xlIjoidXNlciIsImlzcyI6Imczbm90eXBlIiwic3ViIjoidXNlckBleGFtcGxlLmNvbSIsImF1ZCI6WyJtaXMtdXN1YXJpb3MiXSwiZXhwIjoxNzUyMTIxMTI4LCJuYmYiOjE3NTIwMzQ3MjgsImlhdCI6MTc1MjAzNDcyOCwiY3JlYXRlZF9hdCI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIiwidXBkYXRlZF9hdCI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIn0.VnLqVSRx4Qy_Dw43pRKPyhlFuWa5yKWZIrVv95WBrwo" \
        -d @cloudtrail_sample.json \
        -w "\nHTTP Code: %{http_code}\n" \
        http://localhost:9090/v1/enrichment | jq


## Enrichment (Get)

    curl -X GET \
        -H "Content-Type: application/json" \
        http://localhost:9090/v1/enrichment

    curl -X GET \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6InVzdWFyaW9AZXhhbXBsZS5jb20iLCJ0b2tlbiI6IiIsInRva2VuX2hhc2giOiIiLCJleHBpcnkiOiIyMDI1LTA3LTA5VDIxOjI0OjQwLjg2MzU2MDAyOFoiLCJyb2xlIjoidXNlciIsImlzcyI6Imczbm90eXBlIiwic3ViIjoidXN1YXJpb0BleGFtcGxlLmNvbSIsImF1ZCI6WyJtaXMtdXN1YXJpb3MiXSwiZXhwIjoxNzUyMDk2MjgwLCJuYmYiOjE3NTIwMDk4ODAsImlhdCI6MTc1MjAwOTg4MCwiY3JlYXRlZF9hdCI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIiwidXBkYXRlZF9hdCI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIn0.y0ZmklIsF_8V07-oGvCp3SJTIM5GyudXXB6kUSNdry4" \
        -w "\nHTTP Code: %{http_code}\n" \
        http://localhost:9090/v1/enrichment | jq





