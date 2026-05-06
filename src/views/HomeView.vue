<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header Kısmı (Aynı Kalıyor) -->
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
            <span class="cursor-pointer opacity-70 hover:opacity-100 transition-opacity">Sanatçılar</span>
          </nav>
          <div class="bg-white text-galeri-yesil p-3 rounded-2xl shadow-inner cursor-pointer">
            <router-link to="/profile" class="cursor-pointer">
              <span class="text-2xl">👤</span>
            </router-link>
          </div>
        </div>
      </div>
    </header>

    <!-- Çıkış Butonu -->
    <div class="mt-10 px-10">
        <button 
          @click="handleLogout" 
          class="flex justify-center gap-2 px-5 bg-red-50 text-red-600 font-bold py-3 rounded-2xl hover:bg-red-100 transition-all border border-red-100"
        >
          <span>🚪</span> Çıkış Yap
        </button>
    </div>

    <main class="container mx-auto -mt-8 px-4 pb-20">
      <!-- Üst Bar: Arama ve Filtreleme -->
      <div class="bg-white p-4 rounded-2xl shadow-sm mb-12 flex flex-col md:flex-row gap-4 items-center justify-between">
        <div class="relative w-full md:w-96">
          <input 
            v-model="searchQuery" 
            type="text" 
            placeholder="Eser veya sanatçı ara..." 
            class="w-full p-3 pl-12 bg-gray-50 rounded-xl focus:outline-none focus:ring-2 focus:ring-galeri-yesil/20" 
          />
          <span class="absolute left-4 top-3">🔍</span>
        </div>
        
        <div class="flex gap-4">
          <button class="px-6 py-2 bg-gray-100 rounded-xl font-semibold hover:bg-gray-200">Kategoriler</button>
          <button class="px-6 py-2 bg-gray-100 rounded-xl font-semibold hover:bg-gray-200">Sırala</button>
        </div>
      </div>

      <!-- Eser Listesi (Grid Arama Barının Dışında Olmalı) -->
      <div v-if="filteredArtworks.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
        <ArtworkCard 
          v-for="artwork in filteredArtworks" 
          :key="artwork.id"
          :id="artwork.id"
          :title="artwork.title" 
          :artist="artwork.artist" 
          :price="artwork.price" 
          :image="artwork.image"
        />
      </div>

      <!-- Veri yoksa veya Arama sonucu boşsa gösterilecek mesaj -->
      <div v-else class="col-span-full text-center py-20">
        <span class="text-6xl block mb-4">🏜️</span>
        <p class="text-gray-400 text-xl font-medium">
          {{ searchQuery ? '"' + searchQuery + '" ile eşleşen eser bulunamadı.' : 'Henüz sergilenecek eser bulunamadı...' }}
        </p>
      </div>
    </main>
  </div>
</template>



<script setup>
import { ref, onMounted, computed } from 'vue' // computed eklendi
import { useRouter } from 'vue-router'
import { jwtDecode } from 'jwt-decode'
import axios from 'axios'
import ArtworkCard from '../components/ArtworkCard.vue'

const router = useRouter()
const userName = ref('Misafir')
const artworks = ref([])
const searchQuery = ref('') // Arama terimi için

// API'den eserleri çek
const fetchArtworks = async () => {
  try {
    const response = await axios.get('http://localhost:8080/artworks')
    // Eğer gelen veri null ise boş bir dizi ata, değilse gelen veriyi ata
    artworks.value = response.data || [] 
  } catch (error) {
    console.error("Eserler yüklenirken bir hata oluştu:", error)
    artworks.value = [] // Hata durumunda da boş dizi olsun ki sayfa çökmesin
  }
}

// ÖDEV MADDE 2: Arama ve Filtreleme Mantığı
const filteredArtworks = computed(() => {
  // Veri henüz gelmemişse veya null ise boş dizi döndür
  if (!artworks.value) return [] 

  const search = searchQuery.value.toLowerCase().trim()
  return artworks.value.filter(artwork => {
    return artwork.title.toLowerCase().includes(search) || 
           artwork.artist.toLowerCase().includes(search)
  })
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
  fetchArtworks()
})

const handleLogout = () => {
  localStorage.removeItem('userToken')
  localStorage.removeItem('userEmail') // Bunu da temizleyelim
  router.push('/login')
}

// NOT: goToDetail fonksiyonunu buradan sildik çünkü ArtworkCard içinde tanımlı.

</script>