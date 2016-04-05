package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"math/rand"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
)

func abs(n int) int {
	if n >= 0 {
		return n
	}
	return -n
}

func changeDesktopBackground(path string) error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s/Library/Application Support/Dock/desktoppicture.db", usr.HomeDir))
	if err != nil {
		return err
	}
	defer db.Close()

	sqlSmt := fmt.Sprintf("update data set value = '%s';", path)

	_, err = db.Exec(sqlSmt)
	if err != nil {
		return err
	}

	cmd := exec.Command("killall", "Dock")
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func getScreenResolution() (int, int) {
	cmd := "system_profiler SPDisplaysDataType |grep Resolution |tr 'x' '\n' |sed 's/@.*//' |sed 's/[^0-9]//g'"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return 1920, 1080
	}
	s := strings.Split(string(out), "\n")
	if len(s) >= 2 {
		width := 1920
		height := 1080
		for i := 1; i < len(s); i += 2 {
			w, err := strconv.Atoi(s[i-1])
			if err != nil {
				return 1920, 1080
			}
			h, err := strconv.Atoi(s[i])
			if err != nil {
				return 1920, 1080
			}
			if w*h > width*height {
				width = w
				height = h
			}
		}
		return width, height
	}
	return 1920, 1080
}

func randBool() bool {
	return rand.Intn(2) == 0
}

func randMinMax(min int, max int) int {
	if min == max {
		return min
	}
	return rand.Intn(max-min) + min
}
