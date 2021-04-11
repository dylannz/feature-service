package cfg

type Config struct {
	Version  string             `yaml:"version"`
	Features map[string]Feature `yaml:"features"`
}

type Feature struct {
	Rules Rules `yaml:"rules"`
}

type Rules struct {
	Enable  []EnableRule  `yaml:"enable"`
	Disable []DisableRule `yaml:"disable"`
}

type EnableRule struct {
	Field  string   `yaml:"field"`
	Fields []string `yaml:"fields"`

	Values MatchValues `yaml:"values"`
	Weight int         `yaml:"weight"`
}

type DisableRule struct {
	Field  string   `yaml:"field"`
	Fields []string `yaml:"fields"`

	Values MatchValues `yaml:"values"`
}

type MatchValues struct {
	Eq []string `json:"eq"`
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
			f.Rules.Enable = append(f.Rules.Enable, a.Features[name].Rules.Enable...)
			f.Rules.Disable = append(f.Rules.Disable, a.Features[name].Rules.Disable...)
			c.Features[name] = f
		} else {
			c.Features[name] = feature
		}
	}
}
