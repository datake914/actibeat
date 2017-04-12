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
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}
		// Create spy.
		// TODO Use appropriate SPY according to system architecture.
		actispy, err := newActispyWin32()
		if err != nil {
			bt.HandleError(b, err)
			continue
		}

		procid, err := actispy.getProcessID()
		if err != nil {
			bt.HandleError(b, err)
			continue
		}
		procname, err := actispy.getProcessName()
		if err != nil {
			bt.HandleError(b, err)
			continue
		}
		windowname, err := actispy.getWindowName()
		if err != nil {
			bt.HandleError(b, err)
			continue
		}
		username, err := actispy.getUserName()
		if err != nil {
			bt.HandleError(b, err)
			continue
		}

		event := common.MapStr{
			"@timestamp": common.Time(time.Now()),
			"type":       b.Name,
			"procid":     procid,
			"procname":   procname,
			"windowname": windowname,
			"username":   username,
			"interval": common.MapStr{
				"sec":    bt.config.Period.Seconds(),
				"minute": bt.config.Period.Minutes(),
				"hour":   bt.config.Period.Hours(),
			},
		}
		bt.client.PublishEvent(event)
		logp.Info("Event sent")
	}
}

func (bt *Actibeat) HandleError(b *beat.Beat, err error) {
	event := common.MapStr{
		"@timestamp": common.Time(time.Now()),
		"type":       b.Name + "error",
		"error":      err.Error(),
	}
	bt.client.PublishEvent(event)
	logp.Info("Error Event sent")
}

func (bt *Actibeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
