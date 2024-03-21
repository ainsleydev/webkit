package main

import (
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

type Command struct {
	Exec string
	Dir  string
}

type WebKit struct {
}

var commands = []Command{
	{
		Exec: "npx create-strapi-app@latest cms --no-run --use-npm --dbclient=sqlite --template=./templates/strapi",
		Dir:  "./",
	},
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	for _, cmd := range commands {
		a := strings.Split(cmd.Exec, " ")
		if err := StdExec(exec.Command(a[0], a[1:]...), cmd.Dir); err != nil {
			logger.Error(err.Error())
		}
	}
}

func StdExec(cmd *exec.Cmd, dir string) error {
	cmd.Dir = dir

	cmd.Stdout = io.MultiWriter(os.Stdout)
	cmd.Stderr = io.MultiWriter(os.Stderr)

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
