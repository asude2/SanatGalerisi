<template>
  <div class="min-h-screen bg-gray-50 py-12 px-6">
    <div v-if="workshop" class="max-w-4xl mx-auto bg-white rounded-3xl shadow-2xl overflow-hidden">
      <div class="relative h-96">
        <img :src="workshop.image" class="w-full h-full object-cover" />
        <button @click="router.back()" class="absolute top-6 left-6 bg-white/90 p-3 rounded-full shadow-lg">
          ⬅️ Geri Dön
        </button>
      </div>

      <div class="p-10">
        <div class="flex justify-between items-start mb-6">
          <div>
            <h1 class="text-4xl font-black text-gray-900">{{ workshop.title }}</h1>
            <p class="text-blue-600 font-bold text-lg mt-2 italic">Eğitmen: {{ workshop.instructorName }}</p>
          </div>
          <div class="bg-blue-50 text-blue-700 px-6 py-3 rounded-2xl font-black text-2xl">
            {{ workshop.price > 0 ? workshop.price + ' ₺' : 'Ücretsiz' }}
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-10 border-y py-8 border-gray-100">
          <div class="flex items-center gap-3">
            <span class="text-3xl">📅</span>
            <div>
              <p class="text-gray-400 text-xs font-bold uppercase">Tarih & Saat</p>
              <p class="font-bold text-gray-800">{{ workshop.date }}</p>
            </div>
          </div>
          <div class="flex items-center gap-3">
            <span class="text-3xl">📍</span>
            <div>
              <p class="text-gray-400 text-xs font-bold uppercase">Konum</p>
              <p class="font-bold text-gray-800">{{ workshop.location }}</p>
            </div>
          </div>
          <div class="flex items-center gap-3">
            <span class="text-3xl">👥</span>
            <div>
              <p class="text-gray-400 text-xs font-bold uppercase">Kontenjan</p>
              <p class="font-bold text-gray-800">{{ workshop.capacity }} Kişi</p>
            </div>
          </div>
        </div>

        <div class="mb-10">
          <h3 class="text-xl font-bold text-gray-800 mb-4">Etkinlik Hakkında</h3>
          <p class="text-gray-600 leading-relaxed text-lg whitespace-pre-line">
            {{ workshop.description }}
          </p>
        </div>

        <button class="w-full py-5 bg-blue-600 text-white rounded-2xl font-black text-xl hover:bg-blue-700 transition-all shadow-xl shadow-blue-100">
          Atölyeye Kayıt Ol ✨
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const router = useRouter()
const workshop = ref(null)

onMounted(async () => {
  try {
    const response = await axios.get('http://localhost:8080/workshops')
    const all = response.data
    workshop.value = all.find(w => w.id === parseInt(route.params.id))
  } catch (error) {
    console.error("Detay yüklenemedi:", error)
  }
})
</script>