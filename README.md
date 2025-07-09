# cloudtrail-enrichment-api-golang

Security monitoring REST API that performs IP geolocation enrichment based on AWS CloudTrail logs.

Although it might seem like an anti-pattern, it is not in this context. There is a clear separation of concerns in data persistence. However, some potential failure risks may still exist. The architecture aims to leverage the strengths of MongoDB to store and process CloudTrail logs.


<p align="center">
  <img src="diagram.png" width="500"/>
</p>


## Run application

    docker-compose down -v --rmi all
    docker-compose up --build -d

## API Endpoints


The available API endpoints are outlined below.
<br >

<details><summary><code> Health GET /v1/health   </code></summary>

## 

Health Endpoint (pingpong)


Request


```
curl -XGET 'localhost:9090/v1/health'
```

Success Response:

 - Status Code: 200

</summary></details>

-----------------------------------------------------------

<details><summary><code> Signup POST /v1/signup  </code></summary>

## 
This endpoint is used to register a new user using an email and password for authentication.


Request

Body:

```json
{
  "email": "usuario@example.com",
  "password": "secret",
  "role": "user"
}
```

Success Response:

 - Status Code: 201

 - Body:

```json
{
  "error": false,
  "message": "1 registros procesados y enriquecidos exitosamente",
  "data": []
}
```


- Usage

```
    curl -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "email": "usuario@example.com",
        "password": "secret"
    }' \
    http://localhost:9090/v1/signup
```
</summary></details>

-----------------------------------------------------------

<details><summary><code> Login POST /v1/login  </code></summary>

## 
Authenticates a user by email and password, returning a JWT for accessing protected endpoints.


Request

Body:

```json
{
  "email": "usuario@example.com",
  "password": "secret"
}
```

Success Response:

 - Status Code: 201

 - Body:

```json
{
  "error": false,
  "message": "Autenticación exitosa",
  "data": {
    "email": "usuario@example.com",
    "expiry": "2025-07-09T21:24:40.863560028Z",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...."
  }
}
```


- Usage

```
curl -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "email": "usuario@example.com",
        "password": "secret"
        }' \
    http://localhost:9090/v1/login
```
</summary></details>

-----------------------------------------------------------
<details><summary><code> Log ingestion POST /v1/enrichment  </code></summary>

## 
This endpoint is used for ingesting logs from CloudTrail for geolocation enrichment based on event IP addresses.


Request

Body:

```json
{
    "Records":[
        {
            "eventVersion":"1.0",
            "userIdentity":{
                "type":"IAMUser",
                "principalId":"EX_PRINCIPAL_ID",
                "arn":"arn:aws:iam::123456789012:user/Alice",
                "accessKeyId":"EXAMPLE_KEY_ID",
                "accountId":"123456789012",
                "userName":"Alice"
            },
            "eventTime":"2014-03-06T21:22:54Z",
            "eventSource":"ec2.amazonaws.com",
            "eventName":"StartInstances",
            "awsRegion":"us-east-2",
            "sourceIPAddress":"205.251.233.176",
            "userAgent":"ec2-api-tools 1.6.12.2",
            "requestParameters":{
                "instancesSet":{
                "items":[
                    {
                        "instanceId":"i-ebeaf9e2"
                    }
                ]
                }
            },
            "responseElements":{
                "instancesSet":{
                "items":[
                    {
                        "instanceId":"i-ebeaf9e2",
                        "currentState":{
                            "code":0,
                            "name":"pending"
                        },
                        "previousState":{
                            "code":80,
                            "name":"stopped"
                        }
                    }
                ]
                }
            }
        }
    ]
    }
```

Success Response:

 - Status Code: 201

 - Body:

```json
{
  "error": false,
  "message": "1 registros procesados y enriquecidos exitosamente",
  "data": []
}
```


- Usage

```
curl -X POST \
    -H "Content-Type: application/json" \
    -d @cloudtrail_sample.json \
    -w "%{http_code}\n" \
    http://localhost:9090/v1/enrichment | jq

```
</summary></details>

-----------------------------------------------------------

<details><summary><code> Retrieve last 10 records GET /v1/enrichment </code></summary>

## 
This endpoint is used to retrieve the last 10 enriched CloudTrail events stored in the database.

Request

```
/v1/enrichment
```

Success Response:

 - Status Code: 200


- Usage

```
    curl http://localhost:9090/v1/enrichment
    curl http://localhost:9090/v1/enrichment | jq '.data.events | length'
```
</summary></details>

-----------------------------------------------------------

## Database Querys

        db.enriched_events.find()
        
        db.enriched_events.find(
        { "sourceIPAddress": "205.251.233.176" },
        { "enrichment": 1, "_id": 0 }
        )


## Secure coding practices

✅ Authentication and password management (JWT-based authentication middleware)

✅ Input validation _(Under construction)_

✅ Error Handling and Logging

✅ Cryptographic practices (JWT token implementation & design)

✅ Communication security (SSL/TLS - CORS - security headers)

✅ Database security _(Under construction)_


## Secure Deployment Practices 

✅  No hardcoded credentials in docker-compose.yml

✅  Container security using Bitnami images  _(Under construction)_

✅  



## References

[OWASP Secure Coding Practices - Checklist](https://owasp.org/www-project-secure-coding-practices-quick-reference-guide/stable-en/02-checklist/05-checklist)








