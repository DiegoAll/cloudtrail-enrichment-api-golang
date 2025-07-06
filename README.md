# cloudtrail-enrichment-api-golang
Security monitoring REST API that performs IP geolocation enrichment based on AWS CloudTrail logs.



# Scopes:

    export SCOPE=local
    SCOPE=local go run cmd/api/main.go


# Execute application

    docker-compose down -v --rmi all
    docker-compose up --build -d

    docker exec -it enrichment_api_db psql -U postgres -d booksdb
    docker exec -it enrichment_api_db psql -U postgres -d booksdb -c "\dt"
    docker exec -it enrichment_api_db psql -U postgres -d booksdb -c "SELECT * FROM users;"

    docker exec -it enrichment_api_db psql -U postgres -d booksdb -c "SELECT * FROM tokens;"

    docker build --tag rest_app .



## Secure Coding Practices

> Disclaimer

- JWT Based authentication using middleware (token type, custom claims)
- Logs
- Secrets Management
- HTTPS Transport Cipher
- HTTP security headers, CORS
- Strong cryptography 4 JWT, Hashing password etc
- CSRF Token Header <> X-CSRF-Token

## Secure Deployment Practices 

- Container security (Bitnami images, )
- Kubernetes security (Security context, )

## Security by Design

- TM


Tener una DB en la nube es super caro 

GCLOUD TOMBAS CIBERSFISICAS CAGADAS COLOMBIA
BUSCAR AWS

HOMOLOGAR ESTRUCTURA CONNECTION STRING MONGO CON PG. (Global var)


ALgunas variables no se pueden agregar de forma parametrica, si se puede hacer una sustitucion directa si.

De lo contrario son parametros para un generador.