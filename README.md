# cloudtrail-enrichment-api-golang

Security monitoring REST API that performs IP geolocation enrichment based on AWS CloudTrail logs.

> Disclaimer


## Run application

    docker-compose down -v --rmi all
    docker-compose up --build -d


## Secure Coding Practices

> Disclaimer

- JWT Based authentication using middleware:
 (token type=HS256, custom claims=role, iss, sub,aud, exp,nbf,iat)
- ✅ j.Config.AuthConfig.TokenDuration IMPORTANTE
- Logs
- Secrets Management
- HTTPS Transport Cipher
- HTTP security headers, CORS
- Strong cryptography 4 JWT, Hashing password etc
- CSRF Token Header <> X-CSRF-Token
- Authorization Role Based (Under construction)

## Secure Deployment Practices 

- docker-compose.yml remove environment variables (.env it's not neccesary)
- Container security (Bitnami images, )
- Kubernetes security (Security context, )

## Security by Design

- TM
- Auditory fields
- Token Design:
- Config (scaffold_config) ✅ Es ideal cuando no estás corriendo dentro de Docker. [Componente config para propagar las variables]
- .env para Docker y produccion

Parametros en texto plano, es mejor estandarizar y elegir uno.
- Podria servir para emular un test unitario del componente config.
- Redundante (escoger uno) config.go se presta para los dos.


UUID public API (mas seguro)
Si tienes una arquitectura monolítica y no estás preocupado por seguridad a ese nivel.
Pensar que siempre sera publica
Token con id


## Software Engineering

No es un  antipattern. Una separación clara de responsabilidades en la persistencia de datos, lo cual es una buena práctica de diseño. Se  aprovechan las fortalezas de PostgreSQL para datos relacionales y de MongoDB para datos de documentos.

## Otros


Tener una DB en la nube es super caro

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


El mensaje InvalidNamespace) Invalid namespace specified 'mydatabase.' proviene directamente del driver de MongoDB cuando intenta realizar una operación. Un "namespace" en MongoDB es la combinación de database.collection (ej., mydatabase.mycollection). El error explícitamente dice mydatabase., indicando que el problema está en el nombre de la base de datos.

    docker build --tag rest_app .