package common

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"protocol"
	"strconv"
	"time"
	"unicode"
)

func ParseTimeDuration(value string) (int64, error) {
	parsedFloat, err := strconv.ParseFloat(value[:len(value)-1], 64)
	if err != nil {
		return 0, err
	}

	switch value[len(value)-1] {
	case 'u':
		return int64(parsedFloat * float64(time.Microsecond)), nil
	case 's':
		return int64(parsedFloat * float64(time.Second)), nil
	case 'm':
		return int64(parsedFloat * float64(time.Minute)), nil
	case 'h':
		return int64(parsedFloat * float64(time.Hour)), nil
	case 'd':
		return int64(parsedFloat * 24 * float64(time.Hour)), nil
	case 'w':
		return int64(parsedFloat * 7 * 24 * float64(time.Hour)), nil
	}

	lastChar := value[len(value)-1]
	if !unicode.IsDigit(rune(lastChar)) && lastChar != '.' {
		return 0, fmt.Errorf("Invalid character '%c'", lastChar)
	}

	if value[len(value)-2] != '.' {
		extraDigit := float64(lastChar - '0')
		parsedFloat = parsedFloat*10 + extraDigit
	}
	return int64(parsedFloat), nil
}

func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func StringToSeriesArray(seriesString string, args ...interface{}) ([]*protocol.Series, error) {
	seriesString = fmt.Sprintf(seriesString, args...)
	series := []*protocol.Series{}
	err := json.Unmarshal([]byte(seriesString), &series)
	return series, err
}

func CurrentTime() int64 {
	return time.Now().UnixNano() / int64(1000)
}

func RingLocation(database *string, timeSeries *string, time *int64) int {
	hasher := sha1.New()
	hasher.Write([]byte(fmt.Sprintf("%s%s%d", *database, *timeSeries, *time)))
	buf := bytes.NewBuffer(hasher.Sum(nil))
	var n int64
	binary.Read(buf, binary.LittleEndian, &n)
	nInt := int(n)
	if nInt < 0 {
		nInt = nInt * -1
	}
	return nInt
}
