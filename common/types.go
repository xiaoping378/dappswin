package common

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// JSONTime 自定义时间，为了给前端返回符合2006-01-02 15:04:05的格式
type JSONTime struct {
	Time time.Time
}

const (
	timeFormart = "2006-01-02 15:04:05"
)

// func (t *JSONTime) UnmarshalJSON(data []byte) (err error) {
// 	now, err := time.ParseInLocation(`"`+timeFormart+`"`, string(data), time.Local)
// 	t.Time = now
// 	return
// }

// MarshalJSON 返回前端的json格式化方法
func (t JSONTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormart)+2)
	b = append(b, '"')
	b = t.Time.AppendFormat(b, timeFormart)
	b = append(b, '"')
	return b, nil
}

// Value insert timestamp into mysql need this function.
func (t JSONTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof time.Time
func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
