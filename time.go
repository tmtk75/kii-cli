package kiicli

import "time"

func timeFromUnix(s int64) time.Time {
	t := time.Unix(s/1000, s%1000*1000*1000)
	if globalConfig.UTC {
		return t.UTC()
	}
	return t
}

func timeFromStringInUTC(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return t, err
	}
	if globalConfig.UTC {
		return t, nil
	}
	return t.Local(), nil
}
