<template>
  <div class="bg-white rounded-3xl shadow-sm hover:shadow-xl border border-gray-100 overflow-hidden transition-all duration-300 group" @click="goToDetail">
    <!-- Eser Görseli -->
    <div class="relative overflow-hidden h-64">
      <img :src="image" :alt="title" class="w-full h-full object-cover group-hover:scale-110 transition-duration-500" />
      <button 
        @click.stop="toggleFavorite" 
        class="absolute top-4 right-4 p-3 rounded-full shadow-lg transition-all cursor-pointer"
        :class="isFavorite ? 'bg-red-500 text-white' : 'bg-white/90 text-gray-400 hover:text-red-500'"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" :fill="isFavorite ? 'currentColor' : 'none'" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
        </svg>
      </button>
    </div>
    <!-- Eser Bilgileri -->
    <div class="p-6">
      <div class="flex justify-between items-start mb-2">
        <div>
          <h3 class="text-xl font-bold text-gray-800">{{ title }}</h3>
          <p class="text-gray-500 text-sm italic">{{ artist }}</p>
        </div>
        <div class="bg-green-50 text-galeri-yesil px-3 py-1 rounded-lg text-sm font-bold">
          4.8 ⭐
        </div>
      </div>

      <div class="mt-4 flex items-center justify-between">
        <span class="text-2xl font-extrabold text-gray-900">{{ price }} TL</span>
        <button class="bg-galeri-yesil text-white px-4 py-2 rounded-xl font-semibold hover:bg-green-600 transition-all shadow-md cursor-pointer">
          Detaylar
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import { jwtDecode } from 'jwt-decode';

const props = defineProps(['id', 'title', 'artist', 'price', 'image']);
const router = useRouter();
const isFavorite = ref(false);

const getEmailFromToken = () => {
  const token = localStorage.getItem('userToken');
  if (!token) return null;
  try {
    const decoded = jwtDecode(token);
    return decoded.email;
  } catch (e) {
    return null;
  }
};

onMounted(async () => {
  const email = getEmailFromToken();
  if (email && props.id) {
    try {
      const response = await axios.get('http://localhost:8080/favorites/check', {
        params: {
          email: email,
          artworkId: props.id
        }
      });
      
      // Go'dan gelen {"isFavorite": true/false} cevabını değişkene atıyoruz
      isFavorite.value = response.data.isFavorite;
    } catch (error) {
      console.error("Favori durumu kontrol hatası:", error);
    }
  }
});

const goToDetail = () => {
  router.push(`/artwork/${props.id}`);
};

const toggleFavorite = async () => {
  const email = getEmailFromToken();
  if (!email) {
    alert("Lütfen önce giriş yapın! 👤");
    return;
  }

  try {
    if (isFavorite.value) {
      // Favoriden ÇIKAR
      await axios.post('http://localhost:8080/favorites/remove', {
        email: email,
        artworkId: props.id
      });
      isFavorite.value = false;
    } else {
      // Favoriye EKLE
      await axios.post('http://localhost:8080/favorites/add', {
        email: email,
        artworkId: props.id
      });
      isFavorite.value = true;
    }
  } catch (error) {
    alert("İşlem sırasında bir hata oluştu.");
  }
};
</script>