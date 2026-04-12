package telegram

import (
	"SplitCore/internal/domain"
	"SplitCore/internal/pkg/utils"
	"SplitCore/internal/repository"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	tele "gopkg.in/telebot.v4"
)

type BotHandler struct {
	userState map[int64]*UserContext
	userRepo  repository.UserRepository
	fundRepo  repository.FundRepository
}

type State int

type UserContext struct {
	State     State
	LastMsgID int
}

const (
	StateNone State = iota
	StateFundName
	StateFundJoinCode
)
const (
	CommandCreateFund = "create_fund"
	CommandMyFund     = "my_fund"
	CommandJoinFund   = "join_fund"
	CommandBack       = "back"
)

func NewBotHandler(userRepository repository.UserRepository, fundRepository repository.FundRepository) *BotHandler {
	slog.Info("Setting up telegram bot")
	return &BotHandler{
		userState: make(map[int64]*UserContext),
		userRepo:  userRepository,
		fundRepo:  fundRepository,
	}
}

//--------------------Menu-----------------------

func (h *BotHandler) MainMenu() *tele.ReplyMarkup {
	menu := tele.ReplyMarkup{ResizeKeyboard: true}

	btnCreateFund := menu.Data("Create Fund", CommandCreateFund)
	btnMyFund := menu.Data("My Fund", CommandMyFund)
	btnJoinFund := menu.Data("Join Fund", CommandJoinFund)

	menu.Inline(
		menu.Row(btnCreateFund),
		menu.Row(btnMyFund),
		menu.Row(btnJoinFund),
	)
	return &menu
}

func (h *BotHandler) BackMenu() *tele.ReplyMarkup {
	menu := tele.ReplyMarkup{ResizeKeyboard: true}
	btnBack := menu.Data("Back", CommandBack)
	menu.Inline(menu.Row(btnBack))
	return &menu
}

//--------------Router--------------

func (h *BotHandler) SetupRegister(b *tele.Bot) {
	b.Use(LoggingMiddleware())
	b.Handle("/start", h.HandleStart)
	b.Handle("\f"+CommandCreateFund, h.HandleCreateFund)
	b.Handle("\f"+CommandMyFund, h.HandleMyFund)
	b.Handle("\f"+CommandJoinFund, h.HandleJoinFund)
	b.Handle("\f"+CommandBack, h.HandleBack)
	b.Handle(tele.OnText, h.OnText)
	slog.Info("Setting up handlers")
}

//-------------Handlers-----------

func (h *BotHandler) HandleStart(c tele.Context) error {
	ctx := context.Background()
	var user domain.User
	user.TgID = c.Sender().ID
	user.Username = c.Sender().Username
	user.FirstName = c.Sender().FirstName

	_, err := h.userRepo.Create(ctx, &user)
	if err != nil {
		slog.Warn("START", "error to add user to DB", err.Error())
	}
	args := c.Args()
	slog.Debug("START", "arg", args, "len", len(args))

	// if url invite code
	if len(args) == 1 {
		arg := args[0]
		if len(arg) != 6 {
			err = errors.New("invalid argument")
			slog.Warn("startInvalidArg", "error", err.Error())
			return err
		}
		fund, err := h.fundRepo.GetByInviteCode(ctx, arg)
		if err != nil {
			slog.Warn("startGet", "error to get by ic", err.Error())
			return err
		}
		slog.Debug("START", "fund", fund)
		err = h.fundRepo.AddMember(ctx, fund, user.TgID)
		if err != nil {
			slog.Warn("startAddMEMBER", "error to add member", err.Error())
			return err
		}
		slog.Debug("START", "err", err)
		msg := "Congratulations🎉\n\nYou have successfully joined to the fund😊!\n" +
			"You can see them in <b>My Funds</b>⬇️"
		return c.Reply(msg, h.MainMenu(), tele.ModeHTML)
	}
	// if url invite code
	msg := "Hello, it's helper:"
	return c.Reply(msg, h.MainMenu(), tele.ModeHTML)
}

func (h *BotHandler) HandleCreateFund(c tele.Context) error {
	id := c.Sender().ID
	if h.userState[id] == nil {
		h.userState[id] = &UserContext{}
	}

	h.userState[id].State = StateFundName
	h.userState[id].LastMsgID = c.Message().ID
	msg := "Input Name Fund:"
	return c.Edit(msg, h.BackMenu(), tele.ModeHTML)
}

func (h *BotHandler) HandleBack(c tele.Context) error {
	id := c.Sender().ID
	if h.userState[id] == nil {
		h.userState[id] = &UserContext{
			State: StateNone,
		}
	}
	h.userState[id].State = StateNone
	msg := "Menu:"
	return c.Edit(msg, h.MainMenu(), tele.ModeHTML)
}

func (h *BotHandler) HandleMyFund(c tele.Context) error {
	ctx := context.Background()
	id := c.Sender().ID

	offset := "0"
	limit := "5"
	_, err := h.fundRepo.GetByUserID(ctx, id, limit, offset)
	if err != nil {
		slog.Warn("MYFUND", "error to find fund by userID", err.Error())
		return err
	}

	msg := "Your funds:"
	return c.Edit(msg, h.BackMenu(), tele.ModeHTML)
}

func (h *BotHandler) HandleJoinFund(c tele.Context) error {
	id := c.Sender().ID
	if h.userState[id] == nil {
		h.userState[id] = &UserContext{
			State: StateFundJoinCode,
		}
	}
	h.userState[id].State = StateFundJoinCode
	h.userState[id].LastMsgID = c.Message().ID
	msg := "Input Join Code:"
	return c.Edit(msg, h.BackMenu(), tele.ModeHTML)
}

func (h *BotHandler) OnText(c tele.Context) error {
	err := c.Delete()
	if err != nil {
		slog.Warn("ONTXT", "error delete message", err.Error())
		return err
	}
	id := c.Sender().ID
	if h.userState[id] == nil {
		return nil
	}
	text := c.Text()
	ctx := context.Background()
	switch h.userState[id].State {
	case StateFundName:
		InviteCode := utils.GenerateInviteCode(6)
		botName := os.Getenv("BOT_NAME")
		InviteCodeInviteURL := utils.GenerateInviteCodeURL(InviteCode, botName)
		fund := domain.Fund{
			AuthorID:   id,
			Name:       text,
			InviteCode: InviteCode,
		}
		slog.Info("Setting up fund", "fund", fund)
		_, err = h.fundRepo.Create(ctx, &fund)
		if err != nil {
			slog.Warn("STATEFN", "error to create fund", err.Error())
			return err
		}

		storedMsg := &tele.Message{ID: h.userState[id].LastMsgID, Chat: c.Chat()}
		msg := fmt.Sprintf("Fund Created🎉!\n\n Fund Code🔑: <code>%s</code>\nFund URL🌐: <code>%s</code>\n\n You can share URL or Code with users your fund👍", fund.InviteCode, InviteCodeInviteURL)
		ctxMsg, err := c.Bot().Edit(storedMsg, msg, h.BackMenu(), tele.ModeHTML)
		if err != nil {
			slog.Warn("STATEFN", "error to edit message", err.Error())
			return err
		}
		h.userState[id].LastMsgID = ctxMsg.ID
	case StateFundJoinCode:
		fund, err := h.fundRepo.GetByInviteCode(ctx, text)
		if err != nil {
			slog.Warn("STATEJS", "error to get fund by inviteCode", err.Error(), "inviteCode", text)
			return err
		}

		err = h.fundRepo.AddMember(ctx, fund, id)
		if err != nil {
			slog.Warn("STATEJS", "error to add member", err.Error())
			return err
		}

		storedMsg := &tele.Message{ID: h.userState[id].LastMsgID, Chat: c.Chat()}
		msg := "You successfully joined to fund🎉\n\n" +
			"Go to <b>My Fund</b> to see this⬇️."
		ctxMsg, err := c.Bot().Edit(storedMsg, msg, h.BackMenu(), tele.ModeHTML)

		if err != nil {
			slog.Warn("STATEJC", "error to edit message", err.Error())
			return err
		}
		h.userState[id].LastMsgID = ctxMsg.ID
		slog.Info("Setting up fund join code")
	case StateNone:
		storedMsg := &tele.Message{ID: h.userState[id].LastMsgID, Chat: c.Chat()}
		msg := "No answer"
		_, err = c.Bot().Edit(storedMsg, msg, h.BackMenu(), tele.ModeHTML)
		if err != nil {
			slog.Warn("STATENONE", "error to edit message", err.Error())
			return err
		}
	}
	h.userState[id].State = StateNone
	return nil
}
