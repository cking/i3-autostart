package main

import (
	"errors"
	"github.com/keegancsmith/shell"
	"github.com/rkoesters/xdg/basedir"
	"github.com/rkoesters/xdg/desktop"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func main() {
	configDirs := append([]string{}, /*basedir.ConfigDirs,*/ basedir.ConfigHome)

	for _, d := range configDirs {
		autostartDir := path.Join(d, "autostart")
		println("checking", autostartDir, "for autostart files...")
		filepath.Walk(autostartDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				// skip directories
				if path == autostartDir {
					return nil
				}

				return filepath.SkipDir
			}

			r, err := os.Open(path)
			if err != nil {
				return err
			}

			d, err := desktop.New(r)
			if err != nil {
				return err
			}

			if d.Type != desktop.Application {
				return errors.New("unexpected type")
			}

			if d.TryExec != "" {
				if ext, err := exec.LookPath(d.TryExec); err != nil {
					return nil
				} else if _, err := os.Open(ext); err != nil {
					return nil
				}
			}

			println("found", d.Name, "at", path)
			cmd := shell.Commandf(d.Exec)
			cmd.Dir = d.Path
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Start(); err != nil {
				println("!!! ERROR: ", err.Error())
				return err
			}

			return nil
		})
	}
}
