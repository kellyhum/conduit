package shared

const (
	Header              = "CONDUIT\n-----------------------\n(press q to quit)"
	CheckFirstTime      = "Are you a first time user? (y/n)"
	FirstTimeUserPrompt = "Enter the username that you would like: "
	ReturningUserPrompt = "Enter your existing username: "
	InvalidUserInput    = "Invalid input. Please try again."
	SetupComplete       = "Setup complete! Public key (hashed): :"
	WelcomeBack         = "Welcome back, "
	Error               = "Error occurred: "
	UsernameNotFound    = "Username not found."
	PrivateKeyNotFound  = "Private key not found."
	CommandList         = "\nAvailable commands: upload, download, help"
	NoKey               = "No key pair found. Setting up Conduit for first-time use..."
	UploadFile          = "Enter the file to upload (must be in the current directory): "
	UploadFileSuccess   = "File uploaded successfully! Your file ID is: "
	UploadFileError     = "Unable to upload file, please try again. File id: "
	QuitProgram         = "Goodbye! Thanks for using Conduit."
	ServerErrorResp     = "Server returned non-201 status: "
)
