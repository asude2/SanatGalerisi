<template>
  <div class="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div v-if="artwork" class="max-w-6xl mx-auto bg-white rounded-3xl shadow-2xl overflow-hidden border border-gray-100">
      <div class="flex flex-col lg:flex-row">
        
        <div class="lg:w-1/2 relative bg-gray-200">
          <img :src="artwork.imageUrl" :alt="artwork.title" class="w-full h-full object-cover min-h-[500px]" />
          
          <div v-if="artwork.IsCampaign" class="absolute top-6 right-6 bg-red-600 text-white px-6 py-2 rounded-full font-black shadow-2xl animate-pulse">
            %{{ artwork.DiscountRate }} KAMPANYA!
          </div>

          <button @click="router.back()" class="absolute top-6 left-6 bg-white/90 p-3 rounded-full shadow-lg hover:bg-white transition-all cursor-pointer">
            ⬅️ Geri Dön
          </button>
        </div>

        <div class="lg:w-1/2 p-12 flex flex-col justify-center">
          <div class="space-y-6">
            <div>
              <span class="text-blue-600 font-bold tracking-widest uppercase text-sm">Sanat Eseri Detayı</span>
              <h1 class="text-5xl font-black text-gray-900 mt-2">{{ artwork.title }}</h1>
              <p class="text-2xl text-gray-500 font-medium mt-1 flex items-center gap-4">Sanatçı: <span class="text-gray-800">{{ artwork.artist }}</span>
                <button @click="showArtistModal = true" class="text-sm bg-blue-100 text-blue-700 px-4 py-2 rounded-xl font-bold hover:bg-blue-600 hover:text-white transition-all shadow-sm">
                  Sanatçıyı Görüntüle 🔍
                </button>
              </p>
            </div>

            <div class="border-y border-gray-100 py-8">
              <h3 class="text-lg font-bold text-gray-800 mb-3 underline decoration-blue-500 underline-offset-4">Eser Açıklaması</h3>
              <p class="text-gray-600 leading-relaxed text-xl">
                {{ artwork.description || 'Bu eser için henüz bir açıklama eklenmemiş.' }}
              </p>
            </div>

            <div class="pt-4 p-6 bg-blue-50 rounded-2xl border border-blue-100">
              <p class="text-gray-400 text-sm uppercase font-bold tracking-tighter">Ödenecek Tutar</p>
              <div class="flex flex-col">
                <div v-if="artwork.IsCampaign" class="space-y-1">
                  <div class="flex items-center gap-3">
                    <span class="text-xl text-gray-400 line-through font-bold">{{ artwork.price.toLocaleString() }} ₺</span>
                    <span class="bg-red-100 text-red-600 text-xs px-2 py-1 rounded-md font-black">-%{{ artwork.DiscountRate }}</span>
                  </div>
                  <p class="text-5xl font-black text-blue-600">{{ finalCalculatedPrice.toLocaleString() }} ₺</p>
                  <p class="text-green-600 text-sm font-bold flex items-center gap-1">
                    ✨ Kampanya indirimi uygulandı!
                  </p>
                </div>
                
                <div v-else>
                   <p class="text-5xl font-black text-blue-600">{{ (artwork.price - appliedDiscount).toLocaleString() }} ₺</p>
                   <p v-if="appliedDiscount > 0" class="text-green-600 text-sm font-bold mt-2">✅ Kupon indirimi uygulandı!</p>
                </div>
              </div>
            </div>

            <div class="space-y-4 my-6">
              <div class="p-4 bg-gray-50 rounded-xl border border-gray-200">
                <label class="block text-xs font-bold text-gray-500 uppercase mb-2">Ödeme Yöntemi</label>
                <select v-model="selectedPaymentMethod" class="w-full p-2 border border-gray-300 rounded-lg outline-none focus:ring-2 focus:ring-blue-500">
                  <option value="Kredi Kartı">Kredi Kartı</option>
                  <option value="Havale/EFT">Havale / EFT</option>
                  <option value="Cüzdan">Uygulama Bakiyesi</option>
                </select>
              </div>

              <div class="p-4 bg-gray-50 rounded-xl border-2 border-dashed border-gray-300">
                <label class="block text-xs font-bold text-gray-500 uppercase mb-2">Ekstra İndirim Kuponu</label>
                <div class="flex gap-2">
                  <input v-model="couponCode" type="text" placeholder="Kupon Kodunuz" class="flex-1 p-2 border border-gray-300 rounded-lg outline-none">
                  <button @click="applyCoupon" class="bg-gray-800 text-white px-4 py-2 rounded-lg font-bold hover:bg-black transition-all">Uygula</button>
                </div>
              </div>
            </div>

            <button @click="buyArtwork" class="w-full bg-blue-600 text-white py-4 rounded-xl font-bold text-lg hover:bg-blue-700 transition-all shadow-lg">
              Satın Almayı Tamamla ({{ finalCalculatedPrice.toLocaleString() }} ₺) 🛒
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="flex flex-col items-center justify-center h-[60vh] space-y-4">
      <div class="animate-spin rounded-full h-16 w-16 border-b-4 border-blue-600"></div>
      <p class="text-2xl text-gray-500 font-medium tracking-tight">Eser detayları yükleniyor...</p>
    </div>

    </div>
</template>
<script setup>
import { ref, onMounted, watch, computed } from 'vue'; // computed eklendi
import { useRoute, useRouter } from 'vue-router';
import axios from 'axios';

const route = useRoute();
const router = useRouter();
const artwork = ref(null);
const showArtistModal = ref(false);
const artistInfo = ref(null);

// Satın alma ve kupon state'leri
const couponCode = ref('');
const appliedDiscount = ref(0);
const selectedPaymentMethod = ref('Kredi Kartı');

// --- HESAPLAMA MANTIĞI (YENİ) ---
// Bu özellik, hem kampanya indirimini hem de kuponu tek tek hesaplar
const finalCalculatedPrice = computed(() => {
  if (!artwork.value) return 0;
  
  let price = artwork.value.price;
  
  // 1. Önce Veritabanındaki Kampanya İndirimini uygula
  if (artwork.value.IsCampaign) {
    price = price * (1 - artwork.value.DiscountRate / 100);
  }
  
  // 2. Sonra varsa kupon indirimini düş (Eksiye düşmemesi için Max kullandık)
  const final = price - (appliedDiscount.value || 0);
  return Math.max(0, final); 
});

// --- Kupon Kontrolü ---
const applyCoupon = async () => {
  if (!couponCode.value) return;
  try {
    const response = await axios.get(`http://localhost:8080/check-coupon?code=${couponCode.value}`);
    appliedDiscount.value = response.data.discount;
    alert(`🎫 Kupon Uygulandı: ${appliedDiscount.value} TL indirim kazandın!`);
  } catch (error) {
    appliedDiscount.value = 0;
    alert("Geçersiz kupon kodu! 😔");
  }
};

const buyArtwork = async () => {
  const userEmail = localStorage.getItem('userEmail');
  if (!userEmail) {
    alert("Lütfen önce giriş yapın! 👤");
    return;
  }

  // Eserin ID'sini ve HESAPLANMIŞ fiyatını alıyoruz
  const artworkId = artwork.value.id; 
  const priceToPay = finalCalculatedPrice.value; // Artık kampanya dahil fiyat gidiyor!

  if (!confirm(`${artwork.value.title} eserini ${priceToPay.toLocaleString()} ₺ karşılığında satın almak istiyor musunuz?`)) return;

  try {
    const response = await axios.post(`http://localhost:8080/buy-artwork/${artworkId}`, {
      email: userEmail,
      artworkId: artworkId,
      price: priceToPay, // Backend'e indirimli fiyatı yolluyoruz
      paymentMethod: selectedPaymentMethod.value || "Cüzdan"
    });

    alert("Satın alma başarılı! 🎨");
    router.push('/profile');
  } catch (error) {
    console.error("Hata detayı:", error.response);
    const status = error.response?.status;
    
    if (status === 400) {
      if (confirm("Bakiyeniz yetersiz! 💸 Bakiye yüklemek için profile gitmek ister misiniz?")) {
        router.push('/profile');
      }
    } else {
      alert("İşlem başarısız: " + (error.response?.data?.message || "Bir hata oluştu."));
    }
  }
};

// Sanatçı bilgisini getir
const fetchArtistInfo = async (artistName) => {
  try {
    const response = await axios.get(`http://localhost:8080/artist?name=${encodeURIComponent(artistName)}`);
    artistInfo.value = response.data;
  } catch (error) {
    console.error("Sanatçı bilgisi yüklenemedi:", error);
    artistInfo.value = { name: artistName, biography: '', artworksCount: 0, workshopsCount: 0 };
  }
};

watch(showArtistModal, (newVal) => {
  if (newVal && artwork.value) fetchArtistInfo(artwork.value.artist);
});

onMounted(async () => {
  try {
    // Eserleri getirirken backend'den yeni kolonların geldiğinden eminiz
    const response = await axios.get('http://localhost:8080/artworks');
    artwork.value = response.data.find(a => a.id === parseInt(route.params.id));
    
    if (!artwork.value) {
      console.error("Eser bulunamadı!");
    }
  } catch (error) {
    console.error("Detaylar yüklenemedi:", error);
  }
});
</script>