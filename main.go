package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/javier/santaMaria/awsgo"
	"github.com/javier/santaMaria/bd"
	"github.com/javier/santaMaria/handlers"
)

// Resultado del arranque, que se hace una vez por contenedor antes de recibir peticiones
var errInicio error

func main() {
	errInicio = inicializar()
	if errInicio != nil {
		// No se hace panic: el contenedor sigue vivo y la primera petición reintenta el arranque
		fmt.Println("Error al inicializar, se reintentará en la primera petición: " + errInicio.Error())
	}
	
	lambda.Start(EjecutoLambda)
}

func inicializar() error {
	if !ValidoParametros() {
		return errors.New("error en los parámetros: debe enviar 'SecretName' y 'UrlPrefix'")
	}

	if err := awsgo.InicializoAWS(); err != nil {
		return err
	}

	return bd.ReadSecret()
}

func EjecutoLambda(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {
	headersResp := map[string]string{
		"Content-Type": "application/json",
	}

	if errInicio != nil {
		// El arranque falló (por ejemplo, un fallo transitorio de Secrets Manager): se reintenta
		errInicio = inicializar()
		if errInicio != nil {
			fmt.Println("Error al inicializar: " + errInicio.Error())
			return &events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "Error interno del servidor",
				Headers:    headersResp,
			}, nil
		}
	}

	path := strings.Replace(request.RawPath, os.Getenv("UrlPrefix"), "", -1)
	method := request.RequestContext.HTTP.Method
	body := request.Body
	header := request.Headers

	status, message := handlers.Manejadores(path, method, body, header, request)

	return &events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       message,
		Headers:    headersResp,
	}, nil
}

func ValidoParametros() bool{
	_, traeParametro := os.LookupEnv("SecretName")
	if !traeParametro {
		return traeParametro
	}

	_, traeParametro = os.LookupEnv("UrlPrefix")
	if !traeParametro {
		return traeParametro
	}

	return traeParametro
}