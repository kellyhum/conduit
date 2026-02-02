package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/kellyhum/conduit/config"
	"github.com/kellyhum/conduit/shared"
)

const (
	StateFirstTimePrompt = iota
	StateFirstTimeUserCreation
	StateReturningUser
	StateListCommandsPrompt
)

type model struct {
	currState       int
	message         string
	username        string
	userInput       string
	showProgressBar bool
}

func InitialModel() model {
	return model{
		currState:       StateFirstTimePrompt,
		message:         shared.CheckFirstTime,
		username:        "",
		userInput:       "",
		showProgressBar: false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated_model := m

	keyMsg, ok := msg.(tea.KeyMsg)
	if ok {
		// global quitting mechanism
		if keyMsg.String() == "q" || keyMsg.Type == tea.KeyCtrlC {
			fmt.Println(shared.QuitProgram)
			return updated_model, tea.Quit
		}

		// handle regular user input for all the input based states
		if updated_model.currState == StateFirstTimePrompt ||
			updated_model.currState == StateFirstTimeUserCreation ||
			updated_model.currState == StateReturningUser ||
			updated_model.currState == StateListCommandsPrompt {
			updated_model = HandleUserInput(updated_model, keyMsg)
		}

		// handle user input
		if keyMsg.Type == tea.KeyEnter {
			input := m.userInput
			var nextState int
			var nextMessage string
			var nextUsername = ""

			switch m.currState {
			case StateFirstTimePrompt:
				switch input {
				case "y":
					nextState = StateFirstTimeUserCreation
					nextMessage = shared.FirstTimeUserPrompt
				case "n":
					nextState = StateReturningUser
					nextMessage = shared.ReturningUserPrompt
				default:
					nextState = StateFirstTimePrompt
					nextMessage = shared.InvalidUserInput + "\n" + shared.CheckFirstTime
				}

			case StateFirstTimeUserCreation:
				pubKeyStr, err := config.CreateFirstTimeUser(input)
				if err != nil {
					nextState = StateFirstTimeUserCreation
					nextMessage = err.Error()
				} else {
					nextState = StateListCommandsPrompt
					nextMessage = shared.SetupComplete + pubKeyStr + shared.CommandList
				}

			case StateReturningUser:
				userData, err := config.VerifyReturningUser(input)
				if err != nil {
					nextState = StateReturningUser
					nextMessage = err.Error()
				} else {
					nextState = StateListCommandsPrompt
					nextMessage = shared.WelcomeBack + userData.GetUsername() + shared.CommandList
					nextUsername = userData.GetUsername()
				}
			}

			updated_model = m.SetModelValues(nextState, nextMessage, nextUsername, "", false, true)
		}
	}

	return updated_model, nil
}

func (m model) View() string {
	s := shared.Header
	s += "\n"
	s += m.message

	if m.currState == StateFirstTimePrompt ||
		m.currState == StateFirstTimeUserCreation ||
		m.currState == StateReturningUser ||
		m.currState == StateListCommandsPrompt {
		s += "\n"
		s += "> "
		s += m.userInput
	}

	return s
}

func (m model) SetModelValues(currState int, message string, username string, userInput string, showProgressBar bool, shouldAppendMsg bool) model {
	newMessage := message

	if shouldAppendMsg {
		newMessage = m.message + "\n" + message
	}

	return model{
		currState:       currState,
		message:         newMessage,
		username:        username,
		userInput:       userInput,
		showProgressBar: showProgressBar,
	}
}

func HandleUserInput(m model, keyMsg tea.KeyMsg) model {
	switch keyMsg.Type {
	case tea.KeyRunes:
		m.userInput += string(keyMsg.Runes)
	case tea.KeyBackspace:
		if len(m.userInput) > 0 {
			m.userInput = m.userInput[:len(m.userInput)-1]
		}
	}
	return m
}
