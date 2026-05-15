<template>
  <div class="bg-white rounded-3xl shadow-sm hover:shadow-xl border border-gray-100 overflow-hidden transition-all duration-300 group" @click="goToDetail">
    <div class="relative overflow-hidden h-64">
      <img :src="artwork.image || artwork.imageUrl" :alt="artwork.title" class="w-full h-full object-cover group-hover:scale-110 transition-all duration-500" />
      <button 
        @click.stop="toggleFavorite" 
        class="absolute top-4 right-4 p-3 rounded-full shadow-lg transition-all cursor-pointer"
        :class="isFavorite ? 'bg-red-500 text-white' : 'bg-white/90 text-gray-400 hover:text-red-500'"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" :fill="isFavorite ? 'currentColor' : 'none'" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
        </svg>
      </button>

      <div v-if="artwork.IsCampaign" class="absolute top-4 left-4 bg-red-500 text-white px-3 py-1 rounded-full text-xs font-black shadow-lg">
        FIRSAT: %{{ artwork.DiscountRate }} İNDİRİM
      </div>
    </div>

    <div class="p-6">
      <div class="flex justify-between items-start mb-2">
        <div>
          <h3 class="text-xl font-bold text-gray-800">{{ artwork.title }}</h3>
          <p class="text-gray-500 text-sm italic">{{ artwork.artist || artwork.instructorName }}</p>
        </div>
        <div class="bg-green-50 text-green-600 px-3 py-1 rounded-lg text-sm font-bold">
          4.8 ⭐
        </div>
      </div>

        <div class="mt-4 flex flex-col">
          
          <div v-if="artwork.IsCampaign == true || artwork.IsCampaign == 1">
            <span class="text-sm text-gray-400 line-through block">{{ artwork.price }} TL</span>
            <span class="text-2xl font-extrabold text-red-500">
              {{ (artwork.price * (1 - artwork.DiscountRate / 100)).toFixed(2) }} TL
            </span>
          </div>

          <span v-else class="text-2xl font-extrabold text-gray-900">
            {{ artwork.price }} TL
          </span>

        </div>

        <button class="bg-green-600 text-white px-4 py-2 rounded-xl font-semibold hover:bg-green-700 transition-all shadow-md cursor-pointer">
          Detaylar
        </button>
      </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import { jwtDecode } from 'jwt-decode';

const props = defineProps({
  artwork: {
    type: Object,
    required: true,
    default: () => ({
      Id: 0,
      title: '',
      price: 0,
      DiscountRate: 0,
      IsCampaign: false,
      imageUrl: ''
    })
  },
  // Favori durumu gibi ekstra veriler geliyorsa bunları tutmaya devam edebilirsin
  isFavorite: {
    type: Boolean,
    default: false
  }
});

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
  // Veritabanından gelen ID; büyük harf (Id) veya küçük harf (id) olabilir.
  // props içindeki artwork objesinden çekiyoruz.
  const artworkId = props.artwork.Id || props.artwork.id;
  
  if (artworkId) {
    router.push(`/artwork/${artworkId}`);
  } else {
    // Eğer hala undefined geliyorsa konsola yazdıralım ki sorunu görelim
    console.error("Hata: Eser ID'si bulunamadı!", props.artwork);
  }
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