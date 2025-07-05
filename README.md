# cloudtrail-enrichment-api-golang
Security monitoring REST API that performs IP geolocation enrichment based on AWS CloudTrail logs.


# Execute application

    docker-compose down -v --rmi all
    docker-compose up --build -d

    docker exec -it enrichment_api_db psql -U postgres -d booksdb
    docker exec -it enrichment_api_db psql -U postgres -d booksdb -c "\dt"
    docker exec -it enrichment_api_db psql -U postgres -d booksdb -c "SELECT * FROM users;"

    docker exec -it enrichment_api_db psql -U postgres -d booksdb -c "SELECT * FROM tokens;"

    docker build --tag rest_app .



# Security practices

> Disclaimer

- JWT Based authentication using middleware (token type, custom claims)
- Logs
- Secrets Management
- HTTPS Transport Cipher
- HTTP security headers, CORS
- Strong cryptography 4 JWT, Hashing password etc

# Security by Design

