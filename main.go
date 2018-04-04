package main

import (
	"flag"
	"fmt"
	"os"

	lg "github.com/hellozee/cook/logger"
	mg "github.com/hellozee/cook/manager"
	ps "github.com/hellozee/cook/parser"
	wk "github.com/hellozee/cook/worker"
)

var cleanFlag = flag.Bool("clean", false, "To clean the cached data")
var helpFlag = flag.Bool("help", false, "To show this help message")

func main() {

	flag.Parse()

	help := `Usage: cook [OPTIONS]

	--help:
		To show this help message

	--clean:
		To clean the cached data
		
	`
	if *helpFlag == true {
		fmt.Println(help)
		return
	}

	if *cleanFlag == true {
		os.RemoveAll("Cooking/")
		return
	}

	logger := lg.NewLogger()

	//Reading the Recipe File
	manager, err := mg.NewManager(&logger)
	must(err, &logger)

	//Parsing the Recipe File
	parser := ps.NewParser(manager.FileData, &logger)
	err = parser.Parse()
	must(err, &logger)

	worker := wk.NewWorker(&logger)
	err = manager.GenerateFileList(parser, parser.CompilerDetails.Start)
	must(err, &logger)

	if _, err := os.Stat("Cooking/details.json"); err == nil {

		err = manager.ReadDetails()
		must(err, &logger)

		err = worker.CompareAndCompile(parser, &manager)
		must(err, &logger)

	} else {
		_ = os.Mkdir("Cooking", 0755)
		err = manager.GenerateList()
		must(err, &logger)
		err = worker.CompileFirst(parser, manager)
		must(err, &logger)
	}

	err = worker.Link(parser)
	must(err, &logger)

	err = manager.WriteDetails()
	must(err, &logger)

	logger.WriteLog()

	fmt.Println("Build finished, logs reported to Cooking/log")

}

func must(err error, log *lg.Logger) {
	if err != nil {
		fmt.Println("Something went wrong, please check Cooking/log/build.errors")
		log.WriteLog()
		os.Exit(-1)
	}
}
