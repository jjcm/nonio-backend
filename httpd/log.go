package main

func log(message interface{}) {
	sociConfig.Logger.Info(message)
}

func logError(message interface{}) {
	sociConfig.Logger.Error(message)
}
