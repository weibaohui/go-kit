package datekit

import "time"

// FormatToTime 将字符串按格式转换为时间
func StringToTime(dateStr, layout string) (time.Time, error) {
	location, err := time.LoadLocation("Local")
	if err != nil {
		return time.Time{}, err
	}
	timeTmp, err := time.ParseInLocation(layout, dateStr, location)
	if err != nil {
		return time.Time{}, err
	}
	return timeTmp, nil
}

func TimeToString(t time.Time, layout string) string {
	return t.Format(layout)
}

func TimeToStringNormal(t time.Time) string {
	return t.Format(normalLayout)
}
