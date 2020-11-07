package utils

import (
	"log"
	"strings"
)

//Validator - validate files extension
type Validator struct {
	logger *log.Logger
	config *Configuration
}

//CreateValidator - create extension validator
func CreateValidator(logger *log.Logger, config *Configuration) *Validator {
	validator := Validator{logger: logger, config: config}

	logger.Println("File validator created")

	return &validator
}

//ValidateScaledFileExtension - validate which image can be scaled
func (validator *Validator) ValidateScaledFileExtension(extension string) bool {
	for _, ext := range strings.Split(validator.config.ScaledImageRestoreExtension, ",") {
		if ext == extension {
			return true
		}
	}

	return false
}

//ValidateSavedFileExtension - validate which image servise can save
func (validator *Validator) ValidateSavedFileExtension(extension string) bool {
	for _, ext := range strings.Split(validator.config.FileSaveExtensionList, ",") {
		if ext == extension {
			return true
		}
	}

	return false
}
