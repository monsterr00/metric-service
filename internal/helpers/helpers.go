package helpers

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/monsterr00/metric-service.gittest_client/internal/config"
)

// AbsolutePath возвращает путь до файла в зависимости от режима запуска программы.
func AbsolutePath(pathStart string, pathEnd string) string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Server: cant get rooted path")
	}
	p := path.Base(cwd)

	if p == "metric-service" { //проверяем, запущено из корня бинарником тестов или нет
		return path.Join(pathStart, cwd, pathEnd)
	} else if p == "httplayer" || p == "appplayer" || p == "storelayer" {
		absPath, _ := url.JoinPath(pathStart, cwd, "../../..", pathEnd)
		return absPath
	} else {
		absPath, _ := url.JoinPath(pathStart, cwd, "../..", pathEnd)
		return absPath
	}
}

// PrintBuildInfo выводит информацию в stdout о версии сборки.
// В случае, если нет данных о версии, выводится значение "N/A".
//
// Пример выполнения команды build.
//
// BUILD_DATE="$(date -u +'%Y-%m-%dT%H:%M:%SZ')"
// BUILD_COMMIT="$(git rev-parse HEAD)"
// BUILD_VERSION="$(git describe --all --abbrev=0)"
//
// go build -C cmd/server -buildvcs=false -o server -ldflags="-X 'github.com/monsterr00/metric-service.gittest_client/internal/config.buildDate=${BUILD_DATE}' -X 'github.com/monsterr00/metric-service.gittest_client/internal/config.buildCommit=${BUILD_COMMIT}' -X 'github.com/monsterr00/metric-service.gittest_client/internal/config.buildVersion=${BUILD_VERSION}'"
func PrintBuildInfo() {
	versionInfo := config.GetVersionInfo()
	fmt.Printf("Build version: %s\n", versionInfo.BuildVersion)
	fmt.Printf("Build date: %s\n", versionInfo.BuildDate)
	fmt.Printf("Build commit: %s\n", versionInfo.BuildCommit)
}

// ReadConfigJSON считывает json-файл с настройками и записывает результат в мапу
func ReadConfigJSON(path string) (map[string]string, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	env_conf := make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		splitedStr := strings.Split(scanner.Text(), `"`)

		if len(splitedStr) > 1 {
			if len(splitedStr) > 3 {
				env_conf[splitedStr[1]] = splitedStr[3]
			} else {
				env_conf[splitedStr[1]] = strings.TrimRight(strings.TrimLeft(splitedStr[2], ": "), ",")
			}
		}
	}
	if len(env_conf) == 0 {
		return nil, errors.New("server config json is empty")
	}
	return env_conf, nil
}

func IsFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
