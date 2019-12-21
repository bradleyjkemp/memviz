package memviz

const (
	InlineSizeInfinite = 0
	InlineSizeLow      = 2
	InlineSizeMed      = 5
	InlineSizeHigh     = 10
)

type InlineSize uint8

type Config struct {
	maxInliningSize      InlineSize
	includePrivateFields bool
	abbreviatedTypeNames bool
	maxDepth             uint32
}

func (c *Config) MaxDepth() uint32 {
	return c.maxDepth
}

func (c *Config) SetMaxDepth(maxDepth uint32) {
	c.maxDepth = maxDepth
}

func (c *Config) AbbreviatedTypeNames() bool {
	return c.abbreviatedTypeNames
}

func (c *Config) SetAbbreviatedTypeNames(abbreviatedTypeNames bool) {
	c.abbreviatedTypeNames = abbreviatedTypeNames
}

func (c *Config) IncludePrivateFields() bool {
	return c.includePrivateFields
}

func (c *Config) SetIncludePrivateFields(includePrivateFields bool) {
	c.includePrivateFields = includePrivateFields
}

func (c *Config) MaxInliningSize() InlineSize {
	return c.maxInliningSize
}

func (c *Config) SetMaxInliningSize(maxInliningSize InlineSize) {
	c.maxInliningSize = maxInliningSize
}

type Configurator func(*Config)

// Default configuration
//
// Max inlining size: InlineSizeInfininte // unlimited
// Include private fields: true
// Abbreviated type names: true
// Max Depth 0 // unlimited
func defaultConfig() *Config {
	return &Config{
		maxInliningSize:      InlineSizeLow,
		includePrivateFields: true,
		abbreviatedTypeNames: true,
		maxDepth:             0,
	}
}

func New(configurators ...Configurator) *Config {
	config := defaultConfig()
	for _, configurator := range configurators {
		configurator(config)
	}
	return config
}
