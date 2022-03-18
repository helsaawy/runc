package main

import (
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/opencontainers/runc/libcontainer"
	_ "github.com/opencontainers/runc/libcontainer/nsenter"
	"github.com/sirupsen/logrus"
)

const ddir = "/run/"

func init() {

	pre := time.Now()
	if len(os.Args) > 1 && os.Args[1] == "init" {
		// This is the golang entry point for runc init, executed
		// before main() but after libcontainer/nsenter's nsexec().
		runtime.GOMAXPROCS(1)
		runtime.LockOSThread()
		pst := time.Now()

		level, err := strconv.Atoi(os.Getenv("_LIBCONTAINER_LOGLEVEL"))
		if err != nil {
			panic(err)
		}

		logPipeFd, err := strconv.Atoi(os.Getenv("_LIBCONTAINER_LOGPIPE"))
		if err != nil {
			panic(err)
		}

		logrus.SetLevel(logrus.Level(level))
		logrus.SetOutput(os.NewFile(uintptr(logPipeFd), "logpipe"))
		logrus.SetFormatter(new(logrus.JSONFormatter))
		logrus.Debug("child process pre  ", pre.Format(time.RFC3339Nano))
		logrus.Debug("child process post ", pst.Format(time.RFC3339Nano))
		logrus.Debug("child process in init()")
		logrus.Debugf("child process env %q", os.Environ())

		if d, err := os.ReadFile("/proc/self/stat"); err == nil {
			logrus.Info("init self stat:\n" + string(d))
		}

		// if d, err := os.ReadFile("/proc/self/status"); err == nil {
		// 	logrus.Info("init self status:\n" + string(d))
		// }

		if d, err := os.ReadFile("/proc/self/sched"); err == nil {
			logrus.Info("init self sched:\n" + string(d))
		}

		if d, err := os.ReadFile("/proc/self/schedstat"); err == nil {
			logrus.Info("init self sched stat:\n" + string(d))
		}

		if d, err := os.ReadFile("/proc/schedstat"); err == nil {
			logrus.Info("init proc sched stat:\n" + string(d))
		}

		factory, _ := libcontainer.New("")
		if err := factory.StartInitialization(); err != nil {
			// as the error is sent back to the parent there is no need to log
			// or write it to stderr because the parent process will handle this
			os.Exit(1)
		}
		panic("libcontainer: container init failed to exec")
	}
}
