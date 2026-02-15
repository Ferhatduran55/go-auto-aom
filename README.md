# AutoManagement

**İstemci tabanlı ilişkisel yönetim uygulaması.** Sipariş, stok ve raporlama sistemini tek bir masaüstü uygulamasında birleştirir.

![Version](https://img.shields.io/badge/version-26.2.1-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![Platform](https://img.shields.io/badge/platform-Windows-lightgrey)

---

## ✨ Ne Yapar?

- 📦 **Sipariş Yönetimi** — Hızlı sipariş oluşturma, düzenleme ve takip
- 📊 **Stok Takibi** — Ürün giriş/çıkış, kritik stok uyarıları, hareket geçmişi
- 📈 **Raporlar** — Satış ve stok raporları, grafik görünümler
- 📤 **Kolay Paylaşım** — TXT ve PNG çıktıları, WhatsApp entegrasyonu
- 🌍 **Çoklu Dil** — Türkçe ve İngilizce dil desteği

---

## 🚀 Başlangıç

1. [Releases](https://github.com/Ferhatduran55/go-auto-aom/releases) sayfasından EXE dosyasını indirin.
2. Çift tıklayarak çalıştırın — **kurulum gerektirmez**.
3. Verileriniz otomatik olarak `%APPDATA%\AutoManagement` klasöründe güvenle saklanır.

---

## 💡 Özellikler

| Özellik | Açıklama |
|---------|----------|
| **Offline Çalışma** | İnternet bağlantısı gerektirmez |
| **Tek Dosya** | Kurulum ve sunucu gerektirmez (Portable) |
| **Otomatik Kayıt** | Veriler yerel diskte JSON formatında saklanır |
| **Tema Desteği** | Sistem temasına duyarlı (Açık/Koyu) |
| **Hızlı Arama** | Ürün, müşteri ve siparişlerde anlık arama (Bleve tabanlı) |

---

## 🛠️ Geliştirme

### Gereksinimler

- Go 1.21+
- Node.js 18+ (veya Bun)
- WebView2 Runtime (Windows'ta varsayılan gelir)

### Derleme

```powershell
.\build.ps1
```

Bu script otomatik olarak:
- Frontend'i derler (Vite)
- Kaynak dosyalarını gömer
- Tek bir çalıştırılabilir dosya (EXE) oluşturur

### Dağıtım

```
📁 Çıktı
└── AutoManagement-x64-v26.2.1.exe
```

> **Tek EXE = Tam İşlevsellik!** Başka dosya veya DLL gerekmez.

---

## 📄 Lisans

MIT License — Açık kaynak kodludur ve kullanımı serbesttir.

---

## 👤 Geliştirici

**Ferhat Duran** · [GitHub](https://github.com/Ferhatduran55)
