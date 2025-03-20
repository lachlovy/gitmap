package pkg

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"time"
)

func getDateRange(dateRange string) (int, error) {
	currentDate := time.Now()
	switch dateRange {
	case "Year":
		lastYearDate := currentDate.AddDate(-1, 0, 0)
		lastYearWeekDay := int(lastYearDate.Weekday())
		if lastYearWeekDay == 0 {
			lastYearWeekDay = 7
		}
		daysToMonday := 1 - lastYearWeekDay
		lastYearMonday := lastYearDate.AddDate(0, 0, daysToMonday)
		lastYearMonday = time.Date(lastYearMonday.Year(), lastYearMonday.Month(), lastYearMonday.Day(), 0, 0, 0, 0, lastYearMonday.Location())
		currentDateStart := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, currentDate.Location())
		daysDiff := currentDateStart.Sub(lastYearMonday).Hours()/24 + 1
		return int(daysDiff), nil
	case "SixMonth":
		lastSixMonthDate := currentDate.AddDate(0, -6, 0)
		lastSixMonthWeekDay := int(lastSixMonthDate.Weekday())
		if lastSixMonthWeekDay == 0 {
			lastSixMonthWeekDay = 7
		}
		daysToMonday := 1 - lastSixMonthWeekDay
		lastSixMonthMonday := lastSixMonthDate.AddDate(0, 0, daysToMonday)
		lastSixMonthMonday = time.Date(lastSixMonthMonday.Year(), lastSixMonthMonday.Month(), lastSixMonthMonday.Day(), 0, 0, 0, 0, lastSixMonthMonday.Location())
		currentDateStart := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, currentDate.Location())
		daysDiff := currentDateStart.Sub(lastSixMonthMonday).Hours()/24 + 1
		return int(daysDiff), nil
	case "Month":
		lastMonthDate := currentDate.AddDate(0, -1, 0)
		lastMonthWeekDay := int(lastMonthDate.Weekday())
		if lastMonthWeekDay == 0 {
			lastMonthWeekDay = 7
		}
		daysToMonday := 1 - lastMonthWeekDay
		lastMonthMonday := lastMonthDate.AddDate(0, 0, daysToMonday)
		lastMonthMonday = time.Date(lastMonthMonday.Year(), lastMonthMonday.Month(), lastMonthMonday.Day(), 0, 0, 0, 0, lastMonthMonday.Location())
		currentDateStart := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, currentDate.Location())
		daysDiff := currentDateStart.Sub(lastMonthMonday).Hours()/24 + 1
		return int(daysDiff), nil
	default:
		return 0, errors.New("invalid date range")
	}
}

func getGitRepositoryContributions(statistics []int, gitRepository string) ([]int, error) {
	gitRepo, err := git.PlainOpen(gitRepository)
	if err != nil {
		return nil, err
	}
	gitHead, err := gitRepo.Head()
	if err != nil {
		return statistics, nil
	}
	gitLog, err := gitRepo.Log(&git.LogOptions{From: gitHead.Hash()})
	if err != nil {
		return nil, err
	}
	err = gitLog.ForEach(func(c *object.Commit) error {
		commitDate := c.Author.When
		daysDiff := int(time.Now().Sub(commitDate).Hours() / 24)
		if daysDiff < len(statistics) {
			statistics[daysDiff]++
		}
		return nil
	})
	return statistics, nil
}

func DrawContributionPlot(statistics []int) {
	weekdays := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	total_days := len(statistics)

	var numColumns int
	if total_days%7 == 0 {
		numColumns = total_days / 7
	} else {
		numColumns = total_days/7 + 1
	}
	numRows := 7

	maxVal := 0
	for _, val := range statistics {
		if val > maxVal {
			maxVal = val
		}
	}

	for row := 0; row < numRows; row++ {
		if row < len(weekdays) {
			fmt.Printf("%-3s ", weekdays[row])
		} else {
			fmt.Print("    ")
		}
		for col := 0; col < numColumns; col++ {
			index := total_days - row - col*7 - 1
			if index >= 0 {
				val := statistics[index]
				bgColor, textColor := getColors(val, maxVal)
				fmt.Print(bgColor + textColor)
				fmt.Printf("%2d ", val)
				fmt.Print("\033[0m")
			} else {
				fmt.Print("   ")
			}
		}
		fmt.Println()
	}
}

func getColors(val, maxVal int) (string, string) {
	textColor := "\033[30m"

	var bgColor string
	if val == 0 {
		bgColor = "\033[48;5;240m"
		textColor = "\033[97m"
	} else {
		greenLevels := []string{
			"\033[48;5;151m",
			"\033[48;5;114m",
			"\033[48;5;78m",
			"\033[48;5;71m",
			"\033[48;5;28m",
		}

		if maxVal <= 0 {
			bgColor = greenLevels[0]
		} else {
			index := int(float64(val) / float64(maxVal) * float64(len(greenLevels)-1))
			bgColor = greenLevels[index]

			if index > 2 {
				textColor = "\033[30m"
			}
		}
	}

	return bgColor, textColor
}

func GetGitRepositoriesStatistics(gitRepositories []string, dateRange string) ([]int, error) {
	countDays, err := getDateRange(dateRange)
	if err != nil {
		return nil, err
	}

	var statistics = make([]int, countDays)

	for _, gitRepository := range gitRepositories {
		statistics, err = getGitRepositoryContributions(statistics, gitRepository)
		if err != nil {
			fmt.Printf("Warning: Error processing repository %s: %v\n", gitRepository, err)
			continue
		}
	}
	return statistics, nil
}
