package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	go_amqp_lib "github.com/lan143/go-amqp-lib"
)

type Config struct {
	amqpConfig         go_amqp_lib.AmqpConfig
	webServerAddress   string
	healthCheckAddress string
	isDebug            bool
}

func (c *Config) Init() error {
	log.Println("config: init")

	_ = godotenv.Load()

	err := c.loadAmqpConfig()
	if err != nil {
		return err
	}

	err = c.loadWebServerAddress()
	if err != nil {
		return err
	}

	err = c.loadHealthCheckAddress()
	if err != nil {
		return err
	}

	err = c.loadIsDebug()
	if err != nil {
		return err
	}

	log.Println("config.init: successful")

	return nil
}

func (c *Config) GetAmqpConfig() go_amqp_lib.AmqpConfig {
	return c.amqpConfig
}

func (c *Config) GetWebServerAddress() string {
	return c.webServerAddress
}

func (c *Config) GetHealthCheckAddress() string {
	return c.healthCheckAddress
}

func (c *Config) IsDebug() bool {
	return c.isDebug
}

func (c *Config) loadAmqpConfig() error {
	c.amqpConfig = go_amqp_lib.AmqpConfig{
		Host:     os.Getenv("AMQP_HOST"),
		Port:     os.Getenv("AMQP_PORT"),
		Username: os.Getenv("AMQP_USER"),
		Password: os.Getenv("AMQP_PASSWORD"),
		VHost:    os.Getenv("AMQP_VHOST"),
	}

	return nil
}

func (c *Config) loadWebServerAddress() error {
	c.webServerAddress = os.Getenv("WEB_SERVER_ADDR")
	if len(c.healthCheckAddress) == 0 {
		c.webServerAddress = "0.0.0.0:80"
	}

	return nil
}

func (c *Config) loadHealthCheckAddress() error {
	c.healthCheckAddress = os.Getenv("HEALTH_CHECK_ADDR")
	if len(c.healthCheckAddress) == 0 {
		c.healthCheckAddress = "0.0.0.0:8080"
	}

	return nil
}

func (c *Config) loadIsDebug() error {
	isDebugVal := os.Getenv("DEBUG")

	if strings.Compare(isDebugVal, "true") == 0 || strings.Compare(isDebugVal, "1") == 0 {
		c.isDebug = true
	} else {
		c.isDebug = false
	}

	return nil
}

func NewConfig() *Config {
	return &Config{}
}
