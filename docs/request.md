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

character varying(512)
text

Nadie almacena los tokens en la base de datos, por eso restfull. DB in memory.
Concepto de sesion.


FUNCIONANDO


root@pho3nix:/home/diegoall/Projects/cloudtrail-enrichment-api-golang# curl -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "email": "usuario@example.com",
        "password": "123123123",
        "role": "user"
    }' \
    http://localhost:9090/v1/register
{"error":false,"message":"Usuario registrado exitosamente","data":{"email":"usuario@example.com","role":"user","uuid":"a3621bf8-46ba-4261-bcbd-4ae193e5c0cd"}}root@pho3nixroot@pho3nix:/home/diegoall/Projects/cloudtrail-enrichment-api-golang# curl -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "email": "usuario@example.com",
            "password": "123123123"
            }' \
        http://localhost:9090/v1/authenticate
{"error":false,"message":"Autenticación exitosa","data":{"email":"usuario@example.com","expiry":"2025-07-08T18:19:15.862758168Z","role":"user","token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6InVzdWFyaW9AZXhhbXBsZS5jb20iLCJ0b2tlbiI6IiIsInRva2VuX2hhc2giOiIiLCJleHBpcnkiOiIyMDI1LTA3LTA4VDE4OjE5OjE1Ljg2Mjc1ODE2OFoiLCJyb2xlIjoidXNlciIsImlzcyI6Imczbm90eXBlIiwic3ViIjoidXN1YXJpb0BleGFtcGxlLmNvbSIsImF1ZCI6WyJtaXMtdXN1YXJpb3MiXSwiZXhwIjoxNzUxOTk4NzU1LCJuYmYiOjE3NTE5MTIzNTUsImlhdCI6MTc1MTkxMjM1NSwiY3JlYXRlZF9hdCI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIiwidXBkYXRlZF9hdCI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIn0.drpjdYUtcq4lHxRSRq9WbWESb_mFWrSs0wKVvDAYvbk","user_uuid":"a3621bf8-46ba-4261-bcbd-4ae193e5c0cd"}}


Tamaño del token frente a ataques.


curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6InVzdWFyaW9AZXhhbXBsZS5jb20iLCJ0b2tlbiI6IiIsInRva2VuX2hhc2giOiIiLCJleHBpcnkiOiIyMDI1LTA3LTA4VDE4OjE5OjE1Ljg2Mjc1ODE2OFoiLCJyb2xlIjoidXNlciIsImlzcyI6Imczbm90eXBlIiwic3ViIjoidXN1YXJpb0BleGFtcGxlLmNvbSIsImF1ZCI6WyJtaXMtdXN1YXJpb3MiXSwiZXhwIjoxNzUxOTk4NzU1LCJuYmYiOjE3NTE5MTIzNTUsImlhdCI6MTc1MTkxMjM1NSwiY3JlYXRlZF9hdCI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIiwidXBkYXRlZF9hdCI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIn0.drpjdYUtcq4lHxRSRq9WbWESb_mFWrSs0wKVvDAYvbk" \
    -d @cloudtrail_sample.json \
    -w "\nHTTP Code: %{http_code}\n" \
    http://localhost:9090/v1/enrichment/input | jq
