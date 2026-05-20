<template>
  <div 
    class="bg-white rounded-3xl shadow-sm hover:shadow-xl border border-gray-100 overflow-hidden transition-all duration-300 group relative cursor-pointer" 
    :class="{'opacity-90': isSold}"
    @click="goToDetail()"
  >
    <div v-if="isSold" class="absolute inset-0 z-30 flex items-center justify-center pointer-events-none">
      <div class="rotate-[-15deg] border-8 border-red-600/80 px-8 py-3 rounded-2xl bg-white/10 backdrop-blur-sm shadow-2xl">
        <span class="text-red-600 text-5xl font-black tracking-tighter uppercase italic">SATILDI</span>
      </div>
    </div>

    <div class="relative overflow-hidden h-64">
      <img 
        :src="artwork.image || artwork.imageUrl" 
        :alt="artwork.title" 
        class="w-full h-full object-cover transition-all duration-500" 
        :class="isSold ? 'grayscale blur-[1px] opacity-60' : 'group-hover:scale-110'"
      />
      
      <div v-if="!isSold" class="absolute top-4 right-4 flex flex-col gap-2 z-40">
        <button 
          @click.stop="toggleFavorite" 
          class="p-3 rounded-full shadow-lg transition-all cursor-pointer"
          :class="isFavorite ? 'bg-red-500 text-white' : 'bg-white/90 text-gray-400 hover:text-red-500'"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" :fill="isFavorite ? 'currentColor' : 'none'" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
          </svg>
        </button>

        <button 
          @click.stop="addToCompare" 
          class="p-3 bg-white/90 rounded-full shadow-lg transition-all cursor-pointer text-gray-400"
          :class="isInCompareList ? 'bg-blue-600 text-white hover:bg-blue-700' : 'hover:bg-blue-600 hover:text-white'"
          title="Karşılaştır"
        >
          ⚖️
        </button>
      </div>

      <div v-if="!isSold && !artwork.specialPrice && (artwork.IsCampaign || artwork.iscampaign)" class="absolute top-4 left-4 bg-red-500 text-white px-3 py-1 rounded-full text-xs font-black shadow-lg z-10">
        FIRSAT: %{{ artwork.DiscountRate || artwork.discountrate }} İNDİRİM
      </div>

      <div v-if="!isSold && artwork.specialPrice" class="absolute top-4 left-4 bg-[#e9c46a] text-[#1a3a3a] px-3 py-1 rounded-full text-xs font-black shadow-lg z-10 animate-bounce">
        SANA ÖZEL %5 EKSTRA!
      </div>
    </div>

    <div class="p-6" :class="{'opacity-50': isSold}">
      <div class="flex justify-between items-start mb-2">
        <div>
          <h3 class="text-xl font-bold text-gray-800">{{ artwork.title }}</h3>
          <p class="text-gray-500 text-sm italic">{{ artwork.artist || artwork.instructorName || 'Bilinmeyen Sanatçı' }}</p>
        </div>
        <div class="bg-green-50 text-green-600 px-3 py-1 rounded-lg text-sm font-bold">
          {{ rating }} ⭐
        </div>
      </div>

      <div class="mt-4 flex flex-col mb-4">
        
        <div v-if="artwork.specialPrice && !isSold">
          <span class="text-sm text-gray-400 line-through block">{{ artwork.price }} TL</span>
          <div class="flex items-center gap-2">
            <span class="text-2xl font-extrabold text-[#d4a373]">
              {{ artwork.specialPrice.toFixed(2) }} TL
            </span>
            <span class="bg-[#e9c46a]/20 text-[#b08d3e] text-[10px] px-2 py-1 rounded-md font-bold border border-[#e9c46a]/30">
              SANA ÖZEL
            </span>
          </div>
        </div>

        <div v-else-if="(artwork.IsCampaign || artwork.iscampaign) && !isSold">
          <span class="text-sm text-gray-400 line-through block">{{ artwork.price }} TL</span>
          <span class="text-2xl font-extrabold text-red-500">
            {{ (artwork.price * (1 - (artwork.DiscountRate || artwork.discountrate) / 100)).toFixed(2) }} TL
          </span>
        </div>

        <span v-else class="text-2xl font-extrabold text-gray-900">
          {{ artwork.price }} TL
        </span>
      </div>

      <button 
        @click.stop="goToDetail()"
        :class="isSold ? 'bg-gray-500 hover:bg-gray-600 cursor-pointer' : 'bg-green-600 hover:bg-green-700 cursor-pointer'"
        class="w-full text-white px-4 py-3 rounded-xl font-bold transition-all shadow-md"
      >
        {{ isSold ? 'SATILDI - Detayları Gör' : 'Detayları Gör' }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import { jwtDecode } from 'jwt-decode';

const props = defineProps({
  artwork: {
    type: Object,
    required: true,
    default: () => ({
      Id: 0,
      id: 0,
      title: '',
      price: 0,
      DiscountRate: 0,
      discountrate: 0,
      IsCampaign: false,
      iscampaign: false,
      imageUrl: '',
      image: '',
      issold: false,
      IsSold: false,
      artist: '',
      instructorName: ''
    })
  }
});

const router = useRouter();
const isFavorite = ref(false);
const rating = ref(null); // will be set from entity-stats or artwork prop
const compareList = ref([]);

const isSold = computed(() => {
  const s = props.artwork.issold !== undefined ? props.artwork.issold : 
            props.artwork.IsSold !== undefined ? props.artwork.IsSold : 
            props.artwork.is_sold !== undefined ? props.artwork.is_sold : 
            props.artwork.Is_sold !== undefined ? props.artwork.Is_sold : false;
  return s == 1 || s == true || s == "1";
});

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

const loadCompareList = () => {
  const list = localStorage.getItem('compareList');
  compareList.value = list ? JSON.parse(list) : [];
};

const isInCompareList = computed(() => {
  const artworkId = props.artwork.Id ?? props.artwork.id;
  return compareList.value.includes(artworkId);
});

onMounted(async () => {
  loadCompareList();
  const email = getEmailFromToken();
  const artworkId = props.artwork.Id ?? props.artwork.id;
  
  if (email && artworkId !== undefined && artworkId !== null) {
    try {
      const response = await axios.get('http://localhost:8080/favorites/check', {
        params: { email: email, artworkId: artworkId }
      });
      isFavorite.value = response.data.isFavorite;
    } catch (error) {
      console.error("Favori durumu kontrol hatası:", error);
    }
  }

  // Fetch avg rating; fall back to artwork-provided rating or 0.0
  if (artworkId !== undefined && artworkId !== null) {
    try {
      const statsRes = await axios.get(`http://localhost:8080/entity-stats?targetId=${artworkId}&targetType=Artwork`);
      if (statsRes.data && statsRes.data.avgRating !== undefined && statsRes.data.avgRating > 0) {
        rating.value = Number(statsRes.data.avgRating).toFixed(1);
      } else if (props.artwork.Rating || props.artwork.rating) {
        rating.value = Number(props.artwork.Rating || props.artwork.rating).toFixed(1);
      } else {
        rating.value = Number(0).toFixed(1);
      }
    } catch (err) {
      // on error, try fallback
      if (props.artwork.Rating || props.artwork.rating) {
        rating.value = Number(props.artwork.Rating || props.artwork.rating).toFixed(1);
      } else {
        rating.value = Number(0).toFixed(1);
      }
      console.error("Rating çekilemedi kanka:", err);
    }
  } else {
    if (props.artwork.Rating || props.artwork.rating) {
      rating.value = Number(props.artwork.Rating || props.artwork.rating).toFixed(1);
    } else {
      rating.value = Number(0).toFixed(1);
    }
  }
});

const goToDetail = () => {
  const artworkId = props.artwork.Id ?? props.artwork.id;
  if (artworkId !== undefined && artworkId !== null) {
    router.push(`/artwork/${artworkId}`);
  } else {
    console.error("Hata: Eser ID'si bulunamadı!", props.artwork);
  }
};

// 🚀 ENES'İN SEPETE EKLEME TERAZİSİ (Senin paket nesne mimarinle harmanlandı)
const addToCompare = () => {
  const artworkId = props.artwork.Id ?? props.artwork.id;
  if (artworkId === undefined || artworkId === null) return;

  const currentList = JSON.parse(localStorage.getItem('compareList') || '[]');
  const currentType = localStorage.getItem('compareType') || 'Artwork';

  if (currentType !== 'Artwork') {
    if (!confirm('Karşılaştırma listenizdeki atölyeler temizlenecek. Devam edilsin mi kanka?')) return;
    localStorage.setItem('compareList', JSON.stringify([artworkId]));
    localStorage.setItem('compareType', 'Artwork');
  } else {
    if (!currentList.includes(artworkId)) {
      if (currentList.length >= 4) {
        alert('Kanka yavaş, en fazla 4 ürünü karşılaştırabilirsin! 🛑');
        return;
      }
      currentList.push(artworkId);
      localStorage.setItem('compareList', JSON.stringify(currentList));
      localStorage.setItem('compareType', 'Artwork');
    } else {
      // Çift tıklamada listeden kaldır kanka (Toggle güvencesi)
      const filtered = currentList.filter(item => item !== artworkId);
      localStorage.setItem('compareList', JSON.stringify(filtered));
    }
  }
  router.push('/compare');
};

const toggleFavorite = async () => {
  const email = getEmailFromToken();
  const artworkId = props.artwork.Id ?? props.artwork.id;

  if (!email) {
    alert("Lütfen önce giriş yapın! 👤");
    return;
  }

  try {
    if (isFavorite.value) {
      await axios.post('http://localhost:8080/favorites/remove', { email: email, artworkId: artworkId });
      isFavorite.value = false;
    } else {
      await axios.post('http://localhost:8080/favorites/add', { email: email, artworkId: artworkId });
      isFavorite.value = true;
    }
  } catch (error) {
    alert("İşlem sırasında bir hata oluştu.");
  }
};
</script>