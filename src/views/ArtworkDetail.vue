<template>
  <div class="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div v-if="artwork" class="max-w-6xl mx-auto bg-white rounded-3xl shadow-2xl overflow-hidden border border-gray-100">
      <div class="flex flex-col lg:flex-row">
        
        <!-- Sol Taraf: Büyük Eser Görseli -->
        <div class="lg:w-1/2 relative bg-gray-200">
          <img :src="artwork.image" :alt="artwork.title" class="w-full h-full object-cover min-h-[500px]" />
          <button @click="router.back()" class="absolute top-6 left-6 bg-white/90 p-3 rounded-full shadow-lg hover:bg-white transition-all">
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
                <button @click="router.push('/artist/' + artwork.artist)"class="text-sm bg-blue-100 text-blue-700 px-4 py-2 rounded-xl font-bold hover:bg-blue-600 hover:text-white transition-all shadow-sm">
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

    <!-- Yorumlar Bölümü -->
    <div v-if="artwork" class="max-w-6xl mx-auto mt-8">
      <CommentSection targetType="Artwork" :targetId="artwork.id" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { jwtDecode } from 'jwt-decode';
import { useRoute, useRouter } from 'vue-router';
import axios from 'axios';
import CommentSection from '../components/CommentSection.vue';

const route = useRoute();
const router = useRouter();
const artwork = ref(null);



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




onMounted(async () => {
  try {
    // Backend'den tüm eserleri getirip içinden ID'si eşleşeni buluyoruz
    // Not: İleride sadece tek bir eseri çeken /artworks/:id endpoint'i de yazılabilir
    const response = await axios.get('http://localhost:8080/artworks');
    const allArtworks = response.data;
    artwork.value = allArtworks.find(a => a.id === parseInt(route.params.id));
  } catch (error) {
    console.error("Detaylar yüklenemedi:", error);
  }

});
</script>