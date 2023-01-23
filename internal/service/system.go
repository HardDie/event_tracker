package service

import (
	"errors"
	"io"
	"os"

	"github.com/HardDie/event_tracker/internal/logger"
)

var (
	swaggerCache []byte
)

type ISystem interface {
	GetSwagger() ([]byte, error)
}

type System struct {
}

func NewSystem() *System {
	return &System{}
}

func (s *System) GetSwagger() ([]byte, error) {
	// Check cache
	if swaggerCache != nil {
		return swaggerCache, nil
	}

	// Open file
	file, err := os.Open("swagger.yaml")
	if err != nil {
		logger.Error.Println("error opening swagger.yaml file:", err.Error())
		return nil, errors.New("can't find swagger.yaml file")
	}
	defer file.Close()

	// Read data from file
	data, err := io.ReadAll(file)
	if err != nil {
		logger.Error.Println("error reading swagger.yaml file:", err.Error())
		return nil, errors.New("error reading swagger.yaml file")
	}

	// Save cache and return result
	swaggerCache = data
	return data, nil
}
