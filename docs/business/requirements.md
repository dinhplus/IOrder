# IOrder — Business Requirements Document

> **Author:** Product Owner (PO)
> **Date:** 2026-03-02
> **Version:** 2.0
> **Status:** Approved

---

## 1. Executive Summary

IOrder is a smart table-ordering platform for restaurants, cafés, and food & beverage outlets. The system enables multiple restaurants to operate independently on a shared multi-tenant SaaS platform. Customers place orders by scanning a QR code at their table, and the restaurant manages the entire workflow — from order intake to payment — through role-specific interfaces for every staff position.

---

## 2. Business Goals

| # | Goal | Success Metric |
|---|---|---|
| G1 | Accelerate ordering and payment | Time from QR scan → order confirmed < 3 minutes |
| G2 | Eliminate order transmission errors | Error rate < 1% of all orders |
| G3 | Drive repeat visits via membership | 30% of customers return within 3 months |
| G4 | Support multiple restaurants on one platform | Onboard ≥ 10 restaurants in Q1 |
| G5 | Increase cashless payment adoption | 60% of orders paid electronically |
| G6 | Reduce food waste via inventory visibility | Stockout incidents reduced by 40% |

---

## 3. Stakeholders & Restaurant Staff Roles

| Role | System Role | Primary Responsibilities |
|---|---|---|
| **Platform Admin** | `platform_admin` | Manage entire platform, onboard new restaurants |
| **Restaurant Owner** | `owner` | Manage menu, floor plans, staff accounts, reports |
| **Restaurant Manager** | `manager` | Oversee daily operations, approve refunds, view analytics |
| **Order Taker / Waiter** | `waiter` | Take orders at the table or POS terminal, assist customers |
| **Server** | `server` | Deliver food from kitchen to table, mark orders as SERVED |
| **Kitchen Chef** | `kitchen` | Receive and prepare food orders via Kitchen Display System (KDS) |
| **Bartender / Beverage** | `bartender` | Receive and prepare drink orders via Bar Display System (BDS) |
| **Cashier** | `cashier` | Process payments, issue receipts, handle cash reconciliation |
| **Warehouse Manager** | `warehouse` | Manage ingredient and supply inventory |
| **Accountant** | `accountant` | View financial reports, export accounting data |
| **Customer (Dine-in)** | N/A | Scan QR → browse menu → order → pay |
| **Customer (Member)** | N/A | Earn points, use vouchers, receive promotions |

---

## 4. Staff Role Workflows

This section defines the end-to-end operational flow for each restaurant staff position.

---

### 4.1 Waiter / Order Taker

The waiter assists dine-in customers who prefer to be served traditionally, or takes orders on behalf of customers at tables that do not use QR self-ordering.

**Workflow:**

```
1. Customer is seated → Waiter greets and hands QR or takes order verbally
2. Waiter opens the Staff POS interface (tablet/mobile app)
3. Waiter selects the correct table on the floor map
4. Waiter adds items to the order, including modifiers and special notes
5. Waiter submits the order → system routes to Kitchen and/or Bar
6. Waiter monitors order status on the POS interface
7. When all items are READY, Waiter collects from pass-through and delivers
8. Waiter marks items as SERVED on the POS interface
9. Waiter can add additional items requested by the customer (new order round)
10. Waiter informs the Cashier when the customer requests the bill
```

**Permissions:**
- View floor map and table status
- Create/modify orders for any table (subject to manager approval for modifications after CONFIRMED)
- Mark items as SERVED
- Cannot process payments (cashier role required)

---

### 4.2 Kitchen Chef

The chef receives food orders on the Kitchen Display System (KDS) and coordinates food preparation.

**Workflow:**

```
1. Chef logs in to KDS screen (large display tablet or monitor at kitchen pass)
2. New order appears automatically when order status = CONFIRMED
3. Audible alert sounds for every new incoming order
4. Chef reviews order items, special notes, and table number
5. Chef taps "Start Cooking" → items move to IN_PREPARATION lane
6. Each cook handles their station; head chef tracks overall timing
7. When all items in an order are prepared, chef taps "Ready" → status → READY
8. KDS notifies the server to collect and deliver the order
9. Chef can flag an item as unavailable → triggers REJECTED notification to waiter/customer
10. At end of shift, chef views completed orders count and average preparation time
```

**KDS Display Lanes:**

| Lane | Orders Shown |
|---|---|
| **New** | CONFIRMED orders not yet started |
| **In Progress** | Orders in IN_PREPARATION state |
| **Ready** | Completed orders waiting for pickup |
| **Done** | Served orders (scrollable history) |

**Permissions:**
- View food-category items only (beverage items go to BDS)
- Update item status: pending → preparing → ready
- Flag items as unavailable (sold out)
- View order history for current shift

---

### 4.3 Bartender / Beverage Staff

The bartender receives drink orders on the Bar Display System (BDS), a separate KDS filtered to beverage items only.

**Workflow:**

```
1. Bartender logs in to BDS (dedicated tablet at the bar)
2. Drink orders from new orders appear automatically (filtered by item category = "Beverages")
3. Bartender prepares drinks in parallel with kitchen food preparation
4. Bartender marks each drink as READY when complete
5. Server collects drinks from the bar pass-through
6. Bartender can update stock levels for specific drinks inline
```

**Permissions:**
- View beverage-category items only
- Update drink item status
- Mark beverages as sold out

---

### 4.4 Server / Food Runner

The server collects prepared food and drinks from the kitchen/bar and delivers them to the correct table.

**Workflow:**

```
1. Server monitors the READY lane on the floor-map view (mobile app / wall display)
2. Notification appears when an order moves to READY status
3. Server collects items from kitchen pass-through
4. Server navigates to the correct table using the floor map
5. Server taps "Deliver" on the table order → items marked SERVED
6. Server can split delivery if only some items are ready (partial delivery)
7. Server notifies cashier if customer requests the bill
```

**Permissions:**
- View all orders in READY state across all floor plans
- Mark orders/items as SERVED
- Cannot modify order items

---

### 4.5 Cashier

The cashier handles payment processing, receipt issuance, and end-of-shift cash reconciliation.

**Workflow:**

```
1. Customer or waiter requests bill → order moves to PAYMENT_REQUESTED
2. Cashier opens the bill on the POS interface
3. Cashier reviews itemized bill; can apply:
   a. Member voucher / discount code
   b. Manual discount (requires manager approval > 20%)
   c. Complimentary item (manager approval required)
4. Cashier selects or confirms payment method:
   - Cash → enter amount received, system calculates change
   - VietQR → display dynamic QR to customer
   - MoMo / ZaloPay / ShopeePay → display payment QR or deep link
   - Credit/Debit card (Adyen terminal)
5. Payment confirmed → order status → PAID
6. System auto-generates e-receipt, sends to customer (SMS/email if member)
7. Cashier closes the session for that table
8. At end of shift: cashier runs daily report, reconciles cash drawer
9. Cashier exports shift summary for accountant handoff
```

**Permissions:**
- View and process all PAYMENT_REQUESTED orders
- Apply vouchers and discounts (manual discounts > 20% require manager approval)
- Issue manual refunds up to a configurable limit
- View shift-level sales summary
- Cannot create or modify menu items

---

### 4.6 Warehouse Manager

The warehouse manager controls ingredient and supply inventory to prevent stockouts and reduce waste.

**Workflow:**

```
1. Warehouse manager logs in to the Inventory module
2. View current stock levels for all ingredients and supplies
3. Receive deliveries: scan or manually enter quantities received
4. Set minimum stock thresholds per ingredient
5. System auto-alerts when ingredient falls below threshold
6. Perform stock-take (physical count) and reconcile with system
7. View consumption report: ingredients used per menu item per day
8. Create and send purchase orders to suppliers
9. Link ingredient depletion to sold menu items (recipe mapping)
10. Export inventory reports for accountant
```

**Inventory Features:**
- Ingredient master list with unit of measure (kg, litre, piece)
- Recipe mapping: each menu item consumes X units of ingredient Y
- Auto-deduct from inventory when order is CONFIRMED (reserved) and SERVED (actual)
- Low-stock alerts via push notification and dashboard warning
- Supplier management: contact info, lead times, preferred suppliers
- Purchase order workflow: draft → sent → received → verified

**Permissions:**
- Full access to Inventory module
- Cannot access order or payment data directly
- Can view which menu items are high-consumption

---

### 4.7 Accountant

The accountant has read-only access to financial data and can export reports for external accounting systems.

**Workflow:**

```
1. Accountant logs in with read-only financial role
2. View daily/weekly/monthly revenue summary by payment method
3. View itemized transaction log (all PAID orders with amounts, discounts, taxes)
4. View discount and voucher usage report
5. View refund log
6. Export data in multiple formats:
   - PDF: formatted financial report
   - Excel/CSV: raw transaction data for MISA, Fast Accounting
   - XML: VAT declaration format
7. View membership revenue attribution (revenue from member vs. non-member)
8. View cost of goods sold (COGS) based on inventory consumption
9. Reconcile daily cash drawer report vs. system totals
```

**Permissions:**
- Read-only access to all financial reports
- Export any report to PDF/Excel/CSV
- Cannot create, modify, or cancel orders
- Cannot process payments

---

## 5. Feature Modules

### 5.1 Multi-Tenant Restaurant Management

**Description:** Each restaurant is an independent tenant with its own data and settings.

**User Stories:**

| ID | As a | I want to | So that |
|---|---|---|---|
| US-01 | Platform Admin | Onboard a new restaurant (name, logo, address, timezone) | The restaurant can start operating on the platform |
| US-02 | Restaurant Owner | Configure restaurant details (name, description, opening hours, cover photo) | Customers correctly identify the restaurant |
| US-03 | Restaurant Owner | Manage staff accounts (create/edit/deactivate) with specific roles | Access control is properly enforced |
| US-04 | Platform Admin | View all restaurants and their operational status | Monitor platform health |

**Acceptance Criteria (US-01):**
- [ ] Create new tenant with a unique slug used in QR URLs
- [ ] Upload logo (JPG/PNG, max 2MB)
- [ ] Set timezone and currency
- [ ] Restaurant Owner account is automatically created upon onboarding

---

### 5.2 Menu Management

**Description:** Restaurants manage products, categories, pricing, modifiers, and availability.

**User Stories:**

| ID | As a | I want to | So that |
|---|---|---|---|
| US-05 | Restaurant Owner | Create/edit/delete menu categories (e.g., Starters, Mains, Desserts, Beverages) | Items are logically grouped |
| US-06 | Restaurant Owner | Add a new item with name, description, price, image, allergy tags, and nutritional tags | Customers see accurate item information |
| US-07 | Restaurant Owner | Configure modifiers for an item (e.g., size, spice level, toppings) | Customers can customize their order |
| US-08 | Restaurant Owner / Kitchen | Mark an item as sold out in real time | Customers cannot order unavailable items |
| US-09 | Restaurant Owner | Set serving hours per item (e.g., breakfast items available 06:00–11:00 only) | Menu reflects time-based availability |
| US-10 | Restaurant Owner | Assign items to Beverages or Food category | Order is routed to Kitchen or Bar correctly |

**Acceptance Criteria (US-06):**
- [ ] Name and price are required; description and image are optional
- [ ] Price stored in smallest currency unit (e.g., xu for VND)
- [ ] Image max 5MB, formats: JPG/PNG/WebP
- [ ] Tags: `vegetarian`, `vegan`, `gluten-free`, `spicy`, `contains-nuts`, `alcohol`
- [ ] Items can be activated/deactivated without deletion

---

### 5.3 Table & Floor Map Management

**Description:** Visual floor plan system supporting multiple floors and zones within a single restaurant.

**User Stories:**

| ID | As a | I want to | So that |
|---|---|---|---|
| US-11 | Restaurant Owner | Create a floor plan with a name and description | Tables can be organized by zone |
| US-12 | Restaurant Owner | Add tables to the floor plan with position (x, y), shape, and seat count | The restaurant layout is visualized |
| US-13 | Restaurant Owner | Assign custom names to tables (A1, B2, VIP-01, etc.) | Tables are easily identified |
| US-14 | Restaurant Owner | Create multiple floors/zones (Floor 1, Rooftop, Outdoor Terrace) | Multi-floor restaurants are supported |
| US-15 | Staff (any role) | See real-time table status (available / occupied / ordering) on the floor map | Staff can coordinate service efficiently |
| US-16 | Restaurant Owner | Generate a QR code for each table (download as PNG or PDF) | Customers can scan to order |

**Acceptance Criteria (US-16):**
- [ ] Each table QR encodes a unique URL: `https://order.iorder.vn/{slug}/table/{id}?token={jwt}`
- [ ] QR downloadable as PNG (300×300 px) or PDF A4 (6 QR codes per page with restaurant branding)
- [ ] Expired QR codes can be regenerated by staff
- [ ] QR token is a 6-hour JWT signed with a per-restaurant secret

---

### 5.4 QR Code Ordering Flow

**Description:** Customer journey from QR scan to order confirmation and real-time status tracking.

**User Stories:**

| ID | As a | I want to | So that |
|---|---|---|---|
| US-17 | Customer | Scan the QR code on the table using my phone camera | I can access the menu without downloading an app |
| US-18 | Customer | Browse the full menu with photos, prices, and descriptions | I can make an informed choice |
| US-19 | Customer | Add, remove, and adjust quantities in my cart | I can freely customize my order |
| US-20 | Customer | Add a note to individual items or the whole order | Special requests are communicated clearly |
| US-21 | Customer | Confirm my order and track its real-time status | I know when my food will arrive |
| US-22 | Customer | Order additional items after my first order has been served | I can add to my meal flexibly |
| US-23 | Customer | Request the bill from the app | I do not need to flag down a waiter |

**Acceptance Criteria (US-17):**
- [ ] Scan opens a PWA in the phone browser — no app installation required
- [ ] Browser language is auto-detected (Vietnamese / English)
- [ ] Restaurant name and table number are displayed prominently
- [ ] Page loads in under 2 seconds on a 3G connection

---

### 5.5 Order State Machine

**Description:** A well-defined order lifecycle with explicit states and permitted transitions per role.

**Order States:**

```
DRAFT → SUBMITTED → CONFIRMED → IN_PREPARATION → READY → SERVED → PAYMENT_REQUESTED → PAID → CLOSED
                  ↓                                                      ↓
               REJECTED                                             CANCELLED
```

| State | Description | Permitted Actor |
|---|---|---|
| `DRAFT` | Customer is building the cart; not yet submitted | Customer / Waiter |
| `SUBMITTED` | Order submitted; awaiting restaurant confirmation | Customer / Waiter |
| `CONFIRMED` | Restaurant has accepted the order | Waiter / Auto-confirm |
| `REJECTED` | Restaurant rejected the order (e.g., items out of stock) | Waiter / Manager |
| `IN_PREPARATION` | Kitchen/Bar is actively preparing the order | Chef / Bartender |
| `READY` | All items prepared and waiting for pickup | Chef / Bartender |
| `SERVED` | Order delivered to the table | Server / Waiter |
| `PAYMENT_REQUESTED` | Customer or waiter has requested the bill | Customer / Waiter |
| `PAID` | Payment successfully processed | Cashier / Payment System |
| `CLOSED` | Session closed, table is reset to available | System / Manager |
| `CANCELLED` | Order cancelled before preparation began | Customer / Waiter / Manager |

**User Stories:**

| ID | As a | I want to | So that |
|---|---|---|---|
| US-24 | Chef / Bartender | Receive an immediate alert when a new order is CONFIRMED | I can start preparation without delay |
| US-25 | Chef / Bartender | Update the order status (IN_PREPARATION → READY) | The server knows when to pick up |
| US-26 | Server | See all READY orders on the floor map view | I can deliver them promptly |
| US-27 | Customer | Receive real-time status notifications | I know the progress of my order |
| US-28 | Waiter / Manager | Cancel an order with a reason | Exception cases are handled cleanly |

---

### 5.6 Payment Integration

**Description:** Support for multiple payment methods including cash, bank transfer, and digital wallets.

**Supported Payment Methods:**

| Method | Provider | Notes |
|---|---|---|
| Cash | N/A | Direct payment; cashier confirms manually |
| Bank transfer QR | VietQR | Open Vietnamese inter-bank QR standard |
| Digital wallet | MoMo | MoMo Payment Gateway integration |
| Digital wallet | ZaloPay | ZaloPay SDK integration |
| Digital wallet | ShopeePay | Airpay/ShopeePay integration |
| International card | Adyen | Visa/Mastercard/JCB; Apple Pay, Google Pay |
| International card | Stripe | Fallback for international markets |

**User Stories:**

| ID | As a | I want to | So that |
|---|---|---|---|
| US-29 | Customer | Choose my preferred payment method | I can pay in the way that suits me |
| US-30 | Customer | Scan a VietQR code to pay via my banking app | I can pay quickly without cash |
| US-31 | Customer | Pay via MoMo/ZaloPay/ShopeePay | I can use popular digital wallets |
| US-32 | Customer | Receive an electronic receipt after payment | I have a record of the transaction |
| US-33 | Restaurant Owner / Accountant | View transaction history and revenue reports | I can manage finances accurately |
| US-34 | Restaurant Owner | Configure which payment methods are accepted | I can match my restaurant's setup |
| US-35 | Cashier | Process a split bill for a table | Groups can pay separately |

**Acceptance Criteria (US-30):**
- [ ] Dynamic VietQR generated with exact order amount
- [ ] Real-time status update via webhook or polling
- [ ] 15-minute payment timeout; order reverts to SERVED if unpaid
- [ ] Real-time payment status displayed to customer and cashier

---

### 5.7 Membership & Loyalty Program

**Description:** A flexible member loyalty system helping restaurants build customer retention and run targeted marketing campaigns.

#### 5.7.1 Member Registration & Profiles

**User Stories:**

| ID | As a | I want to | So that |
|---|---|---|---|
| US-36 | Customer | Register as a member using my phone number or email | I can earn points and receive perks |
| US-37 | Customer | View my points balance, tier status, and transaction history | I know where I stand |
| US-38 | Customer | Link my profile to my order by logging in before or after ordering | My spend is always tracked |

#### 5.7.2 Points & Tiers

**Membership Tiers:**

| Tier | Qualification | Points Earning | Default Discount | Benefits |
|---|---|---|---|---|
| **Bronze** | Register | 1 pt / 10,000 VND | None | Access to member-only promotions |
| **Silver** | Cumulative spend ≥ 2,000,000 VND | 1.5 pts / 10,000 VND | 5% off total bill | Priority seating notification |
| **Gold** | Cumulative spend ≥ 5,000,000 VND | 2 pts / 10,000 VND | 10% off total bill | Free dessert on birthday, priority service |
| **Platinum** | Cumulative spend ≥ 15,000,000 VND | 3 pts / 10,000 VND | 15% off total bill | Dedicated server, complimentary welcome drink |

**Tier Rules:**
- Tier is evaluated on rolling 12-month spend (not lifetime)
- Tier is downgraded at the annual review if spend threshold is not maintained
- Points never expire while the member remains active (logs in at least once per year)

#### 5.7.3 Vouchers & Discount Codes

**Voucher Types:**

| Type | Description | Example |
|---|---|---|
| **Percentage discount** | Reduces bill by a percentage | 20% off total bill |
| **Fixed-amount discount** | Reduces bill by a flat amount | 50,000 VND off |
| **Free item** | Grants a specific menu item for free | Free dessert on birthday |
| **Buy-X-Get-Y** | Bonus item when a threshold is met | Buy 2 mains, get 1 free drink |
| **Minimum spend** | Discount applies only if bill exceeds a threshold | 100,000 VND off on orders ≥ 500,000 VND |
| **Category discount** | Discount applies to a specific menu category | 30% off all beverages |
| **Points redemption** | Points converted to a discount voucher | 100 pts = 10,000 VND voucher |

**Voucher Lifecycle:**

```
CREATED → ACTIVE → [REDEEMED | EXPIRED | REVOKED]
```

**User Stories:**

| ID | As a | I want to | So that |
|---|---|---|---|
| US-39 | Customer | Enter a discount code at checkout | I can apply a promotional offer |
| US-40 | Customer | Redeem points for a voucher in-app | I get rewarded for my loyalty |
| US-41 | Customer | See all my available vouchers before paying | I do not miss applicable discounts |
| US-42 | Restaurant Owner | Create a voucher with type, value, validity period, and usage limit | I can run targeted promotions |
| US-43 | Restaurant Owner | Set voucher constraints: single-use, per-member limit, minimum order amount | Promotions are used as intended |
| US-44 | Restaurant Owner | Distribute vouchers to a segment of members (e.g., all Gold+ members) | Marketing campaigns are targeted |
| US-45 | Restaurant Owner | View voucher redemption report (redeemed / unused / expired) | I can measure campaign effectiveness |
| US-46 | Cashier | Manually apply a voucher at point of sale | Walk-in promotions can be honored |

**Acceptance Criteria (US-42):**
- [ ] Voucher must have: code (auto-generated or custom), type, value, validity start/end date
- [ ] Optional: minimum order amount, maximum discount cap, per-member usage limit, total usage limit
- [ ] Voucher code is case-insensitive, alphanumeric, 6–12 characters
- [ ] Expired or fully-redeemed vouchers are automatically deactivated

#### 5.7.4 Promotional Campaigns

**Campaign Types:**

| Campaign | Description |
|---|---|
| **Happy Hour** | Time-based discount on selected items or categories |
| **Birthday Reward** | Automatic voucher sent 3 days before member birthday |
| **Referral Program** | Member earns bonus points when a referred friend makes their first order |
| **Milestone Reward** | Bonus voucher when member reaches a spend or visit milestone |
| **Seasonal Promotion** | Date-range-based discount (e.g., Valentine's Day, Tết) |
| **Re-engagement Campaign** | Voucher sent to members who have not visited in 60+ days |

**User Stories:**

| ID | As a | I want to | So that |
|---|---|---|---|
| US-47 | Restaurant Owner | Create an automated campaign triggered by member events | Marketing runs without manual effort |
| US-48 | Restaurant Owner | Preview a campaign before activating it | I can review it before sending |
| US-49 | Restaurant Owner | Send push notifications and SMS to campaign targets | Members are aware of the promotion |
| US-50 | Restaurant Owner | View campaign performance: sends, open rate, redemption rate, revenue attributed | I can optimize future campaigns |

#### 5.7.5 Membership Discount Application Rules

When multiple discounts are applicable, the system applies them in this priority order:

1. **Member tier discount** (e.g., Gold 10%) is calculated on the subtotal
2. **Active voucher** is applied after tier discount (can be additive up to a configurable cap)
3. **Points redemption voucher** is applied last
4. **Minimum total floor**: Bill total can never drop below configurable minimum (e.g., 10% of original value) to prevent abuse

**Example:**
- Bill subtotal: 500,000 VND
- Gold tier discount (10%): −50,000 VND → 450,000 VND
- Applied voucher (20% off, max cap 80,000 VND): −80,000 VND (capped) → 370,000 VND
- Points redemption (50 pts = 5,000 VND): −5,000 VND → **Final: 365,000 VND**

---

### 5.8 Kitchen Display System (KDS) & Bar Display System (BDS)

**Description:** Real-time order display screens for kitchen and bar staff to manage and fulfill orders efficiently.

**User Stories:**

| ID | As a | I want to | So that |
|---|---|---|---|
| US-51 | Chef | See incoming orders sorted by time on the KDS screen | I process orders in correct sequence |
| US-52 | Chef | Update the status of individual items within an order | Partial completions are tracked |
| US-53 | Chef | Receive an audible and visual alert for every new order | I never miss an incoming order |
| US-54 | Chef | See the waiting time for each order | I can prioritize orders that have been waiting longest |
| US-55 | Chef | Mark an item as unavailable directly on the KDS | The waiter and customer are notified immediately |
| US-56 | Bartender | View only beverage-category items on a separate BDS | Kitchen and bar work independently |

---

### 5.9 Inventory Management (Warehouse)

**Description:** Real-time ingredient and supply tracking to prevent waste and enable cost control.

**User Stories:**

| ID | As a | I want to | So that |
|---|---|---|---|
| US-57 | Warehouse Manager | Maintain an ingredient catalog with units of measure | All ingredients are tracked consistently |
| US-58 | Warehouse Manager | Map ingredients to menu items via recipes | Consumption is automatically tracked |
| US-59 | Warehouse Manager | Set minimum stock thresholds per ingredient | I receive alerts before running out |
| US-60 | Warehouse Manager | Record stock receipts when deliveries arrive | Inventory levels stay accurate |
| US-61 | Warehouse Manager | Perform periodic stock-takes and record adjustments | Physical and system counts are reconciled |
| US-62 | Warehouse Manager | Manage a supplier list with contact information and lead times | Purchase orders can be raised quickly |
| US-63 | Restaurant Owner | View a daily consumption vs. theoretical usage report | Waste and theft are detectable |

---

### 5.10 Reservation & Booking

**Description:** Advance table booking via app or phone, with confirmation and management tools.

**User Stories:**

| ID | As a | I want to | So that |
|---|---|---|---|
| US-64 | Customer | Book a table in advance by date, time, and party size | I am guaranteed a seat |
| US-65 | Customer | Receive a booking confirmation via SMS or email | I have a record of my reservation |
| US-66 | Customer | Cancel or reschedule my booking | I can adjust my plans freely |
| US-67 | Restaurant Owner / Manager | View the reservation calendar and confirm or decline bookings | I manage seating capacity |
| US-68 | Restaurant Owner | Configure booking rules (earliest/latest time, max party size, blackout dates) | Booking constraints match operations |

---

### 5.11 Analytics & Reporting

**Description:** Business intelligence dashboard for restaurant owners, managers, and accountants.

**User Stories:**

| ID | As a | I want to | So that |
|---|---|---|---|
| US-69 | Restaurant Owner | View revenue by day/week/month with trend charts | I understand business performance |
| US-70 | Restaurant Owner | View the top 10 best-selling items | I can optimize the menu |
| US-71 | Restaurant Owner | View peak hours heatmap | I schedule staff appropriately |
| US-72 | Accountant | Export revenue report as PDF, Excel, or CSV | I have data for accounting and tax |
| US-73 | Restaurant Owner | View membership KPIs (new sign-ups, tier distribution, points issued/redeemed) | I measure loyalty program effectiveness |
| US-74 | Accountant | View cost of goods sold (COGS) from inventory consumption | I calculate gross margin accurately |
| US-75 | Platform Admin | View aggregated platform-wide metrics | I monitor overall platform health |

---

### 5.12 Notification System

**Description:** Multi-channel notifications delivered to the right person at the right time.

**Notification Events:**

| Event | Recipient | Channel |
|---|---|---|
| New order submitted | Waiter (if manual confirm mode) | POS app push |
| Order confirmed | Kitchen Chef / Bartender | KDS/BDS alert + sound |
| Order status changed | Customer | Push / SMS |
| Item marked unavailable | Waiter, Customer | Push |
| Low stock alert | Warehouse Manager | Push / Email |
| Payment successful | Customer + Cashier | Push / Email |
| Membership voucher available | Customer | Push / SMS / Email |
| Birthday reward issued | Customer | Push / SMS |
| Booking confirmed | Customer | SMS / Email |
| Re-engagement campaign | Inactive members | Push / SMS |

---

### 5.13 Multi-Language Support

**Description:** Multi-language interface for international guests.

| Language | Code | Priority |
|---|---|---|
| Vietnamese | `vi` | P1 |
| English | `en` | P1 |
| Simplified Chinese | `zh-CN` | P2 |
| Korean | `ko` | P2 |
| Japanese | `ja` | P2 |

---

## 6. Non-Functional Requirements

### 6.1 Performance

| Metric | Target |
|---|---|
| API response time (p95) | < 200ms |
| Page load time (mobile 3G) | < 3 seconds |
| QR scan to menu display | < 2 seconds |
| Real-time order update latency | < 1 second |
| System uptime | ≥ 99.9% (monthly SLA) |

### 6.2 Security

- TLS 1.3 for all connections
- JWT token expiry: access token 1 hour, refresh token 30 days
- QR token: 6 hours, per-restaurant signing secret
- Rate limiting: 100 requests/minute per IP
- PCI DSS compliance for card payments (delegated to Adyen/Stripe)
- GDPR/PDPA compliance for personal data

### 6.3 Scalability

- 10,000 concurrent users per restaurant cluster
- Unlimited menu items per restaurant
- Up to 500 tables per restaurant
- Unlimited restaurants on the platform (horizontal scaling)

### 6.4 Availability

- Disaster recovery: RPO < 1 hour, RTO < 4 hours
- Daily automated backups with 30-day retention
- Multi-AZ deployment on AWS

---

## 7. User Journey Maps

### 7.1 Dine-in Customer Journey

```
1. Customer is seated at a table
2. Customer scans QR code on the table
3. PWA loads → restaurant menu is displayed
4. [Optional] Customer logs in as a member to earn points
5. Customer browses menu → adds items to cart
6. Customer reviews cart → edits quantities → adds notes
7. Customer confirms and submits order
8. Customer tracks order status in real time
9. Food and drinks are delivered by the server
10. Customer requests the bill → selects payment method
11. Customer pays → e-receipt is displayed and sent
12. [Optional] Customer rates the experience
```

### 7.2 Kitchen Chef Journey

```
1. Chef logs in to KDS at the kitchen station
2. New order appears with audible alert (status: CONFIRMED)
3. Chef taps "Start" → order moves to IN_PREPARATION
4. Chef prepares each item, updates item-level progress
5. Chef taps "Ready" → server is notified
6. Server collects and marks as SERVED
```

### 7.3 Cashier Journey

```
1. Customer or waiter requests bill → PAYMENT_REQUESTED
2. Cashier opens bill on POS interface
3. Cashier reviews items, applies any member voucher/discount
4. Cashier presents payment options to customer
5. Customer pays (cash / QR / card / digital wallet)
6. System confirms payment → order → PAID
7. E-receipt generated, sent to member profile
8. Table session closed, table status → Available
```

### 7.4 Warehouse Manager Journey

```
1. Warehouse manager reviews stock levels dashboard
2. Low-stock alerts are visible for critical ingredients
3. Manager creates a purchase order for supplier
4. Delivery arrives → manager records quantities received
5. System auto-updates inventory levels
6. Daily consumption report compared against menu sales
7. Adjustments made if discrepancies found
```

### 7.5 Accountant Journey

```
1. Accountant logs in with read-only finance role
2. Views daily revenue summary (breakdown by payment method)
3. Reviews voucher/discount impact on revenue
4. Exports transaction CSV for import into MISA/Fast
5. Reconciles cash drawer report from cashier
6. Generates monthly financial summary PDF for management
```

---

## 8. Out of Scope (v1.0)

- Home delivery (to be considered in v2.0)
- Third-party loyalty platform integrations (loyalty.vn, etc.)
- Employee scheduling and shift management
- Integration with physical POS hardware (Sapo, KiotViet)
- Online ordering for takeaway
- Table pre-payment (deposit at booking)

---

## 9. Platform-Level Acceptance Criteria

- [ ] A restaurant can be onboarded and start receiving orders within 30 minutes
- [ ] Customers can order using only a phone camera — no app installation required
- [ ] The full order lifecycle is tracked and audited with a complete event log
- [ ] Payment success rate ≥ 99% (after automatic retry)
- [ ] Zero data loss under any failure scenario
- [ ] All staff roles can access only the features and data permitted by their role

---

## 10. Glossary

| Term | Definition |
|---|---|
| Tenant | A single restaurant on the IOrder platform |
| QR Token | A short-lived JWT embedded in a table QR code |
| KDS | Kitchen Display System — order display screen at the kitchen pass |
| BDS | Bar Display System — order display screen at the bar |
| Floor Plan | A visual map of table layouts for one floor or zone |
| Modifier | A customization option for a menu item (size, spice level, topping, etc.) |
| E-receipt | An electronic receipt sent to the customer after payment |
| VietQR | Vietnamese inter-bank QR payment standard (NAPAS) |
| Voucher | A discount coupon that can be applied at checkout |
| Points | Loyalty currency earned by members for every purchase |
| Tier | A membership level (Bronze/Silver/Gold/Platinum) based on cumulative spend |
| Session | A table service session from first QR scan to bill closure |
| COGS | Cost of Goods Sold — calculated from recipe-mapped ingredient consumption |
