<template>
    <div class="p-6">
        <div class="relative h-52">
            <img :src="image" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" />
            <div class="absolute top-4 left-4 bg-white/90 backdrop-blur px-3 py-1 rounded-lg shadow text-sm font-bold text-blue-600">
                {{ price > 0 ? price + ' ₺' : 'Ücretsiz' }}
            </div>
        </div>
        <h3 class="text-xl font-bold text-gray-800 mb-1">{{ title }}</h3>
        <div class="grid grid-cols-2 gap-4 text-sm text-gray-500 mb-6">
        <div class="flex items-center gap-2"><span>📅</span> {{ date }}</div>
        </div>
       <button @click="goToDetail" class="w-full py-3 bg-gray-900 text-white rounded-xl font-bold hover:bg-blue-600 transition-colors">
        Detayları Gör & Kayıt Ol
      </button>
    </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
const props = defineProps(['id', 'title', 'instructorName', 'date', 'location', 'capacity', 'price', 'image'])
const router = useRouter()

const goToDetail = () => {
  router.push(`/workshops/${props.id}`)
}
</script>