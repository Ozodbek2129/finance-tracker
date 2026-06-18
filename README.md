# 💰 Finance Tracker

**Finance Tracker** — foydalanuvchilarga daromad va xarajatlarini kuzatib borish imkonini beruvchi to'liq loyiha. Backend (Golang), Frontend (React) va Telegram Bot (Python) dan iborat.

🌐 Demo

| Qism             |                     Link                     |
|------------------|----------------------------------------------|
| 🖥 Web sayt      | https://finance-tracker-five-neon.vercel.app |
| 🤖 Telegram Bot | @finance_tracker_555_bot                     |

## 🧩 Loyiha tuzilmasi

SAVAT/
├── backend/                  # Golang REST API
│   ├── cmd/
│   │   └── main.go
│   ├── handler/
│   │   ├── handler.go
│   │   └── router.go
│   ├── internal/
│   │   ├── config/
│   │   ├── connectiondb/
│   │   │   ├── postgres_crud/
│   │   │   │   ├── stats.go
│   │   │   │   └── transaction.go
│   │   │   └── postress.go
│   │   ├── logger/
│   │   ├── migrations/
│   │   ├── model/
│   │   └── middleware/
│   ├── .env
│   ├── .gitignore
│   ├── app.log
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
│
├── frontend/finance-tracker/  # React frontend
│   ├── public/
│   ├── src/
│   │   ├── api/
│   │   │   └── client.js
│   │   ├── assets/
│   │   ├── features/
│   │   │   ├── categories/
│   │   │   │   └── CategoryManager.jsx
│   │   │   ├── stats/
│   │   │   │   └── DashboardStats.jsx
│   │   │   └── transactions/
│   │   │       ├── TransactionForm.jsx
│   │   │       └── TransactionList.jsx
│   │   ├── App.css
│   │   ├── App.jsx
│   │   ├── index.css
│   │   └── main.jsx
│   ├── index.html
│   ├── package.json
│   ├── vite.config.js
│   └── eslint.config.js
│
└── telegram_bot/              # Python Telegram Bot + Mini App
    ├── bot.py
    ├── index.html             # Mini App UI
    └── requirements.txt

## 🛠 Texnologiyalar

| Qism     | Texnologiya              |
|----------|--------------------------|
| Backend  | Go, PostgreSQL, Docker   |
| Frontend | React, Vite              |
| Telegram | Python, Telegram Bot API |

## ⚙️ O'rnatish

### 1. Repozitoriyani klonlash

```bash
git clone https://github.com/Ozodbek2129/finance-tracker.git
cd finance-tracker
```

### 🔵 Backend (Golang)

```bash
cd backend
```

`.env` faylini sozlang:

```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=finance_tracker
```

**Docker orqali ishga tushirish:**

```bash
docker build -t finance-backend .
docker run --env-file .env -p 8080:8080 finance-backend
```

**Lokal ishga tushirish:**

```bash
go mod tidy
go run cmd/main.go
```

---

### 🟢 Frontend (React)

```bash
cd frontend/finance-tracker
```

**O'rnatish va ishga tushirish:**

```bash
npm install
npm run dev
```

Brauzerda oching: [http://localhost:5173](http://localhost:5173)

**Production build:**

```bash
npm run build
```

---

### 🟡 Telegram Bot (Python)

```bash
cd telegram_bot
```

`.env` yoki `bot.py` ichida sozlang:

```env
BOT_TOKEN=your_telegram_bot_token
API_URL=http://localhost:8080/api
```

**Paketlarni o'rnatish:**

```bash
pip install -r requirements.txt
```

**Botni ishga tushirish:**

```bash
python bot.py
```

---

## 🔌 API Endpointlar

### Tranzaksiyalar

| Metod  |           URL              |        Tavsif         |
|--------|----------------------------|-----------------------|
| GET    | `/api/v1/transactions`     | Barcha tranzaksiyalar |
| POST   | `/api/v1/transactions`     | Yangi tranzaksiya     |
| PUT    | `/api/v1/transactions/:id` | Yangilash             |
| DELETE | `/api/v1/transactions/:id` | O'chirish             |

### Statistika

| Metod |         URL          |      Tavsif       |
|-------|----------------------|-------------------|
| GET   | `/api/v1/stats`      | Umumiy statistika |

### Kategoriyalar

| Metod  |            URL            |      Tavsif          |
|--------|---------------------------|----------------------|
| GET    | `/api/v1/categories`      | Barcha kategoriyalar |
| POST   | `/api/v1/categories`      | Yangi kategoriya     |
| PUT    | `/api/v1/categories/:id`  | Yangilash            |
| DELETE | `/api/v1/categories/:id`  | O'chirish            |

## 📱 Telegram Mini App

1. [@BotFather](https://t.me/BotFather) ga o'ting
2. `/mybots` → botingizni tanlang
3. `Bot Settings` → `Menu Button` → `index.html` hosting URL ni kiriting


## 📝 Litsenziya

MIT License