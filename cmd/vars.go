package cmd

import (
	"github.com/Lord-Y/synker/processing"
)

// List of vars that will be used by Synker
var (
	revision    string
	buildDate   string
	goVersion   string
	Version     = "0.0.1-dev"
	cmdValidate processing.Validate
)
