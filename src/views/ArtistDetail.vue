<template>
  <div class="min-h-screen bg-gray-50 py-12 px-6">
    <div v-if="artistData" class="max-w-6xl mx-auto">
      
      <div class="bg-white rounded-3xl shadow-xl p-10 mb-10 border border-gray-100">
        <button @click="router.push('/')" class="mb-6 flex items-center text-blue-600 font-bold hover:underline">
          ⬅️ Ana Sayfaya Dön
        </button>
        <h1 class="text-6xl font-black text-gray-900 mb-4">{{ artistData.name }}</h1>
        <div class="max-w-3xl p-6 bg-blue-50 rounded-2xl border border-blue-100">
          <h3 class="text-xl font-bold text-blue-700 mb-2 uppercase">Biyografi</h3>
          <p class="text-gray-700 text-lg leading-relaxed italic">
            {{ artistData.info || 'Bu sanatçı hakkında biyografi bilgisi bulunamadı.' }}
          </p>
        </div>
      </div>

      <h2 class="text-3xl font-black text-gray-800 mb-8 flex items-center gap-3">
        🖼️ {{ artistData.name }} Eserleri
      </h2>
      
      <div v-if="artistArtworks.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <ArtworkCard 
          v-for="artwork in artistArtworks" 
          :key="artwork.id"
          :id="artwork.id"
          :title="artwork.title" 
          :artist="artwork.artist" 
          :price="artwork.price" 
          :image="artwork.image"
        />
      </div>
    </div>

    <div v-else class="flex flex-col items-center justify-center h-[50vh]">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mb-4"></div>
      <p class="text-xl text-gray-500">Sanatçı bilgileri yükleniyor...</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import axios from 'axios';
import ArtworkCard from '../components/ArtworkCard.vue'; // Yolun doğruluğundan emin ol

const route = useRoute();
const router = useRouter();
const allArtworks = ref([]);

onMounted(async () => {
  try {
    const response = await axios.get('http://localhost:8080/artworks');
    allArtworks.value = response.data || [];
  } catch (error) {
    console.error("Hata:", error);
  }
});

// Sanatçı bilgilerini filtrele
const artistData = computed(() => {
  if (allArtworks.value.length === 0) return null;
  const art = allArtworks.value.find(a => a.artist === route.params.name);
  if (art) {
    console.log("Bulunan Sanatçı Verisi:", art);
  }
  return art ? { name: art.artist, info: art.artistInfo } : null;
});

// Sanatçının eserlerini filtrele
const artistArtworks = computed(() => {
  return allArtworks.value.filter(a => a.artist === route.params.name);
});
</script>