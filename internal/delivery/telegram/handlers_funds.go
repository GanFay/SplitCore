package telegram

import (
	"SplitCore/internal/domain"
	"context"
	"fmt"
	"log/slog"
	"strconv"

	tele "gopkg.in/telebot.v4"
)

func (h *BotHandler) HandleCreateFund(c tele.Context) error {
	defer c.Respond()

	id := c.Sender().ID
	h.mu.Lock()
	if h.userState[id] == nil {
		h.userState[id] = &UserContext{}
	}
	h.userState[id].State = StateFundName
	h.userState[id].LastMsgID = c.Message().ID
	h.mu.Unlock()
	msg := "Enter the desired fund name (any name)👇"
	return c.Edit(msg, h.BackMenu(), tele.ModeHTML)
}

func (h *BotHandler) HandleMyFund(c tele.Context) error {
	defer c.Respond()

	msg := "Your funds👇"
	h.mu.Lock()
	h.userState[c.Sender().ID] = &UserContext{
		State:     StateFundMenu,
		LastMsgID: c.Message().ID,
	}
	h.mu.Unlock()
	return c.Edit(msg, h.MyFundMenu(c, 0), tele.ModeHTML)
}

func (h *BotHandler) HandleJoinFund(c tele.Context) error {
	defer c.Respond()

	id := c.Sender().ID
	h.mu.Lock()
	if h.userState[id] == nil {
		h.userState[id] = &UserContext{
			State: StateFundJoinCode,
		}
	}
	h.userState[id].State = StateFundJoinCode
	h.userState[id].LastMsgID = c.Message().ID
	h.mu.Unlock()
	msg := "Input Join Code🔑:\n\n" +
		"You can get an invite code by asking the fund's creator🧍‍♂️\nOr create one yourself"
	return c.Edit(msg, h.BackMenu(), tele.ModeHTML)
}

func (h *BotHandler) HandleNext(c tele.Context) error {
	defer c.Respond()

	offset, err := strconv.Atoi(c.Data())
	h.mu.Lock()

	h.userState[c.Sender().ID] = &UserContext{
		LastMsgID: c.Message().ID,
	}
	h.mu.Unlock()
	if err != nil {
		return h.Error(c, "Internal Error, try again later", err.Error(), Edit)
	}
	return c.Edit(h.MyFundMenu(c, offset), tele.ModeHTML)
}

func (h *BotHandler) HandlePrevious(c tele.Context) error {
	defer c.Respond()

	offset, err := strconv.Atoi(c.Data())
	h.mu.Lock()
	h.userState[c.Sender().ID] = &UserContext{
		LastMsgID: c.Message().ID,
	}
	h.mu.Unlock()
	if err != nil {
		return h.Error(c, "Internal Error, try again later", err.Error(), Edit)
	}
	return c.Edit(h.MyFundMenu(c, offset), tele.ModeHTML)
}

func (h *BotHandler) HandleViewFund(c tele.Context) error {
	defer c.Respond()

	ctx := context.Background()
	h.mu.Lock()
	h.userState[c.Sender().ID] = &UserContext{
		State:     StateViewFund,
		LastMsgID: c.Message().ID,
	}
	h.mu.Unlock()
	fundId, err := strconv.Atoi(c.Data())
	if err != nil {

		return h.Error(c, "Internal Error, try again later", err.Error(), Edit)
	}
	fund := &domain.Fund{
		ID: fundId,
	}
	slog.Debug("", "id", fundId)

	fund, err = h.fundRepo.GetInfo(ctx, fund)
	if err != nil {
		return h.Error(c, "Internal Error, failed to get info about this fund, try again later", err.Error(), Edit)
	}
	msg := fmt.Sprintf("Your fund⬇️:\n\nFundName: <code>%s</code>\nAuthorID: <code>%d</code>\nCreatedAt: <code>%s</code>\nInviteCode: <code>%s</code>", fund.Name, fund.AuthorID, fund.CreatedAt.Format(`02.01.2006 15:04`), fund.InviteCode)
	return c.Edit(msg, h.BackMenu(), tele.ModeHTML)
}
