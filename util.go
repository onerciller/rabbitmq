package rmq

import (
	"fmt"
	"math/rand"
)

func generateChannelName(appName string) string {
	return fmt.Sprintf("%s-%d", appName, rand.Int())
}
