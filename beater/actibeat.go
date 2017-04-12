package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/datake914/actibeat/config"
)

type Actibeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Actibeat{
		done:   make(chan struct{}),
		config: config,
	}
	return bt, nil
}

func (bt *Actibeat) Run(b *beat.Beat) error {
	logp.Info("actibeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		var actispy Actispy = new(ActispyWin32)

		event := common.MapStr{
			"@timestamp": common.Time(time.Now()),
			"type":       b.Name,
			"procid":     actispy.getProcessID(),
			"procname":   actispy.getProcessName(),
			"windowname": actispy.getWindowName(),
			"username":   actispy.getUserName(),
			"interval": common.MapStr{
				"sec":    bt.config.Period.Seconds(),
				"minute": bt.config.Period.Minutes(),
				"hour":   bt.config.Period.Hours(),
			},
		}
		bt.client.PublishEvent(event)
		logp.Info("Event sent")
		counter++
	}
}

func (bt *Actibeat) Stop() {
	bt.client.Close()
	close(bt.done)
}