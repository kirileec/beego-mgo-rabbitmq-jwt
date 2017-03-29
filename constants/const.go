package constants

var (
	Secret        = []byte("llinx.me")  //jwt Secret key
	SecureKey     = []byte("secureKey") //custom header for jump over jwt auth
	MainGoRoutine = "main"
)
