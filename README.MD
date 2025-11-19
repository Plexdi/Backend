
# **Plexdi Studio — Backend**

A lightweight, fast, and scalable backend powering the **Plexdi Studio** website.
Built using **Go**, **Gin**, and a **PostgreSQL (Neon)** database, with clean routes, structured services, and email automation.

This backend handles:

* Commission submissions
* Admin dashboard data
* Status updates
* Email notifications (client + studio)
* Basic validation
* JSON persistence fallback (local mode)

---

## 🚀 **Tech Stack**

**Language & Framework**

* Go 1.22+
* Gin Web Framework
* PostgreSQL (Neon serverless)
* pgx driver

**Email Provider**

* Resend API (transactional emails)

**Infrastructure**

* Render (deployment)
* JSON file fallback (`data/commissions.json`)

---

## 🔧 **Core Features**

### 📩 **Commission System**

Handles user submissions from `/commissions` on the frontend.

* **Create Commission**

  * Saves to database
  * Sends client confirmation email
  * Sends admin notification email

* **Get Commissions**
  Returns all stored commissions for the admin dashboard.

* **Update Commission**

  * Allows updating status: `queued`, `in_progress`, `completed`
  * Automatically emails client when:

    * Work begins
    * Work completes (includes revision/purchase info)

* **Delete Commission**
  Removes entry from database.

---

### ✉️ **Email Automation**

When a client submits a commission:

✔️ Client receives a confirmation
✔️ You (Plexdi) receive an admin alert
✔️ When you update status → client gets progress emails

Emails sent using **Resend API** with HTML templates.

---

### 🔐 **Admin Integration (Frontend)**

Backend provides data to `/admins` dashboard:

* Fetch list of commissions
* Patch status
* Delete commission

Frontend handles authentication via token, backend stays stateless & simple.

---

## 📁 **Project Structure**

```
cmd/
 └─ api/
     └─ main.go               # Server entrypoint

internal/
 ├─ handlers/                # HTTP handlers
 ├─ services/                # Business logic (email, commission creation, DB ops)
 ├─ db/                      # PostgreSQL connection
 └─ models/                  # (Optional) structs

data/
 └─ commissions.json         # Local store fallback
```

---

## 📌 **API Routes**

### **POST /commissions**

Create a new commission.

```json
{
  "name": "Mogi",
  "email": "test@example.com",
  "discord": "mogi#1234",
  "type": "banner",
  "details": "Blue theme, anime style",
  "refs": "links…"
}
```

### **GET /commissions**

Returns list of all commissions.

### **PATCH /commissions/:id**

Update a commission's status.

Body example:

```json
{ "status": "in_progress" }
```

### **DELETE /commissions/:id**

Delete a commission from the database.

---

## 🛠️ **Environment Variables**

Create a `.env` file:

```
DATABASE_URL=postgres://...
PORT=10000
RESEND_API_KEY=your_resend_key
```

---

## ▶️ **Running Locally**

### Install dependencies

```bash
go mod tidy
```

### Run server

```bash
go run cmd/api/main.go
```

Default runs at:

```
http://localhost:10000
```

---

## 🧪 **Testing via Postman**

### Create Commission

```
POST http://localhost:10000/commissions
Content-Type: application/json
```

### Update Status

```
PATCH http://localhost:10000/commissions/7
Body:
{ "status": "in_progress" }
```

### Delete Commission

```
DELETE http://localhost:10000/commissions/7
```

---

## 🌐 **Deployment**

Currently set up for:

* **Render** (Go server)
* Neon PostgreSQL

Just push your changes → Render deploys automatically.

---

## 📈 **Future Improvements**

* JWT authentication for admin routes (optional)
* Rate limiting for commission submissions
* File uploads (refs, screenshots)
* Analytics dashboard for client demand
* Full product backend for store page

---

## ✨ **Author**

**Plexdi**
Designer • Developer • Builder
Focused on creating clean visuals + efficient systems.

