# 🚀 SplitCore — Telegram Expense Organizer

**SplitCore** is a Telegram bot designed to automate shared expense tracking for groups of friends, travelers, or event organizers. No more messy Excel sheets or "who owes whom" arguments.

## 🔥 The Core Idea
Users create "Funds" (events), invite friends via unique deep-links, and record their expenses. The bot automatically calculates the balance: who overpaid and who needs to settle their debt.

## 🛠 Tech Stack
* **Language:** Go (Golang) 1.22+
* **Framework:** [telebot.v4](https://github.com/tucnak/telebot) (Telegram Bot API)
* **Database:** PostgreSQL
* **Driver:** [pgx/v5](https://github.com/jackc/pgx) (Connection Pool)
* **In-memory storage:** FSM (Finite State Machine) for user dialog management.
* **Infrastructure:** Docker, Docker Compose, Makefile.
* **Migrations:** [golang-migrate](https://github.com/golang-migrate/migrate).

## 🏗 Architecture (Clean Architecture)
The project is built with a strict separation of concerns, ensuring high testability and scalability:
- `cmd/bot/` — Entry point, initialization, and Dependency Injection.
- `internal/domain/` — Business entities and Repository interfaces.
- `internal/repository/` — Database access layer (PostgreSQL implementation).
- `internal/delivery/telegram/` — Bot-specific logic (handlers, middleware, router).
- `internal/pkg/` — Internal utilities (random generators, etc.).

## 📍 Roadmap (MVP Status)
- [ ] +- (**In work**) Bot skeleton and Inline-button navigation.
- [x] Finite State Machine (FSM) for user input handling.
- [x] Database: migrations, user, and fund storage.
- [x] Fund creation and unique Deep-Link generation.
- [x] Automatic "Creator-to-Member" enrollment (Transactions).
- [ ] Expense logging (Purchases) — **In Progress**.
- [ ] Debt calculation math module — **Next Step**.

## 🚀 Getting Started (Dev)

1. Clone the repository:
   ```bash
   git clone https://github.com/GanFay/new-project.git
   ```
2. Set up environment variables in a `.env` file (Bot token, DB credentials).
3. Start the database using Docker:
   ```bash
   docker-compose up -d
   ```
4. Run migrations:
   ```bash
   make migrate-up
   ```
5. Run the bot:
   ```bash
   go run cmd/bot/main.go
   ```
