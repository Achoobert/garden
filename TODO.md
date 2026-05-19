# iGarden PoC TODO

## Setup
- [x] Initialize Go project with SQLite support `+sqlite3`
- [x] Project structure: `main.go`, `db/`, `web/`

## Database
- [x] Define `User` table: email, time_budget, cash_budget, space_type
- [x] Define `Plant` table: name, instructions (JSON/Relational)
- [x] Define `UserPlant` table (Garden): user_id, plant_id, last_action_time
- [x] Seed data: Dandelion, Knapweed with 5-min intervals

## Backend (Go)
- [x] Server on port 30022
- [x] User signup/login (placeholder/simple)
- [x] API: Update User Profile
- [x] API: Get User Garden
- [x] API: Take action (water/feed/etc)
- [x] Notification engine: Calculate next due tasks

## Frontend (PWA)
- [x] HTML/CSS/JS for User registration/profile
- [x] Dashboard: List of plants and their status
- [x] PWA: `manifest.json`, Service Worker
- [x] Simple UI for "Water", "Feed", etc.

## Push Notifications
- [/] Web Push API implementation or simplified simulated notifications for PoC (Simulated via Browser Notifications)

## High-End Concierge (Backend)
- [ ] Rename 'Username' to 'Full Name & Title' (e.g., Khun Somchai)
- [ ] Implement Thai-Season logic: May 18 = Early Rainy Season (Bangkok)
- [ ] Seed Thai Plants: Holy Basil, Bird's Eye Chili, Nam Dok Mai Mango, Butterfly Pea
- [ ] Instructions: Water frequency adjusted for high humidity

## Ultra-Modern UI (Frontend)
- [ ] Integrate Tailwind CSS & Lucide Icons (CDN)
- [ ] Concierge Greeting: "Good Morning, [Title]. The monsoons have arrived in Sukhumvit."
- [ ] Buy Now buttons: Integrated Shopee/Lazada referral placeholders for premium fertilizer
- [ ] PWA: Support for Push API and Add-to-Home-Screen for iOS

## Growth Strategy
- [ ] 1200 User Roadmap: Focus on Sukhumvit/Thong Lor urbanites, referral loops, and luxury condo partnerships

