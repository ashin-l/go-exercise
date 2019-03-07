package common

type AppConfig struct {
	LogPath    string
	LogLevel   string
	KafkaAddrs []string
	ESAddr     string
	EtcdAddrs  []string
}

type Item struct {
	Id      string
	Type    string
	Content string
}
