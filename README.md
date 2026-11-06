# API serverless en Go

API REST de la [tienda online](https://github.com/javierPerez900/santa-maria-ecommerce). Una única función **AWS Lambda** escrita en Go atiende todas las rutas detrás de **API Gateway**, con **Amazon Cognito** para la autenticación y **Aurora MySQL** como base de datos.

> Esta parte incluye la lógica de negocio del proyecto. El proyecto completo se compone de tres repositorios:

| Repositorio | Contenido |
|---|---|
| [santa-maria-ecommerce](https://github.com/javierPerez900/santa-maria-ecommerce) | Frontend en Next.js |
| **santaMaria** (este) | API REST en Go (Lambda + API Gateway) |
| [santaMariaUser](https://github.com/javierPerez900/santaMariaUser) | Lambda en Go que guarda en la base de datos a cada usuario que confirma su cuenta en Cognito |
