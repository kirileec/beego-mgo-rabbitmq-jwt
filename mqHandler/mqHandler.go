package mqHandler

import (
	log "github.com/goinggo/tracelog"
)

// HandleMsg hanlde the message from a queue
func HandleMsg(body []byte) {
	log.Info("HandleMsg Info", "HandleMsg", "body: %s", body)
}
