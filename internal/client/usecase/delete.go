package usecase

import "github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"

// Deletes the file and metadata information by name.
func (c *ClientUsecase) deleteSecret(secretName string, success bool) {
	if !success {
		return
	}

	if err := c.fileManager.DeleteByName(secretName); err != nil {
		c.ui.PrintErr(ui.DeleteErr)

		c.log.Err(err).Msg("err del file")
	}
}
