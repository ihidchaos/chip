package secure_channel

const (
	SInitialized       = 0
	SSentSigma1        = 1
	SSentSigma2        = 2
	SSentSigma3        = 3
	SSentSigma1Resume  = 4
	kSentSigma2Resume  = 5
	kFinished          = 6
	kFinishedViaResume = 7

	kSigmaParamRandomNumberSize = 32
)
