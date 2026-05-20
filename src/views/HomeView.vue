<template>
  <div v-if="recommendedArtworks.length > 0" class="px-10 mt-10 mb-14 animate-in fade-in slide-in-from-top-4 duration-1000">
    <div class="relative bg-[#1a3a3a] rounded-[45px] p-10 shadow-[0_25px_50px_-12px_rgba(0,0,0,0.5)] overflow-hidden border border-white/5 flex flex-col lg:flex-row items-center gap-10">
      
      <div class="lg:w-1/3 z-10 text-left">
        <div class="flex items-center gap-2 mb-4">
          <span class="text-[#e9c46a] text-2xl animate-pulse">✦</span>
          <h2 class="text-[#e9c46a] text-4xl font-black tracking-tight leading-tight">
            {{ userLastCategory }} Tutkuna <br/> 
            <span class="text-white">Ekstra %5 İndirim!</span>
          </h2>
        </div>
        <p class="text-gray-300 text-lg leading-relaxed opacity-80">
          En son <strong>{{ userLastCategory }}</strong> eserleriyle ilgilendiğini fark ettik. <br/>
          Senin için seçtiğimiz bu özel seçkilerde şansını kaçırma!
        </p>
      </div>

      <div class="lg:w-2/3 w-full z-10">
        <div class="flex gap-6 overflow-x-auto pb-4 no-scrollbar">
          
          <div 
            v-for="art in recommendedArtworks" 
            :key="art.Id ?? art.id" 
            @click="router.push({ path: `/artwork/${art.Id ?? art.id}`, query: { extraDiscount: 'true' } })"
            class="min-w-[240px] group cursor-pointer"
          >
            <div class="bg-[#244a4a] p-3 rounded-[30px] border border-white/10 transition-all duration-500 group-hover:border-[#e9c46a]/50 group-hover:-translate-y-2 shadow-xl">
               <ArtworkCard 
                 :artwork="art" 
                 class="!shadow-none !bg-transparent border-none scale-95 group-hover:scale-100 transition-transform pointer-events-none" 
               />
            </div>
          </div>

        </div>
      </div>

      <div class="absolute top-0 right-0 w-1/2 h-full bg-gradient-to-l from-white/5 to-transparent pointer-events-none"></div>
      <div class="absolute -left-20 -bottom-20 w-80 h-80 bg-[#2a9d8f]/10 rounded-full blur-[100px] pointer-events-none"></div>
    </div>
  </div>

  <div class="min-h-screen bg-gray-50">
    <header class="bg-galeri-yesil text-white shadow-2xl rounded-b-[50px] p-8 pb-12">
      <div class="container mx-auto flex justify-between items-center">
        <div>
          <h1 class="text-4xl font-black italic tracking-tighter">Merhaba, {{ userName }} 🎨</h1>
          <p class="opacity-90 mt-2 text-lg">Bugün hangi sanat eserini keşfetmek istersin?</p>
        </div>
        <div class="flex items-center gap-6">
          <nav class="hidden md:flex gap-8 font-bold text-lg items-center">
            <span class="cursor-pointer border-b-4 border-white pb-1">Eserler</span>
            <span @click="router.push('/workshops')" class="cursor-pointer opacity-70 hover:opacity-100 transition-opacity">Atölyeler</span>
            <span @click="router.push('/artists')" class="cursor-pointer opacity-70 hover:opacity-100 transition-opacity">Sanatçılar</span>
            <span @click="router.push('/support')" class="cursor-pointer opacity-70 hover:opacity-100 transition-opacity">Destek</span>
            <span v-if="userRole === 'Admin'" @click="router.push('/admin/dashboard')" class="cursor-pointer opacity-70 hover:opacity-100 transition-opacity text-yellow-300">Yönetici</span>
          </nav>
          <div class="bg-white text-galeri-yesil p-3 rounded-2xl shadow-inner cursor-pointer hover:scale-105 transition-transform">
            <router-link to="/profile">
              <span class="text-2xl">👤</span>
            </router-link>
          </div>
        </div>
      </div>
    </header>

    <div class="mt-8 px-10 flex flex-wrap justify-between items-center gap-4">
        <button 
          @click="handleLogout" 
          class="flex items-center gap-2 px-5 bg-red-50 text-red-600 font-bold py-3 rounded-2xl hover:bg-red-100 transition-all border border-red-100 cursor-pointer"
        >
          <span>🚪</span> Çıkış Yap
        </button>

        <div v-if="userRole === 'Instructor'" class="flex gap-4">
          <button @click="router.push('/add-artwork')" class="bg-galeri-yesil text-white px-6 py-3 rounded-2xl font-bold shadow-lg shadow-green-100 hover:scale-105 transition-transform cursor-pointer">+ Eser Ekle</button>
          <button @click="router.push('/add-workshop')" class="bg-blue-600 text-white px-6 py-3 rounded-2xl font-bold shadow-lg shadow-blue-100 hover:scale-105 transition-transform cursor-pointer">+ Atölye Oluştur</button>
        </div>
    </div>

    <main class="container mx-auto px-4 pb-20">
      <div class="bg-white p-6 rounded-3xl shadow-sm mb-12 flex flex-col md:flex-row gap-6 items-center justify-between border border-gray-100 mt-6">
        <div class="flex flex-wrap items-center gap-4 w-full">
          <div class="relative w-full md:w-80">
            <input 
              v-model="searchQuery" 
              type="text" 
              placeholder="Eser ara..." 
              class="w-full p-4 pl-12 bg-gray-50 rounded-2xl focus:outline-none focus:ring-2 focus:ring-galeri-yesil/20 border border-gray-100" 
            />
            <span class="absolute left-4 top-4">🔍</span>
          </div>

          <select v-model="activeFilterType" class="p-4 bg-gray-50 text-gray-700 rounded-2xl font-bold border border-gray-100 outline-none focus:ring-2 focus:ring-galeri-yesil/20 min-w-[180px]">
            <option value="none">Filtreleme Yok</option>
            <option value="price">💰 Fiyata Göre</option>
            <option value="artist">👨‍🎨 Sanatçıya Göre</option>
            <option value="category">🖼️ Kategoriye Göre</option>
          </select>

          <div v-if="activeFilterType === 'price'" class="flex items-center bg-gray-50 p-1 rounded-2xl border border-gray-100 animate-in fade-in slide-in-from-left-2">
            <button @click="priceOrder = 'asc'" :class="priceOrder === 'asc' ? 'bg-galeri-yesil text-white shadow-md' : 'text-gray-500'" class="px-4 py-3 rounded-xl font-bold transition-all">📈 Artan</button>
            <button @click="priceOrder = 'desc'" :class="priceOrder === 'desc' ? 'bg-galeri-yesil text-white shadow-md' : 'text-gray-500'" class="px-4 py-3 rounded-xl font-bold transition-all">📉 Azalan</button>
          </div>

          <select v-if="activeFilterType === 'artist'" v-model="selectedArtist" class="p-4 bg-blue-50 text-blue-700 rounded-2xl font-bold border border-blue-100 outline-none">
            <option value="">Tüm Sanatçılar</option>
            <option v-for="artist in artistList" :key="artist" :value="artist">{{ artist }}</option>
          </select>

          <select v-if="activeFilterType === 'category'" v-model="selectedCategory" class="p-4 bg-gray-50 text-gray-700 rounded-2xl font-bold border border-gray-100 outline-none">
            <option value="">Tüm Kategoriler</option>
            <option v-for="cat in categories" :key="cat" :value="cat">{{ cat }}</option>
          </select>

          <button 
            @click="showOnlyCampaigns = !showOnlyCampaigns" 
            :class="showOnlyCampaigns ? 'bg-red-600 text-white' : 'bg-white text-red-600 border-red-100'"
            class="flex items-center gap-2 px-6 py-4 rounded-2xl border-2 font-black transition-all shadow-lg ml-auto"
          >
            <span>🔥</span> Kampanyalı Ürünler
          </button>
        </div>
      </div>

      <div v-if="filteredArtworks.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
        <ArtworkCard 
          v-for="art in filteredArtworks" 
          :key="art.Id ?? art.id"
          :artwork="art" 
        />
      </div>

      <div v-else class="text-center py-32 bg-white rounded-[40px] border border-dashed border-gray-200">
        <span class="text-7xl block mb-6">🏝️</span>
        <h3 class="text-2xl font-bold text-gray-800 mb-2">Eser Bulunamadı</h3>
        <button @click="resetFilters" class="mt-6 text-galeri-yesil font-bold hover:underline">Filtreleri Temizle</button>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue' 
import { useRouter } from 'vue-router'
import { jwtDecode } from 'jwt-decode'
import axios from 'axios'
import ArtworkCard from '../components/ArtworkCard.vue'

const router = useRouter()
const userName = ref('Misafir')
const userRole = ref('')
const artworks = ref([])

// --- FİLTRELEME STATE'LERİ ---
const searchQuery = ref('')
const activeFilterType = ref('none') 
const priceOrder = ref('asc')        
const selectedArtist = ref('')
const selectedCategory = ref('')
const showOnlyCampaigns = ref(false)

const categories = ref(['Manzara', 'Rönesans', 'Realizm', 'Natürmort', 'Modern Sanat'])

// --- 1. KULLANICI TERCİHİ (Duplicate silindi) ---
// --- 1. KULLANICI TERCİHİ VE ÖNERİ SİSTEMİ ---

// Bu computed zaten dinamik çalışıyor: en son ne alındıysa o gelir.
const userLastCategory = computed(() => {
  return localStorage.getItem('lastPurchasedCategory') || '';
});

// Üstteki "Kişiye Özel Seçkiler" Bandı
// --- 1. DİNAMİK ÖNERİ ALGORİTMASI (Gelişmiş Versiyon) ---

const recommendedArtworks = computed(() => {
  if (!userLastCategory.value || !artworks.value || artworks.value.length === 0) return [];

  return artworks.value
    .filter(a => {
      // 1. Kategori eşleşmeli
      const cat = a.category || a.Category;
      const isSameCategory = cat === userLastCategory.value;
      
      // 2. Satılanlar da anasayfada görünsün (müşteri isteği)
      return isSameCategory;
    })
    .map(a => {
      // 🔥 ŞOV BURADA: Her esere o kişiye özel %5 ekstra indirim tanımlıyoruz
      // Mevcut bir indirimi varsa onun üzerine değil, son fiyat üzerinden %5 daha düşüyoruz
      const basePrice = a.price;
      const currentDiscount = a.DiscountRate || a.discountrate || 0;
      const priceAfterFirstDiscount = basePrice * (1 - currentDiscount / 100);
      
      // Kişiye özel "Sadakat İndirimi" (%5)
      const finalSpecialPrice = priceAfterFirstDiscount * 0.95;

      return {
        ...a,
        specialPrice: finalSpecialPrice, // Template'de bunu kullanacağız
        hasSpecialOffer: true
      };
    })
    .slice(0, 4);
});
// --- 2. ANA GALERİ FİLTRELEME MANTIĞI ---
const filteredArtworks = computed(() => {
  if (!artworks.value || artworks.value.length === 0) return [];
  
  // ADIM 1: Tüm eserleri alıyoruz (Satılanlar da anasayfada görünsün)
  let result = [...artworks.value];

  // ADIM 2: Arama
  const search = searchQuery.value?.toLowerCase().trim();
  if (search) {
    result = result.filter(artwork => 
      artwork.title?.toLowerCase().includes(search) || 
      artwork.artist?.toLowerCase().includes(search)
    );
  }

  // ADIM 3: Kampanya
  if (showOnlyCampaigns.value) {
    result = result.filter(a => {
      const isPromo = a.iscampaign !== undefined ? a.iscampaign : a.IsCampaign;
      return isPromo == true || isPromo == 1;
    });
  }

  // ADIM 4: Kategori/Sanatçı
  if (activeFilterType.value === 'artist' && selectedArtist.value) {
    result = result.filter(a => (a.artist || a.Artist) === selectedArtist.value);
  }
  if (activeFilterType.value === 'category' && selectedCategory.value) {
    result = result.filter(a => (a.category || a.Category) === selectedCategory.value);
  }

  // ADIM 5: Fiyat Sıralama
  if (activeFilterType.value === 'price') {
    result.sort((a, b) => {
      const getFinalPrice = (item) => {
        const isPromo = item.iscampaign !== undefined ? item.iscampaign : item.IsCampaign;
        const discount = item.discountrate !== undefined ? item.discountrate : item.DiscountRate;
        return (isPromo == true || isPromo == 1) 
          ? item.price * (1 - (discount || 0) / 100) 
          : item.price;
      };
      const priceA = getFinalPrice(a);
      const priceB = getFinalPrice(b);
      return priceOrder.value === 'asc' ? priceA - priceB : priceB - priceA;
    });
  }
  return result;
});

const fetchArtworks = async () => {
  try {
    const response = await axios.get(`http://localhost:8080/artworks?t=${new Date().getTime()}`);
    artworks.value = response.data;
  } catch (error) {
    console.error("Eserler yüklenirken hata:", error);
  }
};

const artistList = computed(() => {
  const artists = artworks.value.map(a => a.artist || a.Artist)
  return [...new Set(artists)].filter(a => a) 
})

const resetFilters = () => {
  searchQuery.value = ''
  activeFilterType.value = 'none'
  selectedCategory.value = ''
  selectedArtist.value = ''
  showOnlyCampaigns.value = false
  priceOrder.value = 'asc'
}

onMounted(() => {
  const token = localStorage.getItem('userToken')
  if (token) {
    try {
      const decoded = jwtDecode(token)
      const first = decoded.firstName || ''
      const last = decoded.lastName || ''
      if (first || last) {
        userName.value = `${first}${last ? ' ' + last : ''}`
      } else {
        // token doesn't include names; try fetching profile by email stored in localStorage
        const email = localStorage.getItem('userEmail')
        if (email) {
          axios.get(`http://localhost:8080/profile?email=${encodeURIComponent(email)}`)
            .then(res => {
              const u = res.data || {}
              userName.value = (u.firstName || u.FirstName || '') || (u.email || 'Misafir')
              if (u.lastName || u.LastName) {
                userName.value = `${u.firstName || u.FirstName}${u.lastName || u.LastName ? ' ' + (u.lastName || u.LastName) : ''}`
              }
            })
            .catch(err => { console.error('Profil çekilemedi:', err) })
        }
      }
    } catch (error) {
      console.error('Token çözülemedi:', error)
    }
  }
  userRole.value = localStorage.getItem('userRole')
  fetchArtworks()
})

const handleLogout = () => {
  localStorage.clear()
  router.push('/login')
}
</script>