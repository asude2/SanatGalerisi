<template>
  <div class="flex flex-wrap gap-4 mt-4">
    <div class="bg-gray-50 px-4 py-2 rounded-xl flex items-center gap-2 border border-gray-100">
      <span class="text-xl">👁️</span>
      <span class="font-bold text-gray-700">{{ stats.views || 0 }}</span>
      <span class="text-gray-500 text-sm">Görüntülenme</span>
    </div>
    
    <div class="bg-gray-50 px-4 py-2 rounded-xl flex items-center gap-2 border border-gray-100">
      <span class="text-xl">⭐</span>
      <span class="font-bold text-gray-700">{{ stats.avgRating ? stats.avgRating.toFixed(1) : '0.0' }}</span>
      <span class="text-gray-500 text-sm">Puan</span>
    </div>

    <div class="bg-gray-50 px-4 py-2 rounded-xl flex items-center gap-2 border border-gray-100">
      <span class="text-xl">💬</span>
      <span class="font-bold text-gray-700">{{ stats.comments || 0 }}</span>
      <span class="text-gray-500 text-sm">Yorum</span>
    </div>

    <div v-if="targetType === 'Workshop'" class="bg-blue-50 px-4 py-2 rounded-xl flex items-center gap-2 border border-blue-100">
      <span class="text-xl">📊</span>
      <span class="font-bold text-blue-700">%{{ stats.occupancy ? stats.occupancy.toFixed(0) : '0' }}</span>
      <span class="text-blue-600 text-sm">Doluluk</span>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { jwtDecode } from 'jwt-decode'

const props = defineProps({
  targetId: { type: Number, required: true },
  targetType: { type: String, required: true }
})

const stats = ref({})

const fetchStats = async () => {
  try {
    const res = await axios.get(`http://localhost:8080/entity-stats?targetId=${props.targetId}&targetType=${props.targetType}`)
    stats.value = res.data
  } catch (error) {
    console.error("İstatistikler yüklenemedi", error)
  }
}

const logView = async () => {
  try {
    let userId = null
    const token = localStorage.getItem('userToken')
    if (token) {
      userId = jwtDecode(token).userId
    }
    
    await axios.post('http://localhost:8080/log-interaction', {
      userId,
      targetId: props.targetId,
      targetType: props.targetType,
      interactionType: 'View'
    })
  } catch (error) {
    console.error("Görüntülenme kaydedilemedi", error)
  }
}

onMounted(() => {
  fetchStats()
  logView()
})
</script>