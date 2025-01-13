package signout

type SignOutRepository interface{}

type signOutRepository struct{}

func NewSignOutRepository() SignOutRepository {
	return &signOutRepository{}
}
