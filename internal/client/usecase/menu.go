package usecase

import (
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
)

func (c *ClientUsecase) InMenu() {
	switch ui.MenuTable[c.ui.Menu(c.userData.Meta.IsSyncStorage)] {
	case ui.SecretsKey:
		c.getSecrets(c.ui.GetSecret(c.userData.Meta.Data))
	case ui.AddSecretKey:
		c.createSecret(c.ui.ChooseSecretType())
	case ui.UpdateSecretKey:
		c.updateSecret(c.ui.UpdateSecret(c.userData.Meta.Data))
	case ui.SyncSecrets:
		c.sync()
	case ui.DeleteSecretKey:
		c.deleteSecret(c.ui.DeleteSecret(c.userData.Meta.Data))
	case ui.About:
		c.about()
	case ui.ExitKey:
		c.ui.Exit()
	}
}
