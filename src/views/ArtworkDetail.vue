<template>
  <div class="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div v-if="artwork" class="max-w-6xl mx-auto bg-white rounded-3xl shadow-2xl overflow-hidden border border-gray-100">
      <div class="flex flex-col lg:flex-row">
        
        <!-- Sol Taraf: Büyük Eser Görseli -->
        <div class="lg:w-1/2 relative bg-gray-200">
          <img :src="artwork.imageUrl" :alt="artwork.title" class="w-full h-full object-cover min-h-[500px]" />
          <button @click="router.back()" class="absolute top-6 left-6 bg-white/90 p-3 rounded-full shadow-lg hover:bg-white transition-all cursor-pointer">
            ⬅️ Geri Dön
          </button>
        </div>

        <!-- Sağ Taraf: Detay Bilgileri -->
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
                <p class="text-gray-400 text-sm uppercase font-bold tracking-tighter mb-2">Kategoriler</p>
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

            <div class="flex items-center justify-between pt-4">
              <div>
                <p class="text-gray-400 text-sm uppercase font-bold tracking-tighter">Fiyat</p>
                <p class="text-4xl font-black text-blue-600">{{ artwork.price.toLocaleString() }} ₺</p>
              </div>
              <button @click="buyArtwork" class="bg-blue-600 text-white px-10 py-5 rounded-2xl font-bold text-xl hover:bg-blue-700 hover:scale-105 transition-all shadow-xl shadow-blue-100">
                Satın Al 
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Yükleniyor Durumu -->
    <div v-else class="flex flex-col items-center justify-center h-[60vh] space-y-4">
      <div class="animate-spin rounded-full h-16 w-16 border-b-4 border-blue-600"></div>
      <p class="text-2xl text-gray-500 font-medium tracking-tight">Eser detayları yükleniyor...</p>
    </div>

    <!-- Sanatçı Bilgisi Modal -->
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

          <div class="grid grid-cols-2 gap-4 pt-4 border-t border-gray-100">
            <div class="bg-blue-50 rounded-2xl p-4">
              <p class="text-gray-600 text-sm font-medium uppercase mb-1">Toplam Eserler</p>
              <p class="text-3xl font-bold text-blue-600">{{ artistInfo.artworksCount }}</p>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-4 pt-4 border-t border-gray-100">
            <div class="bg-blue-50 rounded-2xl p-4">
              <p class="text-gray-600 text-sm font-medium uppercase mb-1">Toplam Atölyeler</p>
              <p class="text-3xl font-bold text-blue-600">{{ artistInfo.workshopsCount }}</p>
            </div>
          </div>

          <button 
            @click="showArtistModal = false" 
            class="w-full mt-6 px-6 py-3 bg-blue-600 text-white font-bold rounded-xl hover:bg-blue-700 transition-all"
          >
            Kapat
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
import { jwtDecode } from 'jwt-decode';
import { useRoute, useRouter } from 'vue-router';
import axios from 'axios';

const route = useRoute();
const router = useRouter();
const artwork = ref(null);
const showArtistModal = ref(false);
const artistInfo = ref(null);



// Eser Satın AlmaFonksiyonu
const buyArtwork = async () => {
  const token = localStorage.getItem('userToken');
  if (!token) {
    alert("Lütfen önce giriş yapın! 👤");
    return;
  }

  if (!confirm(`${artwork.value.title} eserini satın almak istiyor musunuz?`)) return;

  try {
    const decoded = jwtDecode(token);
    await axios.post('http://localhost:8080/artworks/buy', {
      email: decoded.email,
      artworkId: artwork.value.id,
      price: artwork.value.price
    });

    alert("Satın alma başarılı! Siparişlerinize yönlendiriliyorsunuz.");
    router.push('/profile'); // Hemen profile gitsin ki görsün
  } catch (error) {
    if (error.response && error.response.status === 409) {
      alert("Maalesef bu eser çoktan satılmış! 😔");
    } else {
      alert("Satın alma sırasında bir hata oluştu.");
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
    artistInfo.value = {
      name: artistName,
      biography: '',
      artworksCount: 0,
      workshopsCount: 0
    };
  }
};

// Modal açıldığında sanatçı bilgisini yükle
watch(showArtistModal, (newVal) => {
  if (newVal && artwork.value) {
    fetchArtistInfo(artwork.value.artist);
  }
});




onMounted(async () => {
  try {
    const response = await axios.get('http://localhost:8080/artworks');
    const allArtworks = response.data;
    artwork.value = allArtworks.find(a => a.id === parseInt(route.params.id));
  } catch (error) {
    console.error("Detaylar yüklenemedi:", error);
  }

});
</script>