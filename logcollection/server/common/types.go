package common

type AppConfig struct {
	LogPath    string
	LogLevel   string
	KafkaAddrs []string
	ESAddr     string
	EtcdAddrs  []string
}

type Item struct {
	Content string
}
