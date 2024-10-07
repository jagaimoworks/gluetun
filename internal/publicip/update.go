package publicip

import (
	"fmt"
	"os"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (l *Loop) update(partialUpdate settings.PublicIP) (err error) {
	// No need to lock the mutex since it can only be written
	// in the code below in this goroutine.
	updatedSettings, err := l.settings.UpdateWith(partialUpdate)
	if err != nil {
		return err
	}

	if *l.settings.IPFilepath != *updatedSettings.IPFilepath {
		switch {
		case *l.settings.IPFilepath == "":
			err = persistPublicIP(*updatedSettings.IPFilepath,
				l.ipData.IP.String(), l.puid, l.pgid)
			if err != nil {
				return fmt.Errorf("persisting ip data: %w", err)
			}
		case *updatedSettings.IPFilepath == "":
			err = os.Remove(*l.settings.IPFilepath)
			if err != nil {
				return fmt.Errorf("removing ip data file path: %w", err)
			}
		default:
			err = os.Rename(*l.settings.IPFilepath, *updatedSettings.IPFilepath)
			if err != nil {
				return fmt.Errorf("renaming ip data file path: %w", err)
			}
		}
	}

	l.settingsMutex.Lock()
	l.settings = updatedSettings
	l.settingsMutex.Unlock()

	return nil
}
