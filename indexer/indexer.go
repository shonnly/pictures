package indexer

import (
	"errors"
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/db"
	"github.com/dtylman/pictures/indexer/picture"
	"github.com/jasonwinn/geocoder"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"os"
	"path/filepath"
)

func init() {
	exif.RegisterParsers(mknote.All...)
}

func Start(options Options) error {
	err := SetStarted(options)
	if err != nil {
		return err
	}
	if options.IndexLocation == true {
		//geocoder.SetAPIKey("8cCGEGGioKhpCLPjhAG44NfXYaXs9jCk")
		if conf.Options.MapQuestAPIKey == "" {
			return errors.New("API KEY for map quest is empty")
		}
		geocoder.SetAPIKey(conf.Options.MapQuestAPIKey)
	}
	go indexPictures()
	return nil
}

func Stop() {
	SetDone()
}

func Status() {

}

func indexOne(path string, info os.FileInfo, e1 error) error {
	if e1 != nil {
		AddError(path, e1)
		return nil
	}
	if !IsRunning() {
		return errors.New("Indexer had stopped")
	}
	if !info.IsDir() {
		if info.Size() > 0 {
			i, err := picture.NewIndex(path, info)
			if err != nil {
				AddError(path, err)
			} else {
				saveIndex(path, i)
			}
		}
	}
	return nil
}

func index(rootPath string) error {
	return filepath.Walk(rootPath, indexOne)
}

func indexPictures() {
	defer SetDone()
	for _, folder := range conf.Options.SourceFolders {
		err := index(folder)
		if err != nil {
			AddError(folder, err)
		}
	}
}

func saveIndex(path string, i *picture.Index) {
	if GetOptions().IndexLocation {
		err := i.PopulateLocation()
		if err != nil {
			AddError(path, err)
		}
	}
	err := db.Index(i)
	if err != nil {
		AddError(path, err)
	}
}
