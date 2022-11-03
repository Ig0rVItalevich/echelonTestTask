package cache

import "github.com/go-redis/redis"

type Cache struct {
	DB *redis.Client
}

func NewCache(DB *redis.Client) *Cache {
	return &Cache{DB: DB}
}

func (c *Cache) SetThumbnail(link string, thumbnailBytes []byte) error {
	return c.DB.Set(link, thumbnailBytes, 0).Err()
}

func (c *Cache) GetThumbnail(link string) ([]byte, error) {
	resStruct := c.DB.Get(link)
	thumbnail, err := resStruct.Bytes()

	return thumbnail, err
}
