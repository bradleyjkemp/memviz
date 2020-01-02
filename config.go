package memviz

type InlineSize uint8

type Config struct {
	maxItemsToInline     int
	includePrivateFields bool
	abbreviatedTypeNames bool
	maxDepth             uint32
}

func (c *Config) MaxDepth() uint32 {
	return c.maxDepth
}

// returns whether the mapping uses abbreviated type names or the full pkg+type names.
func (c *Config) UsesAbbreviatedTypeNames() bool {
	return c.abbreviatedTypeNames
}

// returns whether the mapping includes private fields.
func (c *Config) IncludePrivateFields() bool {
	return c.includePrivateFields
}

// returns the maximum number of items to inline.
func (c *Config) MaxItemsToInline() int {
	return c.maxItemsToInline
}

// set the maximum depth to be rendered during the mapping of embedded structs.
func MaxDepth(maxDepth uint32) Configurator {
	return func(config *Config) {
		config.maxDepth = maxDepth
	}
}

// sets the whether the mapping uses abbreviated type names or the full pkg+type names.
func UseAbbreviatedTypeNames(abbreviatedTypeNames bool) Configurator {
	return func(config *Config) {
		config.abbreviatedTypeNames = abbreviatedTypeNames
	}
}

// sets whether private fields are included.
func IncludePrivateFields(includePrivateFields bool) Configurator {
	return func(config *Config) {
		config.includePrivateFields = includePrivateFields
	}
}

// sets the maximum number of items to inline.
func MaxItemsToInline(maxItems int) Configurator {
	return func(config *Config) {
		config.maxItemsToInline = maxItems
	}
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
		maxItemsToInline:     2,
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
