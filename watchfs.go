package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"time"

	"code.google.com/p/go.exp/fsnotify"
)

var (
	timeout = flag.Int64("t", 500, "period to wait after the file change (to wait till the changes settle)")
	expr    = flag.String("f", "^[^\\.].*", "regular expression for file names to monitor for changes, default ignores hidden files. To filter \"*.go\" files use -f=\".*\\.go$\" to filter *.go files")

	fileExp *regexp.Regexp
	cmd     []string
)

func init() {
	//
	flag.Parse()
	cmd = flag.Args()
	fileExp = regexp.MustCompile(*expr)
}

func main() {
	if len(cmd) == 0 {
		flag.Usage()
		return
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			select {
			case event := <-watcher.Event:
				if !fileExp.MatchString(event.Name) {
					continue
				}
				for {
					select {
					case <-watcher.Event:
						continue
					case <-time.After(time.Duration(*timeout) * time.Millisecond):
					}
					runCommand()
					break
				}
			}
		}
	}()
	if err = watcher.Watch("."); err != nil {
		log.Fatal(err)
	}
	runCommand()
	select {}
}

func runCommand() {
	cmd := exec.Command(cmd[0], cmd[1:]...)
	stdout, err1 := cmd.StdoutPipe()
	stderr, err2 := cmd.StderrPipe()
	if err1 != nil || err2 != nil {
		log.Fatal(err1, err2)
	}
	err := cmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}
	go func() { io.Copy(os.Stdout, stdout) }()
	go func() { io.Copy(os.Stderr, stderr) }()
	err = cmd.Wait()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}
