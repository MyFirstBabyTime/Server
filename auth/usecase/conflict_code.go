package usecase

const (
	// use in authUsecase.SendCertifyCodeToPhone
	phoneAlreadyInUse = -101

	// use in authUsecase.CertifyPhoneWithCode
	phoneAlreadyCertified = -111
	incorrectCertifyCode = -112

	// use in authUsecase.SignUpParent
	uncertifiedPhone = -121
	parentIDAlreadyInUse = -122

	// use in authUsecase.LoginParentAuth
	notExistParentID = -131
	incorrectParentPW = -132
)
