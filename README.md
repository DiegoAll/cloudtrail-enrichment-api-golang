# cloudtrail-enrichment-api-golang

Security monitoring REST API that performs IP geolocation enrichment based on AWS CloudTrail logs.


<img src="diagram.png" align="center"/>


## Run application

    docker-compose down -v --rmi all
    docker-compose up --build -d

## Security by Design (Secure Design Principles)

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




## Secure Deployment Practices 

- docker-compose.yml remove environment variables (.env it's not neccesary)
- Container security (Bitnami images, )
- Kubernetes security (Security context, )


## Software Engineering

No es un  antipattern. Una separación clara de responsabilidades en la persistencia de datos, lo cual es una buena práctica de diseño. Se  aprovechan las fortalezas de PostgreSQL para datos relacionales y de MongoDB para datos de documentos.


HAY UN ARCHIVO DE COPIA VALIDAR EL FALLO EN LAS FIRMAS DEL SERVICIO !!!!!

makefile

BACKUP en bash

lea en el folder ultima version
rm -rf db-data
zip 

- Commando para obtener listado de vscode abierto y generar una nueva copia
- discriminar por tesis y relacionados cursos portfolio


> Disclaimer

Aunque pueda parecer un antipatron, no lo es en este contexto, sin embargo existen algunos riesgos. Hay una clara separación de responsabilidades. Aprovechar las ventajas de mongo para logs de cloudtrail.





