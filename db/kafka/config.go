package kafka

import "errors"

type Config struct {
	Brokers   []string
	GroupID   string
	Topic     string
	Partition int64
	User      string
	Password  string
}

func (c *Config) Check() error {
	if c.Brokers == nil || len(c.Brokers) == 0 {
		return errors.New("can not find kafka brokers config")
	}

	if c.GroupID == "" {
		return errors.New("can not find kafka groupID config")
	}

	if c.Partition == 0 {
		return errors.New("can not find kafka partition config or value cannot be set to 0")
	}
	return nil
}
