package qcloud

import "github.com/ghinknet/smsutils/v3/driver"

func init() {
	driver.Register(Name, Driver{})
}
