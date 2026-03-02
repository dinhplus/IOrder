# IOrder — Business Requirements Document

> **Author:** Product Owner (PO)
> **Date:** 2026-03-02
> **Version:** 1.0
> **Status:** Approved

---

## 1. Executive Summary

IOrder là nền tảng đặt bàn và gọi món thông minh dành cho nhà hàng, quán ăn, café và các cơ sở F&B. Hệ thống cho phép nhiều nhà hàng vận hành độc lập trên cùng một nền tảng (multi-tenant SaaS), khách hàng đặt món qua mã QR tại bàn, và nhà hàng quản lý toàn bộ quy trình từ nhận đơn đến thanh toán.

---

## 2. Business Goals

| # | Goal | Success Metric |
|---|---|---|
| G1 | Tăng tốc độ gọi món và thanh toán | Thời gian từ scan QR → xác nhận đơn < 3 phút |
| G2 | Giảm sai sót trong truyền đơn | Tỷ lệ đơn hàng sai < 1% |
| G3 | Hỗ trợ nhà hàng marketing qua membership | 30% khách quay lại sau 3 tháng |
| G4 | Hỗ trợ đa nhà hàng trên cùng nền tảng | Onboard ≥ 10 nhà hàng trong Q1 |
| G5 | Tích hợp thanh toán không tiền mặt | 60% đơn hàng thanh toán điện tử |

---

## 3. Stakeholders

| Role | Description |
|---|---|
| **Platform Admin** | Quản lý toàn bộ hệ thống, onboard nhà hàng mới |
| **Restaurant Owner** | Chủ nhà hàng, quản lý menu, bàn, nhân viên, báo cáo |
| **Restaurant Staff** | Phục vụ, thu ngân, bếp nhận order |
| **Customer (Dine-in)** | Khách ăn tại chỗ, scan QR → gọi món → thanh toán |
| **Customer (Membership)** | Khách thành viên, tích điểm, dùng voucher |

---

## 4. Feature Modules

### 4.1 Multi-Tenant Restaurant Management

**Mô tả:** Mỗi nhà hàng là một tenant độc lập với dữ liệu và cài đặt riêng.

**User Stories:**

| ID | As a | I want to | So that |
|---|---|---|---|
| US-01 | Platform Admin | Onboard nhà hàng mới (tên, logo, địa chỉ, múi giờ) | Nhà hàng bắt đầu hoạt động trên nền tảng |
| US-02 | Restaurant Owner | Cài đặt thông tin nhà hàng (tên, mô tả, giờ mở cửa, ảnh đại diện) | Khách hàng nhận diện đúng nhà hàng |
| US-03 | Restaurant Owner | Quản lý nhân viên (tạo/sửa/xóa tài khoản staff) | Phân quyền truy cập hệ thống |
| US-04 | Platform Admin | Xem danh sách tất cả nhà hàng và trạng thái hoạt động | Giám sát hệ thống |

**Acceptance Criteria (US-01):**
- [ ] Tạo tenant mới với subdomain hoặc mã định danh riêng
- [ ] Upload logo (JPG/PNG, max 2MB)
- [ ] Cài đặt múi giờ, đơn vị tiền tệ
- [ ] Tài khoản Restaurant Owner được tạo tự động sau khi onboard

---

### 4.2 Menu Management

**Mô tả:** Nhà hàng quản lý danh sách sản phẩm, danh mục, giá cả và tồn kho.

**User Stories:**

| ID | As a | I want to | So that |
|---|---|---|---|
| US-05 | Restaurant Owner | Tạo/sửa/xóa danh mục menu (ví dụ: Khai vị, Món chính, Tráng miệng, Đồ uống) | Phân nhóm sản phẩm |
| US-06 | Restaurant Owner | Thêm món mới với tên, mô tả, giá, ảnh, tag dị ứng, tag dinh dưỡng | Hiển thị đúng thông tin cho khách |
| US-07 | Restaurant Owner | Cài đặt modifier cho món (ví dụ: kích thước, mức độ cay, topping) | Khách tùy chỉnh món theo sở thích |
| US-08 | Restaurant Owner | Đánh dấu món hết hàng (sold out) theo thời gian thực | Khách không đặt món hết |
| US-09 | Restaurant Owner | Thiết lập giờ phục vụ cho từng món (ví dụ: chỉ phục vụ buổi sáng) | Quản lý thực đơn theo thời điểm |
| US-10 | Restaurant Owner | Upload ảnh menu dạng PDF hoặc hình ảnh | Đồng bộ thực đơn vật lý |

**Acceptance Criteria (US-06):**
- [ ] Tên, giá bắt buộc; mô tả và ảnh tùy chọn
- [ ] Giá có thể có đơn vị tùy chỉnh (VND, USD)
- [ ] Ảnh món ăn max 5MB, format JPG/PNG/WebP
- [ ] Tag: `vegetarian`, `vegan`, `gluten-free`, `spicy`, `contains-nuts`
- [ ] Có thể kích hoạt/vô hiệu hóa từng món

---

### 4.3 Table & Floor Map Management

**Mô tả:** Hệ thống bản đồ bàn trực quan, hỗ trợ nhiều tầng/khu vực trong cùng một nhà hàng.

**User Stories:**

| ID | As a | I want to | So that |
|---|---|---|---|
| US-11 | Restaurant Owner | Tạo sơ đồ tầng (floor plan) với tên và mô tả | Quản lý bàn theo khu vực |
| US-12 | Restaurant Owner | Thêm bàn vào sơ đồ với vị trí (x, y), hình dạng, số ghế | Hình ảnh hóa nhà hàng |
| US-13 | Restaurant Owner | Đặt tên bàn tùy chỉnh (A1, B2, Bàn VIP 1,...) | Nhận diện bàn dễ dàng |
| US-14 | Restaurant Owner | Tạo nhiều tầng/khu vực (Tầng 1, Tầng 2, Sân thượng, Khu vực ngoài trời) | Quản lý nhà hàng nhiều tầng |
| US-15 | Restaurant Staff | Xem trạng thái bàn (trống/có khách/đang gọi món) theo thời gian thực | Điều phối phục vụ |
| US-16 | Restaurant Owner | Tạo mã QR cho từng bàn (in ấn hoặc download PNG/PDF) | Khách scan để gọi món |

**Acceptance Criteria (US-16):**
- [ ] Mỗi bàn có QR code chứa URL duy nhất (format: `https://app.iorder.vn/{restaurant_slug}/table/{table_id}`)
- [ ] QR code có thể tải về dạng PNG (300x300 px) hoặc PDF A4 (multiple QR per page)
- [ ] QR code hết hạn có thể được regenerate
- [ ] QR code mã hóa JWT ngắn hạn (6 giờ) để bảo mật

---

### 4.4 QR Code Ordering Flow

**Mô tả:** Quy trình khách hàng scan QR → xem menu → gọi món → xác nhận → theo dõi đơn.

**User Stories:**

| ID | As a | I want to | So that |
|---|---|---|---|
| US-17 | Customer | Scan QR trên bàn bằng camera điện thoại | Truy cập menu không cần download app |
| US-18 | Customer | Xem menu đầy đủ với ảnh, giá, mô tả | Chọn món đúng ý |
| US-19 | Customer | Thêm, xóa, chỉnh số lượng món trong giỏ hàng | Tự do tùy chỉnh đơn |
| US-20 | Customer | Thêm ghi chú cho từng món hoặc toàn đơn | Yêu cầu đặc biệt (ít muối, không hành,...) |
| US-21 | Customer | Xác nhận đơn và theo dõi trạng thái (đã nhận, đang chế biến, sẵn sàng phục vụ) | Biết khi nào món đến |
| US-22 | Customer | Gọi thêm món sau khi đã đặt đơn trước | Linh hoạt bổ sung trong bữa ăn |
| US-23 | Customer | Yêu cầu thanh toán (bill) từ ứng dụng | Không cần gọi nhân viên |

**Acceptance Criteria (US-17):**
- [ ] Scan QR mở web browser (PWA), không cần cài app
- [ ] Tự động detect ngôn ngữ từ browser (VN/EN)
- [ ] Hiển thị tên nhà hàng và số bàn ngay trên trang
- [ ] Truy cập được trong vòng < 2 giây (3G connection)

---

### 4.5 Order State Machine

**Mô tả:** Vòng đời đơn hàng với các trạng thái rõ ràng và hành động được phép.

**Trạng thái đơn hàng:**

```
DRAFT → SUBMITTED → CONFIRMED → IN_PREPARATION → READY → SERVED → PAYMENT_REQUESTED → PAID → CLOSED
                  ↓                                                     ↓
               REJECTED                                            CANCELLED
```

| State | Mô tả | Actor |
|---|---|---|
| `DRAFT` | Khách đang tạo đơn, chưa gửi | Customer |
| `SUBMITTED` | Khách xác nhận gửi đơn lên | Customer |
| `CONFIRMED` | Nhà hàng xác nhận nhận đơn | Staff/Auto |
| `REJECTED` | Nhà hàng từ chối đơn (hết hàng, hết bàn) | Staff |
| `IN_PREPARATION` | Bếp đang chế biến | Kitchen Staff |
| `READY` | Món đã sẵn sàng phục vụ | Kitchen Staff |
| `SERVED` | Đã phục vụ tại bàn | Staff |
| `PAYMENT_REQUESTED` | Khách yêu cầu thanh toán | Customer/Staff |
| `PAID` | Thanh toán thành công | Payment System |
| `CLOSED` | Đơn hoàn tất, bàn đã trống | System/Staff |
| `CANCELLED` | Đơn bị hủy trước khi chuẩn bị | Customer/Staff |

**User Stories:**

| ID | As a | I want to | So that |
|---|---|---|---|
| US-24 | Kitchen Staff | Nhận thông báo ngay khi có đơn mới (CONFIRMED) | Bắt đầu chế biến kịp thời |
| US-25 | Kitchen Staff | Cập nhật trạng thái đơn (IN_PREPARATION → READY) | Nhân viên mang món ra bàn |
| US-26 | Staff | Xem tất cả đơn hàng theo trạng thái trên màn hình bếp (Kitchen Display System) | Ưu tiên xử lý đơn |
| US-27 | Customer | Nhận thông báo khi trạng thái đơn thay đổi | Biết tiến độ đơn hàng |
| US-28 | Staff | Hủy đơn hàng với lý do | Quản lý exception cases |

---

### 4.6 Payment Integration

**Mô tả:** Hỗ trợ nhiều phương thức thanh toán bao gồm tiền mặt, thẻ ngân hàng, và ví điện tử.

**Phương thức thanh toán được hỗ trợ:**

| Method | Provider | Ghi chú |
|---|---|---|
| Tiền mặt | N/A | Thanh toán trực tiếp, staff xác nhận |
| Chuyển khoản ngân hàng | VietQR | QR code chuẩn VietQR |
| Ví điện tử | MoMo | Tích hợp MoMo Payment Gateway |
| Ví điện tử | ZaloPay | Tích hợp ZaloPay SDK |
| Ví điện tử | ShopeePay | Tích hợp Airpay/ShopeePay |
| Thẻ quốc tế | Adyen | Visa/Mastercard/JCB |
| Thẻ quốc tế | Stripe | Dự phòng cho thị trường quốc tế |

**User Stories:**

| ID | As a | I want to | So that |
|---|---|---|---|
| US-29 | Customer | Chọn phương thức thanh toán phù hợp | Linh hoạt thanh toán |
| US-30 | Customer | Quét QR VietQR để thanh toán qua app ngân hàng | Thanh toán nhanh không cần tiền mặt |
| US-31 | Customer | Thanh toán qua MoMo/ZaloPay/ShopeePay | Dùng ví điện tử phổ biến |
| US-32 | Customer | Nhận hóa đơn điện tử (e-receipt) sau thanh toán | Lưu lại giao dịch |
| US-33 | Restaurant Owner | Xem lịch sử giao dịch và báo cáo doanh thu | Quản lý tài chính |
| US-34 | Restaurant Owner | Cấu hình phương thức thanh toán được chấp nhận | Phù hợp với nhà hàng |

**Acceptance Criteria (US-30):**
- [ ] Tạo VietQR động chứa số tiền chính xác
- [ ] Polling hoặc webhook khi thanh toán thành công
- [ ] Timeout xử lý 15 phút, sau đó đơn trở về trạng thái chờ
- [ ] Hiển thị trạng thái thanh toán real-time

---

### 4.7 Membership & Marketing

**Mô tả:** Hệ thống thành viên giúp nhà hàng xây dựng lòng trung thành và thực hiện marketing.

**User Stories:**

| ID | As a | I want to | So that |
|---|---|---|---|
| US-35 | Customer | Đăng ký tài khoản thành viên (số điện thoại/email) | Tích điểm và nhận ưu đãi |
| US-36 | Customer | Tích điểm sau mỗi lần thanh toán (10,000 VND = 1 điểm) | Đổi điểm lấy ưu đãi |
| US-37 | Customer | Đổi điểm lấy voucher giảm giá | Được thưởng vì trung thành |
| US-38 | Customer | Nhận thông báo về ưu đãi đặc biệt (sinh nhật, ngày lễ) | Không bỏ lỡ ưu đãi |
| US-39 | Restaurant Owner | Tạo chương trình tích điểm tùy chỉnh | Phù hợp chiến lược marketing |
| US-40 | Restaurant Owner | Tạo voucher (giảm % hoặc giảm cố định, có thời hạn) | Khuyến khích khách quay lại |
| US-41 | Restaurant Owner | Xem báo cáo membership (số thành viên mới, top customers) | Đánh giá hiệu quả marketing |
| US-42 | Restaurant Owner | Gửi push notification hoặc SMS đến nhóm thành viên | Thông báo ưu đãi, sự kiện |

**Membership Tiers:**

| Tier | Điều kiện | Quyền lợi |
|---|---|---|
| **Bronze** | Đăng ký | 1 điểm / 10,000 VND, không có ưu đãi mặc định |
| **Silver** | Chi tiêu ≥ 2,000,000 VND | 1.5 điểm / 10,000 VND, giảm 5% |
| **Gold** | Chi tiêu ≥ 5,000,000 VND | 2 điểm / 10,000 VND, giảm 10%, ưu tiên phục vụ |
| **Platinum** | Chi tiêu ≥ 15,000,000 VND | 3 điểm / 10,000 VND, giảm 15%, dịch vụ VIP |

---

### 4.8 Kitchen Display System (KDS)

**Mô tả:** Màn hình bếp hiển thị đơn hàng theo thời gian thực để bếp xử lý hiệu quả.

**User Stories:**

| ID | As a | I want to | So that |
|---|---|---|---|
| US-43 | Kitchen Staff | Xem màn hình KDS hiển thị các đơn đang chờ | Xử lý theo thứ tự ưu tiên |
| US-44 | Kitchen Staff | Cập nhật trạng thái từng món trong đơn (đang làm/xong) | Theo dõi tiến độ chính xác |
| US-45 | Kitchen Staff | Nhận âm thanh/rung khi có đơn mới | Không bỏ sót đơn |
| US-46 | Kitchen Staff | Xem thời gian chờ của từng đơn | Ưu tiên đơn đang đợi lâu |

---

### 4.9 Reservation & Booking

**Mô tả:** Đặt bàn trước qua app hoặc điện thoại.

**User Stories:**

| ID | As a | I want to | So that |
|---|---|---|---|
| US-47 | Customer | Đặt bàn trước theo ngày, giờ, số người | Chắc chắn có chỗ ngồi |
| US-48 | Customer | Nhận xác nhận đặt bàn qua SMS/email | Lưu lại thông tin |
| US-49 | Customer | Hủy hoặc đổi lịch đặt bàn | Linh hoạt thay đổi kế hoạch |
| US-50 | Restaurant Owner | Xem lịch đặt bàn và xác nhận/từ chối | Quản lý capacity nhà hàng |
| US-51 | Restaurant Owner | Cài đặt quy tắc đặt bàn (thời gian sớm nhất, muộn nhất, số chỗ tối đa) | Kiểm soát booking |

---

### 4.10 Analytics & Reporting

**Mô tả:** Dashboard báo cáo kinh doanh cho chủ nhà hàng.

**User Stories:**

| ID | As a | I want to | So that |
|---|---|---|---|
| US-52 | Restaurant Owner | Xem doanh thu theo ngày/tuần/tháng | Nắm bắt xu hướng kinh doanh |
| US-53 | Restaurant Owner | Xem top 10 món bán chạy | Tối ưu menu |
| US-54 | Restaurant Owner | Xem giờ cao điểm trong ngày/tuần | Sắp xếp nhân viên phù hợp |
| US-55 | Restaurant Owner | Xuất báo cáo dạng PDF/Excel | Lưu trữ và kiểm toán |
| US-56 | Platform Admin | Xem báo cáo tổng hợp toàn platform | Quản lý vận hành |

---

### 4.11 Notification System

**Mô tả:** Thông báo đến đúng người đúng lúc qua đa kênh.

**Notification Events:**

| Event | Recipient | Channel |
|---|---|---|
| Đơn mới | Kitchen Staff | KDS, Push, Sound |
| Trạng thái đơn thay đổi | Customer | Push/SMS |
| Món sắp hết hàng | Restaurant Owner | Push/Email |
| Thanh toán thành công | Customer + Owner | Push/Email |
| Ưu đãi membership | Customer | Push/SMS/Email |
| Đặt bàn được xác nhận | Customer | SMS/Email |

---

### 4.12 Multi-language Support

**Mô tả:** Giao diện hỗ trợ đa ngôn ngữ cho khách hàng quốc tế.

| Language | Code | Priority |
|---|---|---|
| Tiếng Việt | `vi` | P1 |
| English | `en` | P1 |
| 中文 (giản thể) | `zh-CN` | P2 |
| 한국어 | `ko` | P2 |
| 日本語 | `ja` | P2 |

---

## 5. Non-Functional Requirements

### 5.1 Performance

| Metric | Target |
|---|---|
| API response time (p95) | < 200ms |
| Page load time (mobile 3G) | < 3 giây |
| QR scan to menu display | < 2 giây |
| Real-time order update latency | < 1 giây |
| System uptime | ≥ 99.9% (monthly) |

### 5.2 Security

- Mã hóa TLS 1.3 cho tất cả kết nối
- JWT token expiry: access token 1 giờ, refresh token 30 ngày
- QR token: 6 giờ
- Rate limiting: 100 req/minute per IP
- PCI DSS compliance cho thanh toán thẻ (delegated to Adyen/Stripe)
- GDPR/PDPA compliance cho dữ liệu cá nhân

### 5.3 Scalability

- Hỗ trợ đồng thời: 10,000 concurrent users per restaurant cluster
- Menu items per restaurant: không giới hạn
- Tables per restaurant: lên đến 500 bàn
- Restaurants per platform: không giới hạn (horizontal scaling)

### 5.4 Availability

- Disaster recovery: RPO < 1 giờ, RTO < 4 giờ
- Backup hàng ngày với retention 30 ngày
- Multi-AZ deployment trên AWS

---

## 6. User Journey Maps

### 6.1 Dine-in Customer Journey

```
1. Khách ngồi xuống bàn
2. Scan QR code trên bàn
3. Web app tải → hiển thị menu nhà hàng
4. [Optional] Đăng nhập thành viên để tích điểm
5. Browse menu → thêm món vào giỏ
6. Xem giỏ hàng → chỉnh sửa → thêm ghi chú
7. Xác nhận đặt → gửi đơn
8. Theo dõi trạng thái đơn real-time
9. Nhận món
10. Yêu cầu bill → chọn phương thức thanh toán
11. Thanh toán → nhận e-receipt
12. [Optional] Đánh giá trải nghiệm
```

### 6.2 Kitchen Staff Journey

```
1. Đăng nhập KDS (tablet/màn hình tại bếp)
2. Xem đơn hàng mới (thông báo âm thanh)
3. Confirm nhận đơn
4. Cập nhật trạng thái từng món khi chế biến
5. Mark READY khi hoàn thành
6. Nhân viên phục vụ mang món ra
```

---

## 7. Out of Scope (v1.0)

- Delivery (giao hàng đến nhà) — sẽ xem xét v2.0
- Loyalty program tích hợp với bên thứ ba (loyalty.vn, ...)
- Inventory management (quản lý kho nguyên liệu)
- Employee scheduling
- Integration với phần mềm kế toán (MISA, Fast)
- Tích hợp với hệ thống POS vật lý (Sapo, KiotViet)

---

## 8. Acceptance Criteria — Platform Level

- [ ] Nhà hàng có thể onboard và bắt đầu nhận đơn trong < 30 phút
- [ ] Khách hàng có thể gọi món chỉ bằng camera điện thoại (không cần app)
- [ ] Toàn bộ vòng đời đơn hàng được theo dõi và log đầy đủ
- [ ] Thanh toán thành công rate ≥ 99% (sau retry)
- [ ] Zero data loss trong tất cả trường hợp lỗi

---

## 9. Glossary

| Term | Definition |
|---|---|
| Tenant | Một nhà hàng trên nền tảng IOrder |
| QR Token | JWT ngắn hạn nhúng trong QR code của bàn |
| KDS | Kitchen Display System — màn hình hiển thị đơn tại bếp |
| Floor Plan | Sơ đồ bố trí bàn của một tầng/khu vực |
| Modifier | Tùy chọn thêm của món (size, topping, level spicy,...) |
| E-receipt | Hóa đơn điện tử gửi sau thanh toán |
| VietQR | Chuẩn QR thanh toán liên ngân hàng Việt Nam |
