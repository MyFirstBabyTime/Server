package domain

const (
	// use in authUsecase.SendCertifyCodeToPhone
	PhoneAlreadyInUse = -101

	// use in authUsecase.CertifyPhoneWithCode
	PhoneAlreadyCertified = -111
	IncorrectCertifyCode  = -112

	// use in authUsecase.SignUpParent
	UncertifiedPhone     = -121
	ParentIDAlreadyInUse = -122

	// use in authUsecase.LoginParentAuth
	NotExistParentID  = -131
	IncorrectParentPW = -132
)
