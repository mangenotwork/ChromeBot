package host

import (
	"time"

	gt "github.com/mangenotwork/gathertool"
)

func Ping(ip string) (time.Duration, error) {
	return gt.Ping(ip)
}
