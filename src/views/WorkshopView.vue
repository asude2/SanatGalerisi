<template>
  <div class="min-h-screen bg-gray-50 pb-20">
    <header class="bg-white border-b p-10 mb-10">
      <div class="container mx-auto">
        <button @click="router.push('/')" class="text-blue-600 font-bold mb-4 flex items-center gap-2 cursor-pointer">⬅️ Galeriye Dön</button>
        <h1 class="text-5xl font-black text-gray-900">Atölyeler & Etkinlikler</h1>
        <p class="text-gray-500 mt-2 text-lg">Sanatla iç içe eğitimler ve söyleşiler seni bekliyor.</p>
      </div>
    </header>

    <main class="container mx-auto px-6">
      <div class="mb-8">
        <label class="block text-gray-700 font-semibold mb-2" for="workshopSearch">Atölye ara</label>
        <input
          id="workshopSearch"
          v-model="searchTerm"
          type="text"
          placeholder="Başlık ya da eğitmen adı girin"
          class="w-full rounded-xl border border-gray-300 px-4 py-3 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-200"
        />
      </div>

      <div v-if="processedWorkshops.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <WorkshopCard 
          v-for="ws in processedWorkshops"
          :key="ws.id" 
          :id="ws.id"
          :title="ws.title"
          :instructorName="ws.instructorName"
          :date="ws.date"
          :location="ws.location"
          :capacity="ws.capacity"
          :price="ws.price"
          :image="ws.image"
        />
      </div>
      <div v-else class="text-center py-20 text-gray-400">
        {{ workshops.length > 0 ? 'Aramanıza uygun atölye bulunamadı.' : 'Şu an aktif bir atölye bulunmamaktadır...' }}
      </div>

<div class="flex flex-wrap justify-end items-center gap-4 mb-8">
  <div class="flex items-center gap-2">
    <span class="text-sm font-bold text-gray-500 uppercase">Sırala:</span>
    <select 
      v-model="sortCriterion" 
      class="p-3 rounded-2xl border-2 border-gray-100 bg-white text-gray-700 font-bold outline-none focus:border-orange-500 transition-all shadow-sm"
    >
      <option value="default">Varsayılan</option>
      <option value="tarih">📅 Tarih</option>
      <option value="price">💰 Ücret</option>
      <option value="capacity">👥 Kontenjan</option>
    </select>
  </div>

      <div v-if="sortCriterion === 'price' || sortCriterion === 'capacity'" class="flex items-center gap-2 animate-fade-in">
        <span class="text-sm font-bold text-gray-500 uppercase">Yön:</span>
        <select 
          v-model="sortDirection" 
          class="p-3 rounded-2xl border-2 border-gray-100 bg-white text-gray-700 font-bold outline-none focus:border-orange-500 transition-all shadow-sm"
        >
          <option value="asc">Artan (Düşükten Yükseğe)</option>
          <option value="desc">Azalan (Yüksekten Düşüğe)</option>
        </select>
      </div>

      <div v-if="sortCriterion === 'tarih'" class="text-xs font-bold text-orange-600 bg-orange-50 px-3 py-2 rounded-xl border border-orange-100">
        ✨ En yakın etkinlikler en üstte gösteriliyor.
      </div>
    </div>

    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import WorkshopCard from '../components/WorkshopCard.vue'

const router = useRouter()
const workshops = ref([])
const searchTerm = ref('')

const sortCriterion = ref('default') 
const sortDirection = ref('asc')     

// --- ŞOV BURADA: Hem Filtreleme Hem Sıralama Tek Fonksiyonda ---
const processedWorkshops = computed(() => {
  if (!workshops.value) return []
  
  // 1. Önce Arama (Filter) işlemini yapalım
  const search = searchTerm.value.trim().toLowerCase()
  let result = workshops.value.filter(ws => {
    const title = String(ws.title || '').toLowerCase()
    const instructor = String(ws.instructorName || '').toLowerCase()
    return title.includes(search) || instructor.includes(search)
  })

  // 2. Sonra Arama sonuçlarını Sıralayalım (Sort)
  if (sortCriterion.value === 'price') {
    result.sort((a, b) => {
      return sortDirection.value === 'asc' ? a.price - b.price : b.price - a.price
    })
  } 
  else if (sortCriterion.value === 'capacity') {
    result.sort((a, b) => {
      return sortDirection.value === 'asc' ? a.capacity - b.capacity : b.capacity - a.capacity
    })
  } 
  else if (sortCriterion.value === 'tarih') {
    result.sort((a, b) => {
      // split(',')[0] yaparak ilk tarihi baz alıyoruz
      const dateA = new Date(a.availableDates?.split(',')[0] || '9999-12-31')
      const dateB = new Date(b.availableDates?.split(',')[0] || '9999-12-31')
      return dateA - dateB 
    })
  }

  else if (sortCriterion.value === 'tarih') {
    result.sort((a, b) => {
      // Tarihleri güvenli bir şekilde objeye çeviriyoruz
      // availableDates içinde birden fazla tarih varsa ilkini alıyoruz
      const dateA = new Date(a.availableDates?.split(',')[0]).getTime() || 0
      const dateB = new Date(b.availableDates?.split(',')[0]).getTime() || 0

      // sortDirection 'asc' ise En Yakın (Küçük olan tarih) başta
      // sortDirection 'desc' ise En Uzak (Büyük olan tarih) başta
      return sortDirection.value === 'asc' ? dateA - dateB : dateB - dateA
    })
  }

  return result
})

onMounted(async () => {
  try {
    const response = await axios.get('http://localhost:8080/workshops')
    workshops.value = response.data || []
  } catch (error) {
    console.error("Atölyeler yüklenemedi:", error)
  }
})
</script>