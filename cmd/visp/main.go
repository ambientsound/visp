package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/ambientsound/visp/log"
	"github.com/ambientsound/visp/prog"
	"github.com/ambientsound/visp/tokencache"
	"github.com/ambientsound/visp/version"
	"github.com/ambientsound/visp/widgets"
	"github.com/ambientsound/visp/xdg"
)

const (
	ConfigFileName = "visp.conf"
	TokenFileName  = "token.json"
)

const (
	ExitSuccess = iota
	ExitInternalError
	ExitPanic
)

func logAndStderr(line string) {
	log.Errorf(line)
	fmt.Fprintln(os.Stderr, line)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			logAndStderr("*********************************")
			logAndStderr("****** Visp has crashed!!! ******")
			logAndStderr("*********************************")
			logAndStderr("Please report this bug at the Github project and include the following information:")
			logAndStderr(fmt.Sprintf("Error message: %s", err))
			stacktrace := strings.Split(string(debug.Stack()), "\n")
			for _, line := range stacktrace {
				logAndStderr(line)
			}
			os.Exit(ExitPanic)
		}
	}()

	exitCode, err := run()
	if exitCode != ExitSuccess {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	os.Exit(exitCode)
}

//goland:noinspection GoUnhandledErrorResult
func run() (int, error) {
	log.Infof("%s %s starting up", version.Program, version.Version)
	log.Infof("This program was compiled on %s", version.BuildDate().String())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	visp := &prog.Visp{}

	ui, err := widgets.NewApplication(visp)
	if err != nil {
		return ExitInternalError, err
	}

	ui.Init()
	defer ui.Finish()
	go ui.Poll()

	visp.Termui = ui
	visp.Init()

	err = visp.SourceDefaultConfig()
	if err != nil {
		return ExitInternalError, fmt.Errorf("read default configuration: %s", err)
	}

	configDirs := xdg.ConfigDirectories()
	homeConfigDir := configDirs[len(configDirs)-1]
	err = os.MkdirAll(homeConfigDir, 0755)
	if err != nil {
		log.Errorf("Unable to create configuration directory: %s", err)
	}

	// Source configuration files from all XDG standard directories.
	for _, dir := range xdg.ConfigDirectories() {
		configFile := filepath.Join(dir, ConfigFileName)

		err = visp.SourceConfigFile(configFile)

		if errors.Is(err, os.ErrNotExist) {
			log.Debugf("Ignoring non-existing configuration file %s", configFile)
		} else if err != nil {
			log.Errorf("Error in configuration file %s: %s", configFile, err)
		}
	}

	// In case a token has been cached on disk, restore it to memory.
	tokenFile := filepath.Join(configDirs[len(configDirs)-1], TokenFileName)
	visp.Tokencache = tokencache.New(tokenFile)
	token, err := visp.Tokencache.Read()

	if err != nil {
		log.Debugf("Unable to read cached Spotify token: %s", err)
		token = nil
	}

	if token != nil {
		err = visp.Authenticate(token)
		if err != nil {
			log.Errorf(err.Error())
		}
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-signals
		log.Infof("Received signal %s, exiting.", s)
		cancel()
	}()

	log.Infof("Ready.")

	err = visp.Main(ctx)
	if err != nil {
		return ExitInternalError, err
	}

	return ExitSuccess, nil
}
