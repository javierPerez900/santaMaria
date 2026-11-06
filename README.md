# API serverless en Go

API REST de la tienda online [Santa María](https://github.com/javierPerez900/santa-maria-ecommerce). Una única función **AWS Lambda** escrita en Go atiende todas las rutas detrás de **API Gateway**, con **Amazon Cognito** para la autenticación y **Aurora MySQL** como base de datos.

Es la parte de lógica de negocio del proyecto:

| Repositorio | Contenido |
|---|---|
| [santa-maria-ecommerce](https://github.com/javierPerez900/santa-maria-ecommerce) | Frontend en Next.js |
| **santaMaria** (este) | API REST |
| [santaMariaUser](https://github.com/javierPerez900/santaMariaUser) | Lambda que registra a los usuarios nuevos en la base de datos |
