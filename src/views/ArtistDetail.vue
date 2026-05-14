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
            {{ artistData.biography || 'Bu sanatçı hakkında biyografi bilgisi bulunamadı.' }}
          </p>
          <div class="mt-6 grid grid-cols-2 gap-4">
            <div class="bg-white rounded-2xl p-4 border border-gray-200">
              <p class="text-sm uppercase tracking-[0.2em] text-gray-400 mb-2">Toplam Eser</p>
              <p class="text-3xl font-black text-blue-700">{{ artistData.artworksCount }}</p>
            </div>
            <div class="bg-white rounded-2xl p-4 border border-gray-200">
              <p class="text-sm uppercase tracking-[0.2em] text-gray-400 mb-2">Toplam Atölye</p>
              <p class="text-3xl font-black text-blue-700">{{ artistData.workshopsCount }}</p>
            </div>
          </div>
        </div>
      </div>

      <section id="artworks" class="mb-14">
        <div class="flex items-center justify-between mb-8">
          <h2 class="text-3xl font-black text-gray-800 flex items-center gap-3">🖼️ {{ artistData.name }} Eserleri</h2>
          <span class="text-sm text-gray-500">{{ artistArtworks.length }} eser bulundu</span>
        </div>

        <div v-if="artistArtworks.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          <ArtworkCard 
            v-for="artwork in artistArtworks" 
            :key="artwork.id"
            :id="artwork.id"
            :title="artwork.title" 
            :artist="artwork.artist" 
            :price="artwork.price" 
            :image="artwork.imageUrl"
            :rating="artwork.averageRating"
            :commentCount="artwork.commentCount"
          />
        </div>

        <div v-else class="text-center py-16 rounded-3xl bg-white border border-gray-200">
          <p class="text-gray-500">Henüz bu sanatçıya ait bir eser bulunmuyor.</p>
        </div>
      </section>

      <section id="workshops" class="mb-14">
        <div class="flex items-center justify-between mb-8">
          <h2 class="text-3xl font-black text-gray-800 flex items-center gap-3">🎨 {{ artistData.name }} Atölyeleri</h2>
          <span class="text-sm text-gray-500">{{ artistWorkshops.length }} atölye bulundu</span>
        </div>

        <div v-if="artistWorkshops.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          <WorkshopCard
            v-for="workshop in artistWorkshops"
            :key="workshop.id"
            :id="workshop.id"
            :title="workshop.title"
            :instructorName="workshop.instructorName"
            :availableDates="workshop.availableDates"
            :location="workshop.location"
            :capacity="workshop.capacity"
            :price="workshop.price"
            :image="workshop.image"
          />
        </div>

        <div v-else class="text-center py-16 rounded-3xl bg-white border border-gray-200">
          <p class="text-gray-500">Henüz bu sanatçıya ait bir atölye bulunmuyor.</p>
        </div>
      </section>
    </div>

    <div v-else class="flex flex-col items-center justify-center h-[50vh]">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mb-4"></div>
      <p class="text-xl text-gray-500">Sanatçı bilgileri yükleniyor...</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch, nextTick } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import axios from 'axios';
import ArtworkCard from '../components/ArtworkCard.vue';
import WorkshopCard from '../components/WorkshopCard.vue';

const route = useRoute();
const router = useRouter();
const artistData = ref(null);
const allArtworks = ref([]);
const allWorkshops = ref([]);
const activeTab = ref(route.query.tab || 'artworks');

const scrollToSection = async (section) => {
  await nextTick();
  const el = document.getElementById(section);
  if (el) {
    el.scrollIntoView({ behavior: 'smooth', block: 'start' });
  }
};

watch(
  () => route.query.tab,
  (value) => {
    activeTab.value = value || 'artworks';
    scrollToSection(activeTab.value);
  }
);

onMounted(async () => {
  try {
    const [artistRes, artworksRes, workshopsRes] = await Promise.all([
      axios.get(`http://localhost:8080/artist?name=${encodeURIComponent(route.params.name)}`),
      axios.get('http://localhost:8080/artworks'),
      axios.get('http://localhost:8080/workshops')
    ]);

    artistData.value = artistRes.data;
    allArtworks.value = artworksRes.data || [];
    allWorkshops.value = workshopsRes.data || [];

    scrollToSection(activeTab.value);
  } catch (error) {
    console.error('Sanatçı detayları yüklenemedi:', error);
  }
});

const artistArtworks = computed(() => {
  return allArtworks.value.filter(a => a.artist === route.params.name);
});

const artistWorkshops = computed(() => {
  return allWorkshops.value.filter(w => w.instructorName === route.params.name);
});
</script>