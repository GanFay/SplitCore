package telegram

import (
	"SplitCore/internal/domain"
	"SplitCore/internal/pkg/utils"
	"context"
	"fmt"
	"log/slog"
	"os"

	tele "gopkg.in/telebot.v4"
)

func (h *BotHandler) HandleStart(c tele.Context) error {
	if h.userState[c.Sender().ID] == nil {
		h.userState[c.Sender().ID] = &UserContext{State: StateNone}
	}
	ctx := context.Background()
	var user domain.User
	user.TgID = c.Sender().ID
	user.Username = c.Sender().Username
	user.FirstName = c.Sender().FirstName

	_, err := h.userRepo.Create(ctx, &user)
	if err != nil {
		slog.Warn("could not register user", "err", err, "id", user.TgID)
	}
	args := c.Args()

	// if url invite code
	if len(args) == 1 {
		arg := args[0]
		if len(arg) != 6 {
			return c.Send("⚠️ Invalid invite link format.")
		}
		fund := &domain.Fund{
			InviteCode: arg,
		}
		fund, err = h.fundRepo.GetInfo(ctx, fund)
		if err != nil {
			return h.Error(c, "Invite code not found", err.Error(), Reply)
		}

		err = h.fundRepo.AddMember(ctx, fund, user.TgID)
		if err != nil {
			return h.Error(c, "Failed to join the fund", err.Error(), Reply)
		}

		msg := "Congratulations🎉\n\nYou have successfully joined to the fund😊!\n" +
			"You can see them in <b>My Funds</b>⬇️"
		return c.Reply(msg, h.MainMenu(), tele.ModeHTML)
	}
	// if url invite code
	msg := "👋 Hello! I'm your expense helper. Use the menu below:"
	return c.Reply(msg, h.MainMenu(), tele.ModeHTML)
}

func (h *BotHandler) HandleBack(c tele.Context) error {
	id := c.Sender().ID
	if h.userState[id] == nil {
		h.userState[id] = &UserContext{
			State: StateNone,
		}
	}
	h.userState[id].State = StateNone
	msg := "👋 Hello! I'm your expense helper. Use the menu below👇"
	return c.Edit(msg, h.MainMenu(), tele.ModeHTML)
}

func (h *BotHandler) OnText(c tele.Context) error {
	id := c.Sender().ID
	err := c.Delete()
	if err != nil {
		slog.Error("error delete message", "id", id, "err", err.Error())
		return err
	}
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
		slog.Info("Setting up fund", "fund", fund, "id", id)
		_, err = h.fundRepo.Create(ctx, &fund)
		if err != nil {
			return h.Error(c, "Failed to create fund", err.Error(), Edit)
		}

		storedMsg := &tele.Message{ID: h.userState[id].LastMsgID, Chat: c.Chat()}
		msg := fmt.Sprintf("Fund Created🎉!\n\nFund Code🔑: <code>%s</code>\nFund URL🌐: <code>%s</code>\n\n You can share URL or Code with users your fund👍", fund.InviteCode, InviteCodeInviteURL)
		ctxMsg, err := c.Bot().Edit(storedMsg, msg, h.BackMenu(), tele.ModeHTML)
		if err != nil {
			return h.Error(c, "Failed to edit fund", err.Error(), Edit)
		}
		h.userState[id].LastMsgID = ctxMsg.ID
	case StateFundJoinCode:
		fund := &domain.Fund{
			InviteCode: text,
		}
		fund, err = h.fundRepo.GetInfo(ctx, fund)
		if err != nil {
			return h.Error(c, "Failed to get fund", err.Error(), Edit)
		}

		err = h.fundRepo.AddMember(ctx, fund, id)
		if err != nil {
			return h.Error(c, "Failed to join the fund", err.Error(), Edit)
		}

		storedMsg := &tele.Message{ID: h.userState[id].LastMsgID, Chat: c.Chat()}
		msg := "You successfully joined to fund🎉\n\n" +
			"Go to <b>My Fund</b> to see this⬇️."
		ctxMsg, err := c.Bot().Edit(storedMsg, msg, h.BackMenu(), tele.ModeHTML)

		if err != nil {
			return h.Error(c, "Failed to edit fund", err.Error(), Edit)
		}
		h.userState[id].LastMsgID = ctxMsg.ID
		slog.Info("Setting up fund join code", "id", id)
	case StateNone:
		storedMsg := &tele.Message{ID: h.userState[id].LastMsgID, Chat: c.Chat()}
		msg := "No answer"
		_, err = c.Bot().Edit(storedMsg, msg, h.BackMenu(), tele.ModeHTML)
		if err != nil {
			slog.Error("error to edit message", "id", id, "err", err.Error())
			return err
		}
	default:
		panic("You have unstatement case")
	}
	h.userState[id].State = StateNone
	return nil
}

func (h *BotHandler) Error(c tele.Context, userMsg string, techMsg string, mode SendMode) error {
	slog.Error(userMsg, "err", techMsg, "user_id", c.Sender().ID)
	storedMsg := &tele.Message{ID: h.userState[c.Sender().ID].LastMsgID, Chat: c.Chat()}
	switch mode {
	case Edit:
		_, err := c.Bot().Edit(storedMsg, "❌"+userMsg, h.BackMenu(), tele.ModeHTML)
		return err
	case Reply:
		return c.Reply("❌"+userMsg, h.BackMenu(), tele.ModeHTML)
	case Send:
		return c.Send("❌"+userMsg, h.BackMenu(), tele.ModeHTML)
	}
	return fmt.Errorf("send error")
}
