# cloudtrail-enrichment-api-golang
Security monitoring REST API that performs IP geolocation enrichment based on AWS CloudTrail logs.



UUID public API (mas seguro)
Si tienes una arquitectura monolítica y no estás preocupado por seguridad a ese nivel.
Pensar que siempre sera publica

Token con id


# Scopes management

    export SCOPE=local
    SCOPE=local go run cmd/api/main.go


# Run application

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
- Authorization Role Based (Under construction)

## Secure Deployment Practices 

- Container security (Bitnami images, )
- Kubernetes security (Security context, )

## Security by Design

- TM
- Auditory fields


## Software Engineering

No es un  antipattern. Una separación clara de responsabilidades en la persistencia de datos, lo cual es una buena práctica de diseño. Se  aprovechan las fortalezas de PostgreSQL para datos relacionales y de MongoDB para datos de documentos.

## Otros


Tener una DB en la nube es super caro 

GCLOUD TOMBAS CIBERSFISICAS CAGADAS COLOMBIA
BUSCAR AWS

HOMOLOGAR ESTRUCTURA CONNECTION STRING MONGO CON PG. (Global var)


ALgunas variables no se pueden agregar de forma parametrica, si se puede hacer una sustitucion directa si.

De lo contrario son parametros para un generador.


	// Crear el cliente de MongoDB usando la URI y el timeout de la configuración
	mongoClient, err := mongo.NewMongoClient(mongoURI, config.MongoDBConfig.DBTimeout)
	if err != nil {
		log.Fatal("Error al conectar a MongoDB:", err)
		logger.ErrorLog.Fatalf("Error al conectar a MongoDB: %v", err)
	}
	defer func() {
		// Desconectar el cliente de MongoDB al finalizar
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			logger.ErrorLog.Printf("Error al cerrar la conexión a MongoDB: %v", err)
		}
	}()


Capa de composicion main.go
Definir y ensamblar las dependencias.


HAY UN ARCHIVO DE COPIA VALIDAR EL FALLO EN LAS FIRMAS DEL SERVICIO !!!!!

services/enrichment_service.go.txt

makefile

BACKUP en bash

lea en el folder ultima version
rm -rf db-data
zip 


Commando para obtener listado de vscode abierto y generar una nueva cpia

discriminar por
tesis y relacioonados
cursos
portfolio