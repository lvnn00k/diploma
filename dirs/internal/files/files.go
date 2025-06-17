package files

import (
	"fmt"
	"io"
	"os"
)

type Data struct {
	directory string
}

func New(pathToFiles string) *Data {
	return &Data{pathToFiles}
}

func (d *Data) NewFile(name, place string, file io.Reader) error {

	path := fmt.Sprintf("%s/%s/%s", d.directory, place, name)

	if _, err := os.Stat(path); err != nil {
		newfile, err := os.Create(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(newfile, file)
		if err != nil {
			return err
		}

		newfile.Close()

		return nil
	}

	return fmt.Errorf("file exist")

}

func (d *Data) NewDir(path string) error {

	err := os.Mkdir(fmt.Sprintf("%s%s", d.directory, path), 0755)
	if err != nil {
		return err
	}

	return nil
}

func (d *Data) DataRemove(path string) error {

	err := os.RemoveAll(fmt.Sprintf("%s/%s", d.directory, path))
	if err != nil {
		return err
	}

	return nil
}
