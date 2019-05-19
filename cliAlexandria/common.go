package cliAlexandria

import "time"

func ToIoError(err error) string {
	return "IO error while sending command: " + err.Error() + "\n"
}

func formatTime(t int64) string {
	return time.Unix(t, 0).Format(time.UnixDate)
}
