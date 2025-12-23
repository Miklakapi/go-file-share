package ports

type PasswordHasher interface {
	Hash(plain string) (string, error)
	Verify(plain, hash string) bool
}
