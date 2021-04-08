package cfg

type Config struct {
	Version  string             `yaml:"version"`
	Features map[string]Feature `yaml:"features"`
}

type Feature struct {
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	Field  string   `yaml:"field"`
	Fields []string `yaml:"fields"`

	Include []string `yaml:"include"`
	Exclude []string `yaml:"exclude"`
	Weight  int      `yaml:"weight"`
}

func (c *Config) Append(a Config) {
	if c.Version == "" {
		c.Version = a.Version
	}
	if c.Features == nil {
		c.Features = map[string]Feature{}
	}
	for name, feature := range a.Features {
		if f, ok := c.Features[name]; ok {
			f.Rules = append(f.Rules, a.Features[name].Rules...)
			c.Features[name] = f
		} else {
			c.Features[name] = feature
		}
	}
}
