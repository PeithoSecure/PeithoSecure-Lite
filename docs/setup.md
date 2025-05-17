# 🛠️ PeithoSecure Lite™ — Setup Guide  
> A Post-Quantum IAM system that humiliates you before it serves you.

Welcome to the official guide for **PeithoSecure Lite™** — the roast-locked identity backend you can fork, clone, copy… but not run.  
Unless you plan to cry in `make roast`, read this first.

---

## 📦 Prerequisites

You’ll need:

- Go 1.21+ installed
- Git
- Make (`brew install make` on macOS)
- Curiosity, not entitlement

---

## 🚀 1. Clone the Repo

```bash
git clone https://github.com/PeithoSecure/PeithoSecure-Lite.git
cd PeithoSecure-Lite
```

You just downloaded a fully-structured IAM backend.  
But not the part that lets it run. That’s intentional.

---

## 🧱 2. Build the Backend

```bash
make build
```

If it fails, it’s you.  
If it succeeds, it still won’t run.  
That’s me.

---

## 🔓 3. Try to Run It

```bash
make run
```

---

## 🥩 4. Just Want the Roast?

```bash
make roast
```

This runs Roast Engine™ in standalone mode.  
No identity. Just identity crises.

---

## 📸 5. Capture the Pain

```bash
make screenshot
cat peitho-trace.log
```

This saves the roast log to `peitho-trace.log`.  
For sharing, submitting, or simply accepting your fate.

---

## 🧩 6. What's Included

| Feature                      | Status |
|------------------------------|--------|
| JWT Auth (Keycloak)          | ✅     |
| SQLite + Secure Routes       | ✅     |
| Email Verify / Reset         | ✅     |
| Swagger UI                   | ✅     |

---

## 🚫 What's Missing

| Feature                     | Status |
|-----------------------------|--------|
| `unlock.lic`                | ❌     |
| `peitho-core/`              | ❌     |
| `dev_signer.go`             | ❌     |
| PQC Signature Enforcement   | ❌     |

---

## 🧠 This Is Not Broken

If it fails — that means it’s working.  
It’s doing exactly what I designed it to do.

> Clone it. Fork it. Copy it.  
> If it runs, [let me know](mailto:peithoindia.co@gmail.com) — so I can lock it harder.
> If it doesn’t, don’t worry — the Roast Engine™ is working as intended.

---

## 🔓 Coming Soon: The Full PeithoSecure Lite™

This repo currently ships in **roast-only lockdown mode** —  
designed to demonstrate structure, not to be runnable out of the box.

The **complete version** will include:

- 🔐 Fully functional license validation  
- 🧠 Real `peitho-core/` engine  
- 🧪 Signature enforcement using Dilithium2  
- 🥷 Anti-tamper, branding lock, and runtime fingerprint checks  
- ✨ Full authentication, onboarding, and dashboard access

You'll be able to run the backend.  
You'll be able to unlock it with a real license.  
You’ll still get roasted — but only if you try to cheat.

---

## 📊 Final Feature Matrix

| Feature                            | Status Now   | Final Release     |
|------------------------------------|--------------|-------------------|
| JWT auth via Keycloak              | ✅            | ✅                |
| SQLite-based storage               | ✅            | ✅                |
| Email verify + reset               | ✅            | ✅                |
| `/check?email=...` adaptive flow   | ✅            | ✅                |
| `/refresh` token auto-renew        | ✅            | ✅                |
| Swagger docs for all routes        | ✅            | ✅                |
| PQC license structure              | 🧠 stubbed    | ✅ (Dilithium2)    |
| `unlock.lic` real validation       | ❌            | ✅                |
| `peitho-core/` module              | ❌            | ✅ (internal)     |
| `dev_signer.go` (CLI tool)         | ❌            | ✅ (optional)     |
| `observer/` real escalation engine | 🧠 stubbed    | ✅                |
| Branding, engine hash enforcement  | 🧠 simulated  | ✅                |
| TLS 1.3 (real cert support)        | ✅            | ✅                |
| Roast Engine™ (for failed runs)    | ✅            | ✅                |
| Actual onboarding / unlock flow    | ❌            | ✅                |
| SDK + API Swagger                  | ❌            | ✅                |

---

## 📚 Educational Purpose

PeithoSecure Lite™ is designed to **demonstrate** how a secure, modern IAM system is built — from:

- Keycloak-based JWT authentication  
- SQLite storage and token refresh  
- Email flows with deep linking  
- Secure route handling and metrics  
- License validation using post-quantum signatures

This repo exists to **educate**, **inspire**, and **protect**.  
It’s not just code — it’s a blueprint for how security should feel.

> We believe IAM shouldn’t just run — it should **refuse to run** unless it's trusted.

So feel free to:
- Study the structure
- Understand the patterns
- Integrate the ideas into your own systems

But don’t expect it to “just work.”  
That would defeat the entire lesson.

---

## 🧵 Need Help?

Open an issue. But remember:

> You’re not blocked by a bug.  
> You’re blocked by a locked core — and your own overconfidence.
