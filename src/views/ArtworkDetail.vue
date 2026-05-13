<template>
  <div class="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div v-if="artwork" class="max-w-6xl mx-auto bg-white rounded-3xl shadow-2xl overflow-hidden border border-gray-100">
      <div class="flex flex-col lg:flex-row">
        
        <div class="lg:w-1/2 relative bg-gray-200">
          <img :src="artwork.imageUrl" :alt="artwork.title" class="w-full h-full object-cover min-h-[500px]" />
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

            <div class="space-y-4">
              <div>
                <p class="text-gray-400 text-sm uppercase font-bold tracking-tighter mb-2">Kategori</p>
                <div class="flex flex-wrap gap-2">
                  <span 
                    v-if="artwork.category" 
                    class="inline-block bg-blue-100 text-blue-700 px-4 py-2 rounded-full font-semibold text-sm hover:bg-blue-200 transition-all"
                  >
                    {{ artwork.category }}
                  </span>
                  <span v-else class="text-gray-400 text-sm italic">Kategori belirtilmemiş</span>
                </div>
              </div>
            </div>

            <div class="pt-4">
              <p class="text-gray-400 text-sm uppercase font-bold tracking-tighter">Fiyat</p>
              <div class="flex items-center gap-3">
                <p class="text-4xl font-black text-blue-600">{{ (artwork.price - appliedDiscount).toLocaleString() }} ₺</p>
                <span v-if="appliedDiscount > 0" class="text-sm bg-green-100 text-green-700 px-2 py-1 rounded font-bold">
                  İndirim Uygulandı!
                </span>
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
                <label class="block text-xs font-bold text-gray-500 uppercase mb-2">İndirim Kuponu</label>
                <div class="flex gap-2">
                  <input v-model="couponCode" type="text" placeholder="Örn: SANAT100" class="flex-1 p-2 border border-gray-300 rounded-lg outline-none">
                  <button @click="applyCoupon" class="bg-gray-800 text-white px-4 py-2 rounded-lg font-bold hover:bg-black transition-all">Uygula</button>
                </div>
                <p v-if="appliedDiscount > 0" class="text-green-600 text-sm font-bold mt-2">✅ {{ appliedDiscount }} TL indirim uygulandı!</p>
              </div>
            </div>

            <button @click="buyArtwork" class="w-full bg-blue-600 text-white py-4 rounded-xl font-bold text-lg hover:bg-blue-700 transition-all shadow-lg">
              Şimdi Satın Al 🛒
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="flex flex-col items-center justify-center h-[60vh] space-y-4">
      <div class="animate-spin rounded-full h-16 w-16 border-b-4 border-blue-600"></div>
      <p class="text-2xl text-gray-500 font-medium tracking-tight">Eser detayları yükleniyor...</p>
    </div>

    <div v-if="showArtistModal && artistInfo" class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-3xl p-8 max-w-2xl w-full shadow-2xl max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-3xl font-bold text-gray-900">{{ artistInfo.name }}</h3>
          <button @click="showArtistModal = false" class="text-3xl text-gray-400 hover:text-gray-600">✕</button>
        </div>
        <div class="space-y-6">
          <div>
            <h4 class="text-lg font-bold text-gray-700 mb-3">Biyografi</h4>
            <p class="text-gray-600 leading-relaxed text-base whitespace-pre-wrap">
              {{ artistInfo.biography || 'Bu sanatçı henüz bir biyografi eklememiş.' }}
            </p>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div class="bg-blue-50 rounded-2xl p-4">
              <p class="text-gray-600 text-sm font-medium uppercase mb-1">Eser Sayısı</p>
              <p class="text-2xl font-bold text-blue-600">{{ artistInfo.artworksCount }}</p>
            </div>
            <div class="bg-blue-50 rounded-2xl p-4">
              <p class="text-gray-600 text-sm font-medium uppercase mb-1">Atölye Sayısı</p>
              <p class="text-2xl font-bold text-blue-600">{{ artistInfo.workshopsCount }}</p>
            </div>
          </div>
          <button @click="showArtistModal = false" class="w-full mt-4 py-3 bg-blue-600 text-white font-bold rounded-xl">Kapat</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
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

// --- Güncellenmiş Satın Alma Fonksiyonu ---
const buyArtwork = async () => {
  const userEmail = localStorage.getItem('userEmail');
  if (!userEmail) {
    alert("Lütfen önce giriş yapın! 👤");
    return;
  }

  const finalPrice = artwork.value.price - appliedDiscount.value;

  if (!confirm(`${artwork.value.title} eserini ${finalPrice.toLocaleString()} TL karşılığında satın almak istiyor musunuz?`)) return;

  try {
    await axios.post('http://localhost:8080/artworks/buy', {
      email: localStorage.getItem('userEmail'),
      artworkId: artwork.value.id,
      price: finalPrice, 
      paymentMethod: selectedPaymentMethod.value
    });

    alert("Satın alma başarılı! Sanat koleksiyonuna eklendi. 🎨");
    router.push('/profile');
  } catch (error) {
    if (error.response && error.response.status === 400) {
      if (confirm("Bakiyeniz yetersiz! 💸 Bakiye yüklemek için profile gitmek ister misiniz?")) {
        router.push('/profile');
      }
    } else if (error.response && error.response.status === 409) {
      alert("Maalesef bu eser çoktan satılmış! 😔");
    } else {
      alert("Satın alma başarısız: " + (error.response?.data?.message || "Bir hata oluştu."));
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
    const response = await axios.get('http://localhost:8080/artworks');
    artwork.value = response.data.find(a => a.id === parseInt(route.params.id));
  } catch (error) {
    console.error("Detaylar yüklenemedi:", error);
  }
});
</script>