package utils

var influxDbConnection = ""

func SetConnectionUrl(u string) {
	influxDbConnection = u
}
