package telegram

import (
	"SplitCore/internal/repository"
	"log/slog"
	"sync"
)

type BotHandler struct {
	userState map[int64]*UserContext
	userRepo  repository.UserRepository
	fundRepo  repository.FundRepository
	mu        sync.RWMutex
}

type State int

type UserContext struct {
	State     State
	LastMsgID int
}

type SendMode int

const (
	Edit SendMode = iota
	Reply
	Send
)

const (
	StateNone State = iota
	StateFundName
	StateFundJoinCode
	StateFundMenu
	StateViewFund
)
const (
	CommandCreateFund = "create_fund"
	CommandMyFund     = "my_fund"
	CommandJoinFund   = "join_fund"
	CommandBack       = "back"
	CommandNext       = "next"
	CommandPrevious   = "previous"
	CommandFund       = "view_fund"
)

func NewBotHandler(userRepository repository.UserRepository, fundRepository repository.FundRepository) *BotHandler {
	slog.Info("Setting up telegram bot")
	return &BotHandler{
		userState: make(map[int64]*UserContext),
		userRepo:  userRepository,
		fundRepo:  fundRepository,
	}
}
