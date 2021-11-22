package services

type PageHandler interface {
	SetInitialState()
	RehydrateSession()
	SetupBindings()
}
