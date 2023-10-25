package usecase

func (c *ClientUsecase) deleteSecret(secretName string, success bool) {
	defer c.InMenu()
	if !success {
		return
	}

	c.fileManager.DeleteByName(secretName)
}
