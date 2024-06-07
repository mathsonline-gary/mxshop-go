package config

// TODO: Add dynamic configuration package code here, learn from kratos: https://github.com/go-kratos/kratos/tree/main/config

// Loader is the interface that loads configuration.
type Loader interface {
	Load() error
}

// Watcher is the interface that watches configuration changes.
type Watcher interface {
	Watch() error
}

// Closer is the interface that closes the configuration.
type Closer interface {
	Close() error
}
