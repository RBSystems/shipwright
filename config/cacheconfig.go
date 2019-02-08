package config

//Cache .
type Cache struct {
	Name string `json:"name"`

	//Persistence type
	StorageType string `json:"storage-type"`

	//Legacy or Defautl
	CacheType string `json:"cache-type"`

	CouchInfo CouchCache `json:"couch-cache"`
	ELKinfo   ElkCache   `json:"elk-cache"`
	RedisInfo RedisCache `json:"redis-cache"`
}

//CouchCache .
type CouchCache struct {
	DatabaseName string `json:"database-name"`
	URL          string `json:"url"`
}

//ElkCache .
type ElkCache struct {
	DeviceIndex string `json:"device-index"`
	URL         string `json:"url"`
}

//RedisCache .
type RedisCache struct {
	DevDatabsae  int    `json:"device-database"`
	RoomDatabsae int    `json:"room-database"`
	Password     string `json:"password"`
	URL          string `json:"url"`
}
