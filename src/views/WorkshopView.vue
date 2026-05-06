<template>
  <div class="min-h-screen bg-gray-50 pb-20">
    <header class="bg-white border-b p-10 mb-10">
      <div class="container mx-auto">
        <button @click="router.push('/')" class="text-blue-600 font-bold mb-4 flex items-center gap-2">⬅️ Galeriye Dön</button>
        <h1 class="text-5xl font-black text-gray-900">Atölyeler & Etkinlikler</h1>
        <p class="text-gray-500 mt-2 text-lg">Sanatla iç içe eğitimler ve söyleşiler seni bekliyor.</p>
      </div>
    </header>

    <main class="container mx-auto px-6">
      <div v-if="workshops.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <WorkshopCard 
          v-for="ws in workshops" 
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
        Şu an aktif bir atölye bulunmamaktadır...
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import WorkshopCard from '../components/WorkshopCard.vue'

const router = useRouter()
const workshops = ref([])

onMounted(async () => {
  try {
    const response = await axios.get('http://localhost:8080/workshops')
    workshops.value = response.data || []
  } catch (error) {
    console.error("Atölyeler yüklenemedi:", error)
  }
})
</script>