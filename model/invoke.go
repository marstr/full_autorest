package model

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

// InvokeAutoRest executes AutoRest on the local machine.
func InvokeAutoRest(ctx context.Context, language AutoRestLanguage, inputFiles []string, options AutoRestOptions) error {
	invoker := exec.CommandContext(ctx, "autorest")

	invoker.Args = append(invoker.Args, string(language))

	if useVal, hasUseVal := options.Use(); hasUseVal {
		invoker.Args = append(invoker.Args, fmt.Sprintf("--use=%q", useVal))
	}
	if tagVal, hasTagVal := options.Tag(); hasTagVal {
		invoker.Args = append(invoker.Args, fmt.Sprintf("--tag=%q", tagVal))
	}

	outputFolder, _ := options.OutputFolder()
	invoker.Args = append(invoker.Args, fmt.Sprintf("--output-folder=%q", outputFolder))

	invoker.Args = append(invoker.Args, inputFiles...)

	invoker.Stdout, _ = options.Stdout()
	invoker.Stderr, _ = options.Stderr()

	return invoker.Run()
}

// AutoRestLanguage hosts the enumeration of the languages that
// have a generator implemented for AutoRest.
type AutoRestLanguage string

// The items declared here enumerate the well-known language generators.
// The intent is not necesarily to have a comprehensive list. New generators
// get created overtime, and any method taking a dependency on this type should
// be flexible.
const (
	AutoRestLanguageDotNet AutoRestLanguage = "--net"
	AutoRestLanguageGo     AutoRestLanguage = "--go"
	AutoRestLanguageJava   AutoRestLanguage = "--java"
	AutoRestLanguageRuby   AutoRestLanguage = "--ruby"
	AutoRestLanguagePHP    AutoRestLanguage = "--php"
	AutoRestLanguagePython AutoRestLanguage = "--python"
	AutoRestLanguageSwift  AutoRestLanguage = "--swift"
)

// AutoRestOptions encapsulates all of the flags that can be used to control AutoRest
type AutoRestOptions struct {
	outputFolder *string
	tag          *string
	use          *string
	stdout       io.Writer
	stderr       io.Writer
}

// OutputFolder fetches the location AutoRest should output the generated files.
// By default, the local temporary directory plus a unique identifier is used.
func (aro AutoRestOptions) OutputFolder() (string, bool) {
	if aro.outputFolder == nil {
		return path.Join(os.TempDir(), "generated"), false
	}
	return *aro.outputFolder, true
}

func (aro *AutoRestOptions) SetOutputFolder(dir string) *AutoRestOptions {
	aro.outputFolder = &dir
	return aro
}

// Tag fetches the value that should be used for the flag "--tag" when invoking AutoRest.
func (aro AutoRestOptions) Tag() (string, bool) {
	if aro.tag == nil {
		return "", false
	}
	return *aro.tag, true
}

// SetTag overwrites the value that would have previously been used for the flag "--tag" when
// invoking AutoRest.
func (aro *AutoRestOptions) SetTag(value string) *AutoRestOptions {
	aro.tag = &value
	return aro
}

// ClearTag restores the default value that will be used for the flag "--tag" when invoking
// AutoRest.
func (aro *AutoRestOptions) ClearTag() *AutoRestOptions {
	aro.tag = nil
	return aro
}

// Use fetches the value that should be used for the flag "--use" when invoking AutoRest.
// By default, this flag is not used.
func (aro AutoRestOptions) Use() (string, bool) {
	if aro.use == nil {
		return "", false
	}
	return *aro.use, true
}

// SetUse overwrites the value that would have previously been used for the flag "--use" when
// invoking AutoRest.
//
// This should be formatted as an npm package identifier. For example:
// 	 @microsoft.azure/autorest.go@~2
//   @microsoft.azure/autorest.go@2.1.87
func (aro *AutoRestOptions) SetUse(value string) *AutoRestOptions {
	aro.use = &value
	return aro
}

// ClearUse restores the default value that will be used for the flag "--use" when invoking
// AutoRest.
func (aro *AutoRestOptions) ClearUse() *AutoRestOptions {
	aro.use = nil
	return aro
}

// Stdout fetches the writer that will be used when AutoRest is invoked to communicate standard
// status.
func (aro AutoRestOptions) Stdout() (io.Writer, bool) {
	if aro.stdout == nil {
		return ioutil.Discard, false
	}
	return aro.stdout, true
}

// SetStdout overwrites the previous value that would be used
func (aro *AutoRestOptions) SetStdout(w io.Writer) *AutoRestOptions {
	aro.stdout = w
	return aro
}

// ClearStdout restores the default value that should be used for output from AutoRest.
func (aro *AutoRestOptions) ClearStdout() *AutoRestOptions {
	aro.stdout = nil
	return aro
}

// Stderr fetches the Writer that will be used by AutoRest to communicate error status.
// By default, AutoRest messages written to stderr will be discarded.
func (aro AutoRestOptions) Stderr() (io.Writer, bool) {
	if aro.stderr == nil {
		return ioutil.Discard, false
	}
	return aro.stderr, true
}

// SetStderr overwrites the writer that will be used by AutoRest to communicate error status.
func (aro *AutoRestOptions) SetStderr(w io.Writer) *AutoRestOptions {
	aro.stderr = w
	return aro
}

// ClearStderr restores the default functionality of Stderr.
// It is equivalent to calling `aro.SetStderr(nil)`.
func (aro *AutoRestOptions) ClearStderr() *AutoRestOptions {
	aro.stderr = nil
	return aro
}
