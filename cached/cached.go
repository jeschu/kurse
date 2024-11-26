package cached

import (
	"errors"
	"kurse/lang"
	"os"
	"path"
	"time"
)

var (
	cacheDir string
)

func init() {
	var err error
	cacheDir, err = os.UserCacheDir()
	lang.FatalOnError(err)
}

func Load[T any](application string, cache string, maxAge time.Duration, mapper func([]byte) *T) (*T, bool) {
	var (
		fi   os.FileInfo
		data []byte
		err  error
	)
	cacheFile := ensureCacheFile(application, cache)
	fi, err = os.Stat(cacheFile)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false
	}
	age := time.Now().Sub(fi.ModTime())
	if age > maxAge {
		return nil, false
	}
	data, err = os.ReadFile(cacheFile)
	lang.FatalOnError(err)
	return mapper(data), true
}

func Save[T any](application string, cache string, obj *T, mapper func(*T) []byte) {
	cacheFile := ensureCacheFile(application, cache)
	data := mapper(obj)
	err := os.WriteFile(cacheFile, data, 0644)
	lang.FatalOnError(err)
}

func ensureCacheFile(application, cache string) string {
	dir := path.Join(cacheDir, application)
	err := os.MkdirAll(dir, 0744)
	lang.FatalOnError(err)
	return path.Join(dir, cache)
}
