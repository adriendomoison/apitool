package apitool

import (
	"errors"
	"log"
	"os"
)

type Environment int

const (
	Dev Environment = iota
	Test
	Prod
)

var envString = [...]string{"DEV", "TEST", "PROD"}

func (e Environment) String() string {
	return envString[e]
}

type Context struct {
	AppEnvironment   Environment
	AppName          string
	AppDomain        string
	ServiceProtocol  string
	ServiceName      string
	ServicePort      string
	ServiceSubDomain string
	ServiceUrl       string
}

func (ctx *Context) SetUpAsExternal() error {
	if ctx.AppName == "" {
		return errors.New("missing environment variable APP_NAME")
	}
	if ctx.AppDomain == "" {
		return errors.New("missing environment variable APP_DOMAIN")
	}
	if ctx.ServicePort == "" {
		return errors.New("missing environment variable SERVICE_PORT")
	}
	if ctx.ServiceSubDomain == "" {
		return errors.New("missing environment variable SERVICE_SUB_DOMAIN")
	}
	if ctx.ServiceName == "" {
		return errors.New("missing environment variable SERVICE_NAME")
	}
	if err := configEnv(ctx); err != nil {
		return err
	}
	ctx.ServiceUrl = ctx.ServiceProtocol + "://" + ctx.ServiceSubDomain + "." + ctx.AppDomain + ":" + ctx.ServicePort
	return nil
}

func (ctx *Context) SetUpAsInternal() error {
	if ctx.ServiceName == "" {
		return errors.New("missing environment variable SERVICE_NAME")
	}
	if ctx.ServicePort == "" {
		return errors.New("missing environment variable SERVICE_PORT")
	}
	if err := configEnv(ctx); err != nil {
		return err
	}
	ctx.ServiceUrl = ctx.ServiceProtocol + "://" + ctx.ServiceName + ":" + ctx.ServicePort
	return nil
}

func configEnv(ctx *Context) error {
	if env := os.Getenv("ENVIRONMENT"); env == "" {
		return errors.New("missing environment variable ENVIRONMENT")
	} else if env == "DEV" {
		ctx.AppEnvironment = Dev
		ctx.ServiceProtocol = "http"
		log.Println("Dev environnement detected")
	} else if env == "TEST" {
		ctx.AppEnvironment = Test
		ctx.ServiceProtocol = "http"
		log.Println("Test environnement detected")
	} else if env == "PROD" {
		ctx.AppEnvironment = Prod
		ctx.ServiceProtocol = "https"
		log.Println("Prod environement detected")
	} else {
		return errors.New("unknown value " + env + " for ENVIRONMENT variable")
	}
	return nil
}
