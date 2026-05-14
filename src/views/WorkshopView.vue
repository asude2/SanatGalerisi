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

      <div v-if="filteredWorkshops.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <WorkshopCard 
          v-for="ws in filteredWorkshops" 
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

const filteredWorkshops = computed(() => {
  const search = searchTerm.value.trim().toLowerCase()
  if (!search) return workshops.value

  return workshops.value.filter(ws => {
    const title = String(ws.title || '').toLowerCase()
    const instructor = String(ws.instructorName || '').toLowerCase()
    return title.includes(search) || instructor.includes(search)
  })
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