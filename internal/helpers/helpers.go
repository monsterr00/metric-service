package helpers

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path"

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
