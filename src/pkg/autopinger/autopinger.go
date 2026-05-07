package autopinger

import (
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func Start() {
	url := os.Getenv("PING_URL")
	if url == "" {
		logrus.Info("[AutoPinger] PING_URL not set, skipping auto-pinger")
		return
	}

	logrus.Infof("[AutoPinger] Starting → pinging %s every 10 minutes", url)

	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			resp, err := http.Get(url)
			if err != nil {
				logrus.Errorf("[AutoPinger] Ping failed: %v", err)
				continue
			}
			resp.Body.Close()
			logrus.Infof("[AutoPinger] Ping OK → %s [%d]", url, resp.StatusCode)
		}
	}()
}
