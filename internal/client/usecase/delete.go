package usecase

func (c *ClientUsecase) deleteSecret(secretName string, success bool) {
	defer c.InMenu()
	if !success {
		return
	}

	if err := c.fileManager.DeleteByName(secretName); err != nil {

	}

}
