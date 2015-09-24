package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/mgutz/ansi"
	"path/filepath"
	"io"
)

var (
	// FlagSet ListEnvs
	FlagSet, SaveEnv, Help = ChenvFlagSet()
	ChenvPath = ".chenv"
)

// ChenvFlagSet ...
func ChenvFlagSet() (*flag.FlagSet, *bool, *bool) {
	f := flag.NewFlagSet(
		"chenv",
		flag.ContinueOnError,
	)

	sFlag := f.Bool(
		"save",
		false,
		"save the current env to the specified profile",
	)

	hFlag := f.Bool(
		"help",
		false,
		"display usage",
	)

	return f, sFlag, hFlag
}

func main() {
	args := os.Args[1:]

	err := FlagSet.Parse(args)
	PrintError(err)
	ExitOn(err)

	if *Help {
		fmt.Fprintf(
			os.Stderr,
			"%s\n%s\n",
			"Usage: [ENV] set the current .env to [ENV]",
			"Usage: -save [ENV] save the current .env to [ENV]",
		)

		FlagSet.PrintDefaults()
		return
	}

	CheckChenvDir()

	if FlagSet.NArg() < 1 {
		PrintEnvs()
		return
	}

	env := FlagSet.Arg(0)

	if *SaveEnv {
		StoreEnv(env)
		return
	}

	Chenv(env)
}

// StoreEnv ...
func StoreEnv(env string) {
	path := filepath.Join(
		ChenvPath,
		env + ".env",
	)

	_, err := os.Stat(".env")

	if os.IsNotExist(err) {
		fmt.Fprintf(
			os.Stderr,
			".env does not exist.\n%s\n",
			"Please create a current .env",
		)
		ExitOn(err)
	}

	PrintError(err)
	ExitOn(err)

	err = Copy(".env", path)

	PrintError(err)
	ExitOn(err)

	fmt.Println(
		ansi.Color(
			"Saved "+env,
			"cyan",
		),
	)
}

// Chenv ...
func Chenv(env string) {
	path := filepath.Join(
		ChenvPath,
		env + ".env",
	)

	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		fmt.Fprintf(
			os.Stderr,
			"%s does not exist.\n%s\n%s\n",
			env,
			"To create it run:",
			"\tchenv -save "+env,
		)
		ExitOn(err)
	}

	PrintError(err)
	ExitOn(err)

	err = Copy(path, ".env")
	PrintError(err)
	ExitOn(err)

	fmt.Println(
		ansi.Color(
			"Now using "+env,
			"cyan",
		),
	)
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
			"%s\n%s\n%s\n%s\n",
			"The .chenv directory does not exist.",
			"The .chenv directory is used to store env profiles.",
			"\tRun:",
			"\tmkdir .chenv",
		)
		ExitOn(err)
	}

	PrintError(err)
	ExitOn(err)
}

// EnvPaths ...
func EnvPaths() []string {
	filepaths, err := filepath.Glob(
		filepath.Join(
			ChenvPath,
			"*.env",
		),
	)

	PrintError(err)
	ExitOn(err)

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

	_, err = io.Copy(d, s);

	if err != nil {
		d.Close()
		return err
	}

	return d.Close()
}

// PrintEnvs ...
func PrintEnvs() {
	filepaths := EnvPaths()

	fmt.Printf(
		"%s\n%s\n\n",
		ansi.Color(
			"View usage with chenv --help",
			"yellow",
		),
		ansi.Color(
			"Listing Envs...",
			"green",
		),
	)

	for _, path := range filepaths {
		name := filepath.Base(path)
		ext := filepath.Ext(path)

		name = name[:len(name) - len(ext)]

		fmt.Printf(
			"\t%s\n",
			ansi.Color(
				"- "+name,
				"cyan",
			),
		)
	}

	fmt.Fprintln(
		os.Stderr,
	)
}
