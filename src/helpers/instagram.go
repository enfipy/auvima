package helpers

import (
	"os"

	goinsta "github.com/ahmdrz/goinsta/v2"
)

func InitInstagramClient(username, password, credsPath string) *goinsta.Instagram {
	var err error
	var insta *goinsta.Instagram

	if _, err = os.Stat(credsPath); err == nil {
		insta, err = goinsta.Import(credsPath)
		PanicOnError(err)
	} else {
		insta = goinsta.New(username, password)
		err = insta.Login()
		PanicOnError(err)

		err = insta.Export(credsPath)
		PanicOnError(err)
	}

	return insta
}
