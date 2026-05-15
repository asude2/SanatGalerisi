<template>
  <div class="min-h-screen bg-gray-50">
    <header class="bg-galeri-yesil text-white shadow-2xl rounded-b-[50px] p-8 pb-12">
      <div class="container mx-auto flex justify-between items-center">
        <div>
          <h1 class="text-4xl font-black italic tracking-tighter">Merhaba, {{ userName }} 🎨</h1>
          <p class="opacity-90 mt-2 text-lg">Bugün hangi sanat eserini keşfetmek istersin?</p>
        </div>
        <div class="flex items-center gap-6">
          <nav class="hidden md:flex gap-8 font-bold text-lg">
            <span class="cursor-pointer border-b-4 border-white pb-1">Eserler</span>
            <span @click="router.push('/workshops')" class="cursor-pointer opacity-70 hover:opacity-100 transition-opacity">Atölyeler</span>
            <span @click="router.push('/artists')" class="cursor-pointer opacity-70 hover:opacity-100 transition-opacity">Sanatçılar</span>
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
          <button 
            @click="router.push('/add-artwork')" 
            class="bg-galeri-yesil text-white px-6 py-3 rounded-2xl font-bold shadow-lg shadow-green-100 hover:scale-105 transition-transform cursor-pointer"
          >
            + Eser Ekle
          </button>
          <button 
            @click="router.push('/add-workshop')" 
            class="bg-blue-600 text-white px-6 py-3 rounded-2xl font-bold shadow-lg shadow-blue-100 hover:scale-105 transition-transform cursor-pointer"
          >
            + Atölye Oluştur
          </button>
        </div>
    </div>

    <main class="container mx-auto -mt-4 px-4 pb-20">
      <div class="bg-white p-6 rounded-2xl shadow-sm mb-12 flex flex-col md:flex-row gap-4 items-center justify-between border border-gray-100 mt-6">
        
        <div class="relative w-full md:w-96">
          <input 
            v-model="searchQuery" 
            type="text" 
            placeholder="Eser veya sanatçı ara..." 
            class="w-full p-4 pl-12 bg-gray-50 rounded-2xl focus:outline-none focus:ring-2 focus:ring-galeri-yesil/20 border border-transparent focus:border-galeri-yesil/30 transition-all" 
          />
          <span class="absolute left-4 top-4">🔍</span>
        </div>
        

        <div class="relative">
          <button 
            @click.stop="showCategories = !showCategories"
            class="px-6 py-4 bg-gray-50 text-gray-700 rounded-2xl font-bold hover:bg-gray-100 border border-gray-100 transition-all cursor-pointer flex items-center gap-2 min-w-[160px] justify-between"
            :class="{ 'ring-2 ring-galeri-yesil ring-offset-1': selectedCategory }"
          >
            {{ selectedCategory || 'Kategoriler' }}
            <span class="text-xs transition-transform" :class="{ 'rotate-180': showCategories }">▼</span>
          </button>

          <div v-if="showCategories" class="absolute top-full right-0 mt-3 bg-white rounded-2xl shadow-2xl border border-gray-50 z-[100] w-64 overflow-hidden">
            <div class="p-2">
              <button 
                @click="selectCategory(null)"
                class="w-full text-left px-4 py-3 rounded-xl transition-all mb-1 font-medium"
                :class="!selectedCategory ? 'bg-galeri-yesil text-white' : 'hover:bg-gray-50 text-gray-600'"
              >
                Tüm Eserler
              </button>
              <button 
                v-for="category in categories"
                :key="category"
                @click="selectCategory(category)"
                class="w-full text-left px-4 py-3 rounded-xl transition-all mb-1 font-medium capitalize"
                :class="selectedCategory === category ? 'bg-galeri-yesil text-white' : 'hover:bg-gray-50 text-gray-600'"
              >
                {{ category }}
              </button>
            </div>
          </div>
        </div>


                <div class="flex items-center gap-4 mt-6">
          <div class="flex items-center bg-white border-2 border-gray-100 rounded-2xl p-1 shadow-sm">
            <button 
              @click="sortOrder = 'asc'" 
              :class="sortOrder === 'asc' ? 'bg-galeri-yesil text-white' : 'text-gray-500 hover:bg-gray-50'"
              class="px-4 py-2 rounded-xl font-bold transition-all flex items-center gap-2"
            >
              <span>📈</span> Artan Sırada
            </button>
            <button 
              @click="sortOrder = 'desc'" 
              :class="sortOrder === 'desc' ? 'bg-galeri-yesil text-white' : 'text-gray-500 hover:bg-gray-50'"
              class="px-4 py-2 rounded-xl font-bold transition-all flex items-center gap-2"
            >
              <span>📉</span> Azalan Sırada
            </button>
          </div>

          <button 
            v-if="sortOrder" 
            @click="sortOrder = null" 
            class="text-sm text-red-500 font-bold hover:bg-red-50 px-2 py-1 rounded-lg transition-colors"
          >
            Seçimi Temizle
          </button>
        </div>
      </div>

        <div v-if="filteredArtworks.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
          <ArtworkCard 
            v-for="art in filteredArtworks" 
            :key="art.Id || art.id"
            :artwork="art" 
          />
        </div>

        
      <div v-else class="text-center py-32 bg-white rounded-[40px] border border-dashed border-gray-200">
        <span class="text-7xl block mb-6">🏝️</span>
        <h3 class="text-2xl font-bold text-gray-800 mb-2">Sonuç Bulunamadı</h3>
        <p class="text-gray-400 text-lg max-w-md mx-auto">
          {{ selectedCategory ? `"${selectedCategory}" kategorisinde` : 'Aradığınız kriterlerde' }} 
          eşleşen bir eser bulamadık. Lütfen aramayı değiştirmeyi deneyin.
        </p>
        <button @click="resetFilters" class="mt-6 text-galeri-yesil font-bold hover:underline">
          Filtreleri Temizle
        </button>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, onUnmounted } from 'vue' 
import { useRouter } from 'vue-router'
import { jwtDecode } from 'jwt-decode'
import axios from 'axios'
import ArtworkCard from '../components/ArtworkCard.vue'

const router = useRouter()
const userName = ref('Misafir')
const userRole = ref('')
const artworks = ref([])
const searchQuery = ref('')
const selectedCategory = ref(null)
const showCategories = ref(false)

const categories = ref(['Manzara', 'Rönesans', 'Realizm', 'Natürmort', 'Modern Sanat'])


const fetchArtworks = async () => {
  try {
    const response = await axios.get('http://localhost:8080/artworks')
    artworks.value = response.data || [] 
  } catch (error) {
    console.error("Eserler yüklenirken bir hata oluştu:", error)
  }
}

const selectCategory = (category) => {
  selectedCategory.value = category
  showCategories.value = false
}

const resetFilters = () => {
  searchQuery.value = ''
  selectedCategory.value = null
}

// Menü dışına tıklayınca kapatma
const closeMenu = () => { showCategories.value = false }
onMounted(() => { window.addEventListener('click', closeMenu) })
onUnmounted(() => { window.removeEventListener('click', closeMenu) })

const sortOrder = ref(null) // 'asc', 'desc' veya null olabilir

const filteredArtworks = computed(() => {
  if (!artworks.value) return [] 
  let result = [...artworks.value] // Orijinal veriyi bozmamak için kopyasını alıyoruz

  // 1. ARAMA FİLTRESİ
  const search = searchQuery.value.toLowerCase().trim()
  if (search) {
    result = result.filter(artwork => 
      artwork.title?.toLowerCase().includes(search) || 
      artwork.artist?.toLowerCase().includes(search)
    )
  }

  // 2. KATEGORİ FİLTRESİ
  if (selectedCategory.value) {
    result = result.filter(artwork => artwork.category === selectedCategory.value)
  }

  // 3. SIRALAMA ŞOVU (Artan veya Azalan)
  if (sortOrder.value === 'asc') {
    // Küçükten büyüğe (Ucuzdan Pahalıya)
    result.sort((a, b) => a.price - b.price)
  } else if (sortOrder.value === 'desc') {
    // Büyükten küçüğe (Pahalıdan Ucuza)
    result.sort((a, b) => b.price - a.price)
  }

  return result
})

onMounted(() => {
  const token = localStorage.getItem('userToken')
  if (token) {
    try {
      const decoded = jwtDecode(token)
      userName.value = decoded.firstName || 'Kullanıcı'
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