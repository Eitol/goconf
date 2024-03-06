package goconf

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
	"strings"
	"unicode"
)

func loadAllEnvValFromEnvFile(args ExtractorArgs) error {
	if !pathExist(args.Options.EnvFile) {
		if args.Options.OmitEnvFileIfNotExist {
			return nil
		}
		return errors.New("the configuration file doesnt exist: " + args.Options.EnvFile)
	}
	err := godotenv.Load(args.Options.EnvFile)
	if err != nil {
		return err
	}
	return nil
}

func isValidCMDArgStartName(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

func cleanCMDArgName(arg string) string {
	if arg[0] == '-' {
		// case "-param=123" or "-active"
		arg = strings.TrimPrefix(arg, "-")
		// case "--param=123" or "--active"
		arg = strings.TrimPrefix(arg, "-")
	}
	return arg
}

func getEnvValFromCMDArgs(envName string, args []string) string {
	if args == nil || len(args) < 2 {
		return ""
	}
	for _, arg := range args {
		arg = cleanCMDArgName(arg)
		if !isArgStartValid(arg) {
			continue
		}
		// boolean flag case
		if arg == envName {
			return defaultCMDArgValue
		}
		value := checkEnvNameInArg(arg, envName)
		if value != "" {
			return value
		}
	}
	return ""
}

func isArgStartValid(arg string) bool {
	if len(arg) > 0 {
		r := rune(arg[0])
		return isValidCMDArgStartName(r)
	}
	return false
}

func checkEnvNameInArg(arg, envName string) string {
	// look for '=' in arg
	spArgs := strings.SplitN(arg, "=", 2)
	if spArgs[0] != envName {
		return ""
	}
	if len(spArgs) > 1 {
		return removeQuotes(spArgs[1])
	}
	return defaultCMDArgValue
}

func getEnvValFromOSEnv(prefix, envName string) string {
	val, exist := os.LookupEnv(prefix + envName)
	if !exist {
		return ""
	}
	return val
}

func getEnvValuesFromSources(prefix, envName string, opts ExtractorOptions) (string, envSource) {
	val := ""
	var finalSource envSource
	for _, source := range opts.EnvSourcePrecedence {
		switch source {
		case OSEnv:
			envNameKey := changeCase(opts.EnvNameCaseType, envName)
			val = getEnvValFromOSEnv(prefix, envNameKey)
		case CMDArgs:
			envNameKey := changeCase(opts.CMDArgsNameCaseType, envName)
			val = getEnvValFromCMDArgs(envNameKey, os.Args)
		}
		if val == "" {
			continue
		}
		finalSource = source
		break
	}
	return val, finalSource
}
