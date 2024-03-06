package goconf

import "github.com/iancoleman/strcase"

type caseType string

const (
	ScreamingSnake = caseType("ANY_KIND_OF_STRING")
	Snake          = caseType("any_kind_of_string")
	Kebab          = caseType("any-kind-of-string")
	ScreamingKebab = caseType("ANY-KIND-OF-STRING")
	Camel          = caseType("AnyKindOfString")
	LowerCamel     = caseType("anyKindOfString")
)

func changeCase(caseTypeVal caseType, envName string) string {
	m := map[caseType]func(string) string{
		ScreamingSnake: strcase.ToScreamingSnake,
		Snake:          strcase.ToSnake,
		Kebab:          strcase.ToKebab,
		ScreamingKebab: strcase.ToScreamingKebab,
		Camel:          strcase.ToCamel,
		LowerCamel:     strcase.ToLowerCamel,
	}
	f, ok := m[caseTypeVal]
	if !ok {
		return envName
	}
	return f(envName)
}

func removeQuotes(val string) string {
	startIdx := 0
	endIdx := len(val)
	if val[0] == '\'' || val[0] == '"' {
		startIdx = 1
	}
	if val[len(val)-1] == '\'' || val[len(val)-1] == '"' {
		if endIdx > 1 {
			endIdx = len(val) - 1
		}
	}
	return val[startIdx:endIdx]
}
