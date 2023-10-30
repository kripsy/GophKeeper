package usecase

import (
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
)

func (c *ClientUsecase) InMenu() {
	action := ui.MenuTable[c.ui.Menu(c.userData.Meta.IsLocalStorage)]
	switch action {
	case ui.SecretsKey:
		c.getSecrets(c.ui.GetSecret(c.userData.Meta.Data)) // todo сторит обдумать метадату на name
	case ui.AddSecretKey:
		c.createSecret(c.ui.ChooseSecretType())
	case ui.UpdateSecretKey: //todo идут строительные работы
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
