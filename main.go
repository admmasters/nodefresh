package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type nodeFresh struct {
	root       string
	candidates []candidate
	config     *config
}

type candidate struct {
	name string
	path string
}

type logLevel string

const (
	debugLogLevel   logLevel = "debug"
	infoLogLevel    logLevel = "info"
	warningLogLevel logLevel = "warning"
	errorLogLevel   logLevel = "error"
)

type config struct {
	Debug    bool
	LogLevel logLevel
}

func main() {

	n := &nodeFresh{
		config: getConfig(),
	}

	log.SetLevel(log.DebugLevel)

	n.getRoot()
	n.getCandidateFolders()
	n.deleteCandidates()

}

func getConfig() *config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var c config
	err = envconfig.Process("nodefresh", &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &c
}

func (n *nodeFresh) assertLogLevel() {
	logger := logrus.WithFields(logrus.Fields{
		"logLevel": n.config.LogLevel,
	})
	switch n.config.LogLevel {
	case "debug":
	case "info":
	case "warning":
	case "error":
		return
	default:
		logger.Panic("We don't support this loglevel")
	}
}

func (n *nodeFresh) getRoot() {
	root := os.Args[1:]
	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	if len(root) > 1 {
		panic(fmt.Sprintf("We only support a single input folder"))
	}

	if len(root) == 0 {
		root = append(root, wd)
	}

	n.root = root[0]
}

func (n *nodeFresh) print() {
	for _, candidate := range n.candidates {
		log.Infoln(candidate)
	}
}

func (n *nodeFresh) deleteCandidates() {
	for _, candidate := range n.candidates {
		log.Infoln(`Deleting =====>`, candidate)
		candidate.delete()
	}
}

func (c *candidate) delete() {
	os.RemoveAll(filepath.Join(c.path, c.name))
}

func (n *nodeFresh) getCandidateFolders() {
	n.getCandidateFolderAtPath(n.root)
}

func (n *nodeFresh) getCandidateFolderAtPath(path string) {
	log.WithFields(log.Fields{
		"path": path,
	}).Debug("Started getting candidate folders")

	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	foundFolder := func(f os.FileInfo) {
		n.getCandidateFolderAtPath(filepath.Join(path, f.Name()))
	}

	foundNodeModules := func(f os.FileInfo) {
		fp := filepath.Join(path, f.Name())
		log.WithFields(log.Fields{
			"file": fp,
		}).Debug("Added delete candidate")

		c := candidate{
			name: f.Name(),
			path: path,
		}
		n.candidates = append(n.candidates, c)
	}

	for _, file := range files {
		n.getCandidateFolder(file, foundFolder, foundNodeModules)
	}
}

func (n *nodeFresh) getCandidateFolder(file os.FileInfo, f func(f os.FileInfo), g func(f os.FileInfo)) {
	if !file.IsDir() {
		log.WithFields(log.Fields{
			"file": file.Name(),
		}).Debug("Is not a dir")
		return
	}
	if file.Name() != "node_modules" {
		log.WithFields(log.Fields{
			"file": file.Name(),
		}).Debug("Is not node_modules")
		f(file)
		return
	}

	g(file)

}
