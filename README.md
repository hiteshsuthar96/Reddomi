# 🧠 Reddomi

**Reddomi** is an open-source, AI-powered lead generation platform for **Reddit**.\
It automates lead discovery, monitoring, and engagement — helping you find and connect with high-intent users across Reddit while staying fully compliant with community rules.

> 🚀 Automate your Reddit lead generation with AI agents that discover relevant subreddits, monitor discussions, and craft authentic, rule-safe engagement messages.
> 
> This repo is now positioned as the MVP foundation for the Reddomi SaaS launch.

---

<img width="2862" height="1372" alt="image" src="https://github.com/user-attachments/assets/a9ac9b02-f29f-4e4b-b71d-059a7e57c198" />


## 📋 Table of Contents

- [Features](#-features)
- [Architecture Overview](#-architecture-overview)
- [Tech Stack](#-tech-stack)
- [Getting Started](#-getting-started)
    - [Configuration](#configuration)
    - [Backend Setup](#backend-setup)
    - [Frontend Setup](#frontend-setup)
    - [Running the Project](#running-the-project)
- [Integrations](#-integrations)
- [Deployment](#-deployment)
- [Account Safety Strategies](#-account-safety-strategies)
- [Gotchas & TODOs](#-gotchas--todos)
- [Contributing](#-contributing)
- [License](#-license)

---

## ✨ Features

### 🕵️ Lead Generation

- Extract relevant posts for given **keywords**, **subreddits**, and **product details**
- AI-assisted **keyword** and **subreddit** suggestions
- Intelligent post scoring using LLMs

### ⚙️ Automation

- AI-generated **DMs** and **comments** tailored to community guidelines
- **Scheduled auto-replies** on relevant posts
- **Auto DMs** with configurable timing and frequency

### 🔔 Notifications

- Email notifications via [Resend](https://resend.com)
- Slack alerts for daily or weekly summaries
- Email Alerts to users when a connected account is revoked or banned
- Onboarding emails
- Subscription expiry and renewal alerts

### 👥 Multi-Account Management

- Connect multiple Reddit accounts
- Auto-rotation for comments and DMs
- Rotation strategies:
    - `Random` — pick a random account
    - `Specific` — use a chosen account
    - `Most Qualified` — based on karma, age, etc.

### 📊 Reporting

- Daily and weekly engagement summaries

### 💳 Subscription Management

- Simple in-app subscription logic with plan limits

### 💬 Interactions

- Manage all AI-generated **comments** and **DMs**

### 🗓️ Posting

- Generate and schedule Reddit posts
- Posts follow subreddit rules and guidelines

---

## 🏗️ Architecture Overview

This repository contains all components required to run the Reddomi MVP.

## No Reddit API required for the MVP

The fastest path to a working test is the cookie-based Reddit DM login flow:

- Use the Cookie Integration screen in the portal to connect a Reddit account.
- This path is designed to work without official Reddit API credentials for DM automation.
- The API/OAuth path is optional and mainly useful for broader comment automation.

For the launch checklist and readiness notes, see [docs/MVP_READY.md](docs/MVP_READY.md).

```
.
├── backend/             # Go backend services
│   ├── portal-api/      # Public API layer for frontend (gRPC + Connect)
│   └── spooler/         # Core tracking engine for subreddits and posts
├── frontend/            # Frontend mono-repo
│   ├── portal/          # Web app (Next.js + PNPM)
│   └── packages/        # Shared UI, config, and protobuf packages
└── devel/               # Local development setup scripts
```

### Backend Services

- **Portal** — API layer for frontend over gRPC/Connect-Web
- **Spooler** — Tracks relevant posts based on subreddit/keyword pairs every 24h
    - Configurable fanout limits, polling intervals, and daily quotas
    - Limits no of posts to track
    - Limits Max no of relevant posts per day as per the subscribed plan
    - Built-in retry logic for failed posts
    - Send scheduled posts, DM and Comments

### LLM Layer

- Powered by **LiteLLM** (deployed separately)
- Supports **OpenAI** and **Gemini** APIs interchangeably
- Handles scoring, comment/DM generation, and keyword/subreddit suggestions

---

## 🧰 Tech Stack

**Backend**

- Go `1.23+`
- PostgreSQL
- Redis
- Docker
- LiteLLM
- Playwright

**Frontend**

- Node.js `20+`
- PNPM
- Next.js / React
- Tailwind CSS / Material UI

**Auth & APIs**

- Auth0 (passwordless login)
- Reddit OAuth APIs
- Resend (emails)
- Browserless / Steel.dev (CDP automation)
- DODO Payments (subscriptions)

---

## ⚙️ Getting Started

### Prerequisites

Ensure you have installed:

- Docker
- Go `1.23+`
- Node.js `20+`
- PNPM
- [direnv](https://direnv.net/) for environment variables

---

### Configuration

Start local PostgreSQL and Redis:

```bash
./devel/up.sh
```

Copy the environment file and configure:

```bash
cp .envrc.example .envrc
direnv allow
```

Replace placeholders (`<value>`) with your actual secrets and keys.

---

### Backend Setup

Run tests and start the backend:

```bash
cd backend
go test ./...
go build -o redora && ./redora start
```

Initialize the database:

```bash
./backend/script/migrate.sh up
```

Create a new migration:

```bash
./backend/script/migrate.sh new <migration_name>
```

---

### Frontend Setup

Install dependencies:

```bash
cd frontend
pnpm install
```

Start the development server:

```bash
pnpm dev:portal
```

Visit: [http://localhost:3000](http://localhost:3000)

## 🚀 MVP launch checklist

The current repo is ready for a lean MVP launch path when the following are in place:

1. Valid Reddit integration credentials or browser/cookie automation flow
2. Working OpenAI/LiteLLM key for DM/comment generation
3. Auth0 + portal config for real user sign-in
4. Billing/trial flow stubbed for the first customer demo
5. Clean Reddomi branding and onboarding copy

Use this repo as the launch foundation for Reddomi, then upgrade only the features that customers actually ask for.

---

### Running the Project

You’ll need three components running:

1. **Docker** — for Postgres, Redis, Pub/Sub emulator
   ```bash
   ./devel/up.sh
   ```
2. **Backend** — use [reflex](https://github.com/cespare/reflex) for live reload
   ```bash
   reflex -c .reflex
   ```
3. **Frontend** — Next.js app
   ```bash
   cd frontend && pnpm dev:portal
   ```

Visit:

- `http://localhost:8081` - pgweb (Postgres UI)
- `http://localhost:3000` - Reddomi Portal

---

## Exact working startup (tested)

1. Copy example environment values:

```bash
cp .envrc.example .envrc
```

2. Load environment variables:

```bash
direnv allow
```

3. Start local dependencies:

```bash
./devel/up.sh
```

4. Start the app using the built-in startup script:

```bash
./run-local.sh
```

This script does the following:

- sources `.envrc` if present
- exports `NEXT_PUBLIC_API_URL` and `NEXT_PUBLIC_APP_URL` defaults
- starts Docker services for Postgres and Redis
- builds the backend and starts it on port `8787`
- starts the frontend on port `3000`

If the frontend or backend fails to start:

- make sure `docker`, `go`, and `pnpm` are installed
- make sure port `3000` is free
- make sure `.envrc` contains valid Auth0 and Reddit secrets for the portal login flow

### Codespaces-specific notes

If you open this repo in a new GitHub Codespace, do this before visiting the app:

```bash
gh codespace ports visibility 3000:public 8787:public
```

Then open the public Codespaces app URL for port `3000`, for example:

```text
https://<your-codespace-id>-3000.app.github.dev/auth/login
```

The frontend already supports runtime fallback for the Codespaces API URL in:

- `frontend/portal/src/services/grpc.ts`
- `frontend/portal/src/services/config.ts`

So if ports `3000` and `8787` are public, the frontend will resolve the backend automatically.

### If you want to avoid repeating the same errors

Always start from the repo root and use:

```bash
./run-local.sh
```

If you need the exact file set for another AI or a new collaborator, point them here first:

- `.envrc.example` — required environment variables and local defaults
- `run-local.sh` — the full local startup flow
- `frontend/portal/src/services/grpc.ts` — API endpoint resolution logic
- `frontend/portal/src/services/config.ts` — portal and API URL configuration
- `frontend/portal/src/components/pages/Login.tsx` — OTP login UI and flow
- `backend/portal/handler_portal_passwordless_start.go` — passwordless start handler
- `backend/portal/handler_portal_passwordless_verify.go` — passwordless verify handler

If a new AI is taking over the repo, tell it:

- "Run `./run-local.sh` from the repo root after copying `.envrc.example` to `.envrc` and setting secrets."
- "In Codespaces, expose ports `3000` and `8787` publicly."
- "Check `frontend/portal/src/services/grpc.ts` for runtime API resolution."
- "Check `backend/portal/handler_portal_passwordless_start.go` and `...verify.go` for the OTP login route."

---

## �🔌 Integrations

Integrations store external service credentials and configuration.

| Type               | Description                                 |
| ------------------ | ------------------------------------------- |
| **Reddit Cookies** | User-provided cookies for Reddit automation |
| **Slack Webhook**  | Notifications and alerts                    |
| **OAuth Tokens**   | Reddit access/refresh tokens                |

**Manually insert an integration using tools. Example (CLI):**

```bash
doota tools integrations slack_webhook create <org-id> '{"channel":"reddomi-alerts","webhook":"<slack-url>"}'
```
---

## 🛠️ Admin Interface

There is no separate admin interface. Users assigned the role `PLATFORM_ADMIN` can view all organizations and have access to all accounts across the platform.


## ☁️ Deployment

- Hosted on **Railway.app**
- GCP used for LLM and Playwright storage
- Secrets managed via **GCP KMS**
- PostgreSQL and Redis hosted on Railway

---

## 🛡️ Account Safety Strategies

To minimize Reddit account bans:

1. **Age-based limits** — automate only with accounts > 2 weeks old
2. **Gradual scaling** — increase activity slowly
3. **DM-first approach** — DMs are safer than comments
4. **Rule adherence** — generate context-aware replies
5. **Consistent engagement** — reply and post regularly
6. **Occasional posting** — maintain activity score

---

## ⚠️ Gotchas & TODOs

- We currently use OpenAI for keyword/subreddit suggestions during onboarding and LiteLLM for other AI related tasks. Ideally, we should use a single AI provider for all tasks.
- When scoring a post, if the score is >90, we double check it with an advance model. This could be improved by selecting a better default model.
- Comment and DM generation should be moved into separate LLM calls. Right now, scoring, comment and DM generation are all done in a single LLM call.
- We should add the ability to regenerate comments or DMs.
- To avoid getting banned, we only use Reddit accounts that are > 2 weeks old for AI generation. This is a temporary solution and we should come up with a better way to handle account warmup.


---

## 🤝 Contributing

We welcome contributions!

1. Fork the repo
2. Create a new branch:
   ```bash
   git checkout -b feature/your-feature
   ```
3. Commit your changes
4. Open a Pull Request

Please follow Go and JS linting rules before submitting.

---

## 📜 License

Released under the [MIT License](LICENSE).\
© 2025 DoneByAI — building AI tools that work *for* you.

---

### 💡 Maintainers

- [DoneByAI Team](https://donebyai.team)
- [Shashank Agarwal](https://github.com/shank318)
