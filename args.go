package goconf

type envSource string

const (
	OSEnv   = envSource("OSEnv")
	CMDArgs = envSource("CMDArgs")
)

var defaultSourcePrecedence = []envSource{
	OSEnv,
	CMDArgs,
}

type ExtractorOptions struct {
	// The path of the env file
	// i.e: ./env/.env.prod
	EnvFile string

	// If it is true then if the env file does not exist
	// it will be ignored
	OmitEnvFileIfNotExist bool

	// If it is true then all those fields that do not
	// contain the "env" tag will be omitted from the process
	OmitNotTagged bool

	// If the field has no tag and OmitNotTagged is false then
	// the field name is converted to this case type and later
	// is searched in the os env.
	// The default value is ScreamingSnake. i.e: "ANY_KIND_OF_STRING"
	EnvNameCaseType caseType

	// The default value is Snake. i.e: "any_kind_of_string"
	CMDArgsNameCaseType caseType

	// EnvSourcePrecedence is a slice of envSource that represents the order in which different
	// sources will be checked to retrieve values for environment variables.
	// The sources to be checked include OS environment variables and command-line arguments.
	EnvSourcePrecedence []envSource
}

func (e *ExtractorOptions) mergeWithDefault() {
	if e.EnvNameCaseType == "" {
		e.EnvNameCaseType = defaultEnvCase
	}
	if e.CMDArgsNameCaseType == "" {
		e.CMDArgsNameCaseType = defaultCMDArgCase
	}
	if e.EnvSourcePrecedence == nil || len(e.EnvSourcePrecedence) == 0 {
		e.EnvSourcePrecedence = defaultSourcePrecedence
	}
}

type ExtractorArgs struct {
	// Allow to configure the env extraction
	Options ExtractorOptions

	// It must be an even array of elements.
	// For each tuple:
	//  - The first element will be a pointer to the object in which the configuration will be saved.
	//  - The second element will be the prefix for this configuration
	Configs []interface{}
}
