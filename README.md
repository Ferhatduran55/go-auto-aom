# Oto Parça Sipariş Sistemi 🔧

Tek EXE dosyası olarak çalışan, sunucu gerektirmeyen masaüstü uygulaması.

**Vue.js 3** + **TailwindCSS** + **Bleve** (gömülü Elasticsearch alternatifi) + **WebView2** (Windows native) kullanır.

![Version](https://img.shields.io/badge/version-25.12--stable-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![Platform](https://img.shields.io/badge/platform-Windows%2010%2F11-lightgrey)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8)
![Vue.js](https://img.shields.io/badge/Vue.js-3.5+-4FC08D)

---

## 🚀 Özellikler (Son Güncellemeler)

- ✅ **Gelişmiş arama + tarih filtresi birlikte çalışır** — Gelişmiş arama kriterleri seçildiğinde tarih filtresi (Bugün/Tümü/Tarih Aralığı) de uygulanır.
- ✅ **Sipariş kalemlerinde içerik arama** — Ürün adı, OEM, birim fiyat ve durum tek bir arama kutusundan aranabilir; eşleşmeyen ürün sayısı gösterilir.
- ✅ **Çoklu seçim ve toplu silme** — Satır bazında checkbox ile seçip toplu silme yapılabilir.
- ✅ **WhatsApp resim paylaşımı** — PNG, panoya kopyalanır; WhatsApp penceresi açılır ve Ctrl+V ile yapıştırma ile paylaşılır.
- ✅ **Tema uyumlu bileşenler** — Dropdown'lar, toast mesajları ve butonlar karanlık/aydınlık temalarla uyumludur.

---

## 📋 Gereksinimler (Sadece Derleme İçin)

| Araç | Versiyon  | Açıklama                                    |
| ----- | --------- | --------------------------------------------- |
| Go    | 1.21+     | Ana dil                                       |
| GCC   | MinGW-w64 | CGO desteği için                            |
| Bun   | 1.x       | Frontend paket yöneticisi (veya Node.js 18+) |

---

## 🛠️ Derleme

### Otomatik Build Script (Önerilen)

```powershell
.\build.ps1
```

### Manuel Derleme

```powershell
# 1. Frontend build
cd frontend
bun install
bun run build

# 2. EXE metadata güncelle (isteğe bağlı - versioninfo -> resource.syso)
# Goversioninfo kullanarak versioninfo.json'dan resource.syso oluşturun (Windows exe metadata)
# goversioninfo aracı kuruluysa:
goversioninfo -icon="assets/App.ico" -o="resource.syso" versioninfo.json

# 3. Go build
cd ..
go mod tidy
go build -ldflags="-H windowsgui -s -w" -o OtoParcaSiparis.exe .
```

> Not: resource.syso dosyasını oluşturmazsanız, exe içindeki "Ayrıntılar" kısmı eski bilgiler göstermeye devam edebilir.

### Geliştirme Modu

```powershell
# Frontend dev server
cd frontend
bun run dev

# Go (ayrı terminal)
go run .
```

---

## 💻 Kullanım

### Başlarken

1. `OtoParcaSiparis.exe` dosyasını çift tıklayın
2. Uygulama otomatik olarak pencerede açılacaktır

### Sipariş Oluşturma

1. **Sipariş Bilgileri**: Üst çubuktan sipariş başlığı ve müşteri adı girin
2. **Ürün Ekleme**: Sol panelden ürün bilgilerini doldurun
3. **Listeye Ekle**: Butona tıklayın (veya Enter)
4. **Kaydet**: "💾 Siparişi Kaydet" butonuna tıklayın

### Çıktı Alma

- **📄 TXT**: WhatsApp uyumlu düz metin
- **🖼️ PNG**: Tek tuşla panoya kopyalama + WhatsApp ile paylaşma talimatı

### Katalog Yönetimi

- Header'daki **📚 Katalog** butonuna tıklayın
- Ürünleri ve müşterileri görüntüleyin, düzenleyin veya silin

---

## 📁 Proje Yapısı (Özet)

```
...existing code...
```

(Detaylı proje yapısı için dosyaya bakınız.)

---

## 🔧 Geliştirme

### Gerekli Araçlar

```powershell
# Go
go version  # 1.21+

# Bun
bun --version  # 1.x

# GCC (CGO için)
gcc --version
```

### Komutlar

```powershell
# Frontend geliştirme
cd frontend
bun run dev

# Go çalıştır
go run .

# Test
go test ./...

# Build
.\build.ps1
```

---

## 📄 Lisans

MIT License

---

## 👨‍💻 Geliştirici

**Ferhat Duran**

[![GitHub](https://img.shields.io/badge/GitHub-Ferhatduran55-181717?logo=github)](https://github.com/Ferhatduran55)

**Şirket**: Durasoft

---

## 🤝 Katkıda Bulunun

Adımlar:

1. Fork edin
2. Branch oluşturun
3. Değişiklik yapıp PR açın

---

## ⭐ Destek

Proje hoşunuza gittiyse GitHub'da yıldız bırakabilirsiniz!

## 📁 Proje Yapısı

```
go-auto-aom/
├── main.go                  # Ana Go uygulaması + WebView2
├── go.mod                   # Go modül dosyası
├── build.ps1                # Otomatik build script
├── versioninfo.json         # EXE versiyon bilgileri
├── resource.syso            # Windows kaynakları (ikon)
│
├── assets/
│   └── App.ico              # Uygulama ikonu
│
├── storage/
│   └── bleve_store.go       # Bleve veritabanı katmanı
│
├── web/                     # Build çıktısı (Go embed)
│   ├── index.html
│   ├── favicon.ico
│   └── assets/
│
└── frontend/                # Vue.js 3 kaynak kodları
    ├── src/
    │   ├── api/             # Go API wrapper
    │   ├── composables/     # Vue composables
    │   ├── components/      # Vue bileşenleri
    │   │   ├── AppHeader.vue
    │   │   ├── OrderBar.vue
    │   │   ├── ProductForm.vue
    │   │   ├── OrderItems.vue
    │   │   ├── OrdersList.vue
    │   │   ├── CatalogModal.vue
    │   │   ├── Modals.vue
    │   │   └── Toast.vue
    │   ├── App.vue
    │   ├── main.js
    │   └── style.css
    ├── tailwind.config.js
    ├── vite.config.js
    └── package.json
```

---

## 🗄️ Veri Depolama

Veriler otomatik olarak şu konumda saklanır:

**Windows**: `%APPDATA%\OtoParcaSiparis\`

```
OtoParcaSiparis/
├── bleve_index/        # Bleve search index
├── orders/             # Sipariş JSON dosyaları
├── products/           # Ürün JSON dosyaları
└── customers/          # Müşteri JSON dosyaları
```

---

## 📤 Çıktı Örnekleri

### TXT Çıktısı (WhatsApp Uyumlu)

```
📋 *SİPARİŞ: Aralık Ayı İhale*
👤 *Müşteri:* ABC Otomotiv
📅 *Tarih:* 24.12.2024

───────────────────────────────────

*1. Megane 4 Orijinal Far*
   🔢 OEM: 260601234R
   📦 Adet: 2
   🟢 Orijinal
   💰 *₺15.000,00*

*2. Ön Fren Balatası*
   🔢 OEM: 34116850568
   📦 Adet: 4
   🟡 Çıkma
   💰 *₺1.800,00*

───────────────────────────────────

💵 *GENEL TOPLAM: ₺16.800,00*
```

### PNG Çıktısı

Profesyonel görünümlü fiyat teklifi:

- Gradient header (mavi → mor)
- Ürün adı + Adet (mor renkte)
- Fiyat (yeşil renkte)
- Beyaz arka plan

---

## 🎨 Ekran Görüntüleri

### Ana Ekran

```
┌──────────────────────────────────────────────────────────────┐
│  🔧 OTO PARÇA  │  📚 Katalog  📊 Geçmiş  🆕 Yeni  🌙 Tema   │
├──────────────────────────────────────────────────────────────┤
│  📋 Başlık: [________]  👤 Müşteri: [________]  📅 Tarih    │
├────────────────┬─────────────────────────┬───────────────────┤
│  ➕ Ürün Ekle  │  📋 Sipariş Kalemleri   │  📚 Kayıtlı       │
│                │                         │     Siparişler    │
│  Ürün: [____]  │  ┌─────────────────────┐│                   │
│  OEM:  [____]  │  │ 1. Far    2 Adet    ││  [Bugün] [Tümü]  │
│  Adet: [1]     │  │    ₺15.000         ││                   │
│  Fiyat: [__]   │  ├─────────────────────┤│  ┌─────────────┐ │
│  Durum: [▼]    │  │ 2. Balata  4 Adet   ││  │ Sipariş #1  │ │
│                │  │    ₺1.800          ││  │ ₺16.800     │ │
│  [Listeye Ekle]│  └─────────────────────┘│  └─────────────┘ │
│                │                         │                   │
│  [Temizle]     │  TOPLAM: ₺16.800       │                   │
│  [Listeyi Sil] │                         │                   │
│                │  [💾 Kaydet] [📄 TXT] [🖼️ PNG] │            │
└────────────────┴─────────────────────────┴───────────────────┘
```

---

## 🔧 Geliştirme

### Gerekli Araçlar

```powershell
# Go
go version  # 1.21+

# Bun
bun --version  # 1.x

# GCC (CGO için)
gcc --version
```

### Komutlar

```powershell
# Frontend geliştirme
cd frontend
bun run dev

# Go çalıştır
go run .

# Test
go test ./...

# Build
.\build.ps1
```

---

## 📝 Teknolojiler

| Katman                    | Teknoloji   | Açıklama                       |
| ------------------------- | ----------- | -------------------------------- |
| **Backend**         | Go 1.21+    | Ana uygulama dili                |
| **GUI**             | WebView2    | Windows native web görünümü  |
| **Database**        | Bleve       | Gömülü full-text arama motoru |
| **Frontend**        | Vue.js 3    | Reaktif UI framework             |
| **Build Tool**      | Vite        | Hızlı frontend build           |
| **Styling**         | TailwindCSS | Utility-first CSS                |
| **Package Manager** | Bun         | Hızlı JS paket yöneticisi     |

---

## 📄 Lisans

MIT License

---

## 👨‍💻 Geliştirici

**Ferhat Duran**

[![GitHub](https://img.shields.io/badge/GitHub-Ferhatduran55-181717?logo=github)](https://github.com/Ferhatduran55)

**Yapılanma**: Durasoft

Oto parça sektörü için özel olarak tasarlanmış sipariş yönetim sistemi.

**Özellikler**:

- 🔓 Open Source (MIT)
- 🔌 Offline çalışma
- ⚡ Hızlı ve hafif
- 🎯 Kolay kullanım
- 📊 Profesyonel çıktılar

---

## 🤝 Katkıda Bulunun

1. Bu repoyu fork edin
2. Feature branch oluşturun (`git checkout -b feature/amazing-feature`)
3. Değişikliklerinizi commit edin (`git commit -m 'Add some amazing feature'`)
4. Branch'e push edin (`git push origin feature/amazing-feature`)
5. Pull Request açın

---

## ⭐ Destek

Bu projeyi beğendiyseniz, GitHub'da ⭐ vererek destek olabilirsiniz!
