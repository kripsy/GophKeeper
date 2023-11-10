//nolint:staticcheck
package usecase

import (
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
)

// InMenu continuously prompts the user with a menu,
// allowing them to interact with the application based on their choices.
func (c *ClientUsecase) InMenu() {
	switch ui.MenuTable[c.ui.Menu(c.userData.Meta.IsSyncStorage)] {
	// If the user chooses to view secrets, invoke the getSecrets method.
	case ui.SecretsKey:
		c.getSecrets(c.ui.GetSecret(c.userData.Meta.Data))
		// If the user chooses to add a secret, invoke the createSecret method.
	case ui.AddSecretKey:
		c.createSecret(c.ui.ChooseSecretType())
		// If the user chooses to update a secret, invoke the updateSecret method.
	case ui.UpdateSecretKey:
		c.updateSecret(c.ui.UpdateSecret(c.userData.Meta.Data))
		// If the user chooses to sync secrets, invoke the sync method.
	case ui.SyncSecrets:
		c.sync()
		// If the user chooses to delete a secret, invoke the deleteSecret method.
	case ui.DeleteSecretKey:
		c.deleteSecret(c.ui.DeleteSecret(c.userData.Meta.Data))
		// If the user chooses to view information about the application, invoke the about method.
	case ui.About:
		c.about()
	case ui.ExitKey:
		// If the user chooses to exit the application, invoke the exit method from the UI.
		c.ui.Exit()
	}

	// Recursive call to keep the menu loop running.
	c.InMenu()
}
