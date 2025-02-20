package storage

type Storage interface {
	Save(url string) (string, error)
	Get(short string) (string, error)
	// Dump() map[string]string
}
