package usecase

import "github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"

func (c *ClientUsecase) deleteSecret(secretName string, success bool) {
	defer c.InMenu()
	if !success {
		return
	}

	if err := c.fileManager.DeleteByName(secretName); err != nil {
		c.ui.PrintErr(ui.DeleteErr)

		c.log.Err(err).Msg("err del file")
	}
}
