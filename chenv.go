package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/mgutz/ansi"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	// UsageFmt format of usage message
	UsageFmt = "" +
		"%s                    # list env profiles\n" +
		"%s <env>              # load <env> profile\n" +
		"%s save <env>         # save current .env as <env> profile"
	CurrentEnvPath = ".env"
)

var (
	// DisplayHelpFlag to display help
	DisplayHelpFlag = false
	// ChenvPath the path the chenv dir
	ChenvPath = ".chenv"
	// ProgramName name of go program
	ProgramName = "chenv"
	// Usage is the usage command
	Usage = fmt.Sprintf(
		UsageFmt,
		ProgramName,
		ProgramName,
		ProgramName,
	)
)

func init() {
	flag.BoolVar(
		&DisplayHelpFlag,
		"h",
		false,
		"`help` displays help message",
	)
}

// PrintHelpAndExit pring the help message
// and exit
func PrintHelpAndExit() {
	fmt.Fprintln(
		os.Stderr,
		Usage,
	)

	fmt.Fprintln(
		os.Stderr,
		"options:",
	)

	flag.PrintDefaults()
	os.Exit(-1)
}

func main() {
	flag.Parse()
	args := flag.Args()

	if DisplayHelpFlag {
		PrintHelpAndExit()
	}

	CheckChenvDir()

	argL := len(args)

	switch argL {
	case 0:
		PrintEnvs()
	case 1:
		env := args[0]
		Chenv(env)
	case 2:
		cmd := args[0]
		switch cmd {
		case "save":
			env := args[1]
			SaveEnv(env)
		default:
			PrintHelpAndExit()
		}
	default:
		PrintHelpAndExit()
	}
}

// SaveEnv ...
func SaveEnv(env string) {
	path := filepath.Join(
		ChenvPath,
		env+".env",
	)

	_, err := os.Stat(CurrentEnvPath)

	if os.IsNotExist(err) {
		PrintMissingCurrentEnv(err)
		ExitOn(err)
	}

	ErrExit(err)

	err = Copy(
		CurrentEnvPath,
		path,
	)

	ErrExit(err)

	msg := ansi.Color(
		fmt.Sprintf(
			"saved current .env to %s",
			path,
		),
		"green",
	)

	fmt.Println(
		msg,
	)

	msg = ansi.Color(
		fmt.Sprintf(
			"  %s saved!",
			env,
		),
		"cyan",
	)

	fmt.Println(
		msg,
	)

	fmt.Println()
}

// Chenv ...
func Chenv(env string) {
	path := filepath.Join(
		ChenvPath,
		env+".env",
	)

	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		PrintMissingEnv(env, path, err)
		ExitOn(err)
	}

	PrintError(err)
	ExitOn(err)

	err = Copy(path, ".env")
	ErrExit(err)

	loadLabel := ansi.Color(
		fmt.Sprintf(
			"loaded %s",
			path,
		),
		"green",
	)
	fmt.Println(loadLabel)

	usingLabel := ansi.Color(
		"using:",
		"green",
	)
	fmt.Println(usingLabel)

	envMsg := ansi.Color(
		"  * "+env,
		"cyan",
	)
	fmt.Println(envMsg)
}

// ErrExit ...
func ErrExit(err error) {
	PrintError(err)
	ExitOn(err)
}

// PrintError ...
func PrintError(err error) {
	if err != nil {
		fmt.Fprintln(
			os.Stderr,
			err.Error(),
		)
	}
}

// ExitOn ...
func ExitOn(err error) {
	if err != nil {
		os.Exit(-1)
	}
}

// CheckChenvDir ...
func CheckChenvDir() {
	_, err := os.Stat(ChenvPath)

	if os.IsNotExist(err) {
		fmt.Fprintf(
			os.Stderr,
			"%s is missing\n"+
				"creating %s\n"+
				"%s will be used to store .env files",
			ChenvPath,
			ChenvPath,
			ChenvPath,
		)
		err := os.Mkdir(ChenvPath, 0700)
		ErrExit(err)
	}
	ErrExit(err)
}

// EnvPaths ...
func EnvPaths() []string {
	filepaths, err := filepath.Glob(
		filepath.Join(
			ChenvPath,
			"*.env",
		),
	)

	ErrExit(err)

	return filepaths
}

// Copy ...
func Copy(src, dst string) error {
	s, err := os.Open(src)

	if err != nil {
		return err
	}
	defer s.Close()

	d, err := os.Create(dst)

	if err != nil {
		return err
	}

	_, err = io.Copy(d, s)

	if err != nil {
		d.Close()
		return err
	}

	return d.Close()
}

// PrintEnvs ...
func PrintEnvs() {
	currentEnvContent, currEnvErr := ioutil.ReadFile(
		CurrentEnvPath,
	)

	if currEnvErr != nil {
		PrintMissingCurrentEnv(currEnvErr)
		fmt.Fprintln(os.Stderr)
	}

	filepaths := EnvPaths()

	helpUsg := ansi.Color(
		"view usage with: chenv -h",
		"yellow",
	)
	fmt.Println(helpUsg)

	envLabel := ansi.Color(
		"available envs:",
		"green",
	)
	fmt.Println(envLabel)

	for _, path := range filepaths {
		name := filepath.Base(path)
		ext := filepath.Ext(path)

		name = name[:len(name)-len(ext)]

		pathContent, err := ioutil.ReadFile(
			path,
		)

		if err != nil || currEnvErr != nil {
			PrintEnv(name, false)
			continue
		}

		match := bytes.Equal(currentEnvContent, pathContent)
		PrintEnv(name, match)
	}

	fmt.Println()
}

// PrintEnv prints the given env name
func PrintEnv(name string, match bool) {
	sFmt := "    %s"
	if match {
		sFmt = "  * %s"
	}

	envLine := ansi.Color(
		fmt.Sprintf(
			sFmt,
			name,
		),
		"cyan",
	)

	fmt.Println(envLine)
}

// PrintMissingEnv ...
func PrintMissingEnv(name string, path string, err error) {
	if err != nil {
		errLine := ansi.Color(
			fmt.Sprintf(
				"%s",
				err,
			),
			"red",
		)

		fmt.Fprintln(
			os.Stderr,
			errLine,
		)
	}

	helpLine := ansi.Color(
		fmt.Sprintf(
			"missing %s in %s",
			name,
			path,
		),
		"yellow",
	)

	fmt.Fprintln(
		os.Stderr,
		helpLine,
	)
}

// PrintMissingCurrentEnv ...
func PrintMissingCurrentEnv(err error) {
	if err != nil {
		errLine := ansi.Color(
			fmt.Sprintf(
				"%s",
				err,
			),
			"red",
		)

		fmt.Fprintln(
			os.Stderr,
			errLine,
		)
	}

	helpLine := ansi.Color(
		fmt.Sprintf(
			"%s",
			"did you create a .env?",
		),
		"yellow",
	)

	fmt.Fprintln(
		os.Stderr,
		helpLine,
	)
}
