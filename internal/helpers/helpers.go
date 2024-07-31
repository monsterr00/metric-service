package helpers

func AbsolutePath(pathStart string, pathEnd string) string {
	/*cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Server: cant get rooted path")
	}
	p := path.Base(cwd)*/

	return "file:///Users/denis/metric-service/db/migrations"

	/*
		if p == "metric-service" { //проверяем, запущено из корня бинарником тестов или нет
			return path.Join(pathStart, cwd, pathEnd)
		} else if p == "httplayer" {
			absPath, _ := url.JoinPath(pathStart, cwd, "../../..", pathEnd)
			return absPath
		} else {
			absPath, _ := url.JoinPath(pathStart, cwd, "../..", pathEnd)
			return absPath
		}*/
}
