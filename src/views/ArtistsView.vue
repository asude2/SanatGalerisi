<template>
  <div class="min-h-screen bg-gray-50 p-8">
    <div class="container mx-auto">
      
      <button 
        @click="goBack" 
        class="mb-8 p-3 bg-white hover:bg-gray-100 rounded-full shadow-sm border border-gray-100 transition-all group flex items-center justify-center w-12 h-12"
        title="Geri Git"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-gray-600 group-hover:-translate-x-1 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
        </svg>
      </button>

      <header class="mb-12">
        <h1 class="text-4xl font-black italic tracking-tighter text-gray-900">Sanatçılarımız 👩‍🎨</h1>
        <p class="text-gray-500 mt-2 text-lg">Platformumuzdaki değerli sanatçıları ve hikayelerini keşfedin.</p>
      </header>

      <div class="mb-10">
        <label class="block text-sm font-semibold text-gray-600 mb-2">Sanatçı Ara</label>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="İsim veya biyografi ile arayın..."
          class="w-full p-4 rounded-2xl border border-gray-200 bg-white shadow-sm focus:border-blue-400 focus:outline-none"
        />
      </div>

      <div v-if="filteredArtists.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
        <ArtistCard 
          v-for="artist in filteredArtists" 
          :key="artist.id"
          :name="artist.name"
          :biography="artist.biography"
          :nationality="artist.nationality"
          @view-artworks="goToArtistDetail"
          @view-workshops="goToArtistWorkshops"
        />
      </div>
      <div v-else class="text-center py-20">
        <p class="text-gray-400">Aradığınız kriterlere uygun bir sanatçı bulunamadı.</p>
      </div>
      
 
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router'; // Router'ı ekledik
import axios from 'axios';
import ArtistCard from '../components/ArtistCard.vue';

const router = useRouter(); // Router kullanımı için tanımladık
const artists = ref([]);
const searchQuery = ref('');

// Geri gitme fonksiyonu
const goBack = () => {
  router.back();
};

const goToArtistDetail = (name) => {
  router.push({ name: 'ArtistDetail', params: { name } });
};

const goToArtistWorkshops = (name) => {
  router.push({ name: 'ArtistDetail', params: { name }, query: { tab: 'workshops' } });
};

const filteredArtists = computed(() => {
  const query = String(searchQuery.value).trim().toLowerCase();
  if (!query) return artists.value;

  return artists.value.filter(artist => {
    return (
      String(artist.name || '').toLowerCase().includes(query) ||
      String(artist.biography || '').toLowerCase().includes(query) ||
      String(artist.nationality || '').toLowerCase().includes(query)
    );
  });
});

onMounted(async () => {
  try {
    const res = await axios.get('http://localhost:8080/artists');
    artists.value = res.data || [];
  } catch (error) {
    console.error("Sanatçılar yüklenemedi:", error);
  }
});
</script>