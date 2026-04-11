package telegram

import (
	"log/slog"

	tele "gopkg.in/telebot.v4"
)

func LoggingMiddleware() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			user := c.Sender()
			text := c.Text()
			if c.Callback() != nil {
				text = "callback: " + c.Callback().Data
			}

			slog.Info("incoming update",
				"user_id", user.ID,
				"username", user.Username,
				"data", text,
			)

			return next(c)
		}
	}
}
