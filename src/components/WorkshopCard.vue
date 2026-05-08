<template>
    <div class="p-6">
        <div class="relative h-52 mb-4 overflow-hidden rounded-2xl">
            <img :src="image" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" />
            <div class="absolute top-4 left-4 bg-white/90 backdrop-blur px-3 py-1 rounded-lg shadow text-sm font-bold text-blue-600">
                {{ price > 0 ? price + ' ₺' : 'Ücretsiz' }}
            </div>
        </div>

        <h3 class="text-xl font-bold text-gray-800 mb-2 pb-5">{{ title }}</h3>
        

        <button @click="goToDetail" class="w-full py-3 bg-gray-900 text-white rounded-xl font-bold hover:bg-blue-600 transition-colors cursor-pointer">
            Detayları Gör & Kayıt Ol
        </button>
    </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

// Props isimleri Backend'den küçük harf (availableDates) olarak geliyorsa bu şekilde kalmalı
const props = defineProps(['id', 'title', 'instructorName', 'availableDates', 'location', 'capacity', 'price', 'image'])
const router = useRouter()

const goToDetail = () => {
    router.push(`/workshops/${props.id}`)
}
</script>