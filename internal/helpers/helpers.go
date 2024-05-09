package helpers

import (
	"log"
	"net/url"
	"os"
	"path"
)

func AbsolutePath(pathStart string, pathEnd string) string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Server: cant get rooted path")
	}

	//var filePath string
	if path.Base(cwd) == "metric-service" { //проверяем, запущено из корня бинарником тестов или нет
		return path.Join(pathStart, cwd, pathEnd)
	} else {
		//return path.Join(pathStart, cwd, "../..", pathEnd)

		absPath, _ := url.JoinPath(pathStart, cwd, "../..", pathEnd)
		return absPath
	}
}
