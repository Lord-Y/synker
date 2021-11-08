package cmd

import (
	"github.com/Lord-Y/synker/api"
	"github.com/Lord-Y/synker/models"
	"github.com/Lord-Y/synker/validate"
)

// List of vars that will be used by Synker
var (
	revision     string
	buildDate    string
	goVersion    string
	Version      = "0.0.1-dev"
	configStruct models.Configuration
	cmdValidate  validate.Validate
	cmdAPI       api.API
)
