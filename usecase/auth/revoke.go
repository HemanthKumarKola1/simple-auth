package usecase

func (a *authUsecase) RevokeJwt(token string) {
	revokedJwts[token] = true
}
