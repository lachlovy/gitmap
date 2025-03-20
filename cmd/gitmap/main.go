package main

import (
	"flag"
	"github.com/lachlovy/gitmap/pkg"
	"log"
)

func main() {
	var scanDir string
	var configFile string
	var dateRange string

	flag.StringVar(&scanDir, "scan-dir", "", "Directory to be scanned")
	// not implemented
	flag.StringVar(&configFile, "config-file", "", "Configuration file to use")
	flag.StringVar(&dateRange, "date-range", "Year", "Date range to scan, supported values are Year, Month and SixMonth")
	flag.Parse()

	if configFile != "" {
		log.Println("Warning: config-file flag is not implemented yet")
	}

	var allGitRepositories = make([]string, 0)

	if scanDir != "" {
		allGitRepositories = pkg.ScanGitRepositories(scanDir)
	} else if scanDir == "" {
		scanDir = "."
		allGitRepositories = pkg.ScanGitRepositories(scanDir)
	}

	res, err := pkg.GetGitRepositoriesStatistics(allGitRepositories, dateRange)
	if err != nil {
		log.Fatal(err)
	}
	pkg.DrawContributionPlot(res)
}
