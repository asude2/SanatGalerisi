<template>
  <div class="flex flex-wrap gap-4 mt-4">
    <!-- İzlenme -->
    <div class="bg-gray-100 px-4 py-2 rounded-xl flex items-center gap-2 group hover:bg-white hover:shadow-sm transition-all">
      <span class="text-xl">👁️</span>
      <div>
        <p class="text-[10px] font-black text-gray-400 uppercase leading-none">İzlenme</p>
        <p class="text-sm font-bold text-gray-700 leading-tight">{{ stats.views }}</p>
      </div>
    </div>

    <!-- Beğeni -->
    <button 
      @click="logInteraction('Like')"
      class="bg-pink-50 px-4 py-2 rounded-xl flex items-center gap-2 hover:bg-pink-100 transition-all cursor-pointer active:scale-95 shadow-sm shadow-pink-100/20"
    >
      <span class="text-xl">❤️</span>
      <div class="text-left">
        <p class="text-[10px] font-black text-pink-400 uppercase leading-none">Beğeni</p>
        <p class="text-sm font-bold text-pink-600 leading-tight">{{ stats.likes }}</p>
      </div>
    </button>

    <!-- Yorum Sayısı -->
    <div class="bg-blue-50 px-4 py-2 rounded-xl flex items-center gap-2 shadow-sm shadow-blue-100/20">
      <span class="text-xl">💬</span>
      <div>
        <p class="text-[10px] font-black text-blue-400 uppercase leading-none">Yorum</p>
        <p class="text-sm font-bold text-blue-600 leading-tight">{{ stats.comments }}</p>
      </div>
    </div>

    <!-- Puan (Eğer varsa) -->
    <div v-if="stats.avgRating > 0" class="bg-yellow-50 px-4 py-2 rounded-xl flex items-center gap-2 shadow-sm shadow-yellow-100/20">
      <span class="text-xl">⭐</span>
      <div>
        <p class="text-[10px] font-black text-yellow-500 uppercase leading-none">Puan</p>
        <p class="text-sm font-bold text-yellow-600 leading-tight">{{ stats.avgRating }} / 5</p>
      </div>
    </div>

    <!-- Workshop'a Özel: Doluluk ve Rezervasyon -->
    <template v-if="targetType === 'Workshop'">
      <div class="bg-indigo-50 px-4 py-2 rounded-xl flex items-center gap-2 shadow-sm shadow-indigo-100/20">
        <span class="text-xl">🎟️</span>
        <div>
          <p class="text-[10px] font-black text-indigo-400 uppercase leading-none">Rezervasyon</p>
          <p class="text-sm font-bold text-indigo-600 leading-tight">{{ stats.reservations }}</p>
        </div>
      </div>
      
      <div class="bg-green-50 px-4 py-2 rounded-xl flex items-center gap-2 shadow-sm shadow-green-100/20">
        <span class="text-xl">📊</span>
        <div>
          <p class="text-[10px] font-black text-green-400 uppercase leading-none">Doluluk</p>
          <p class="text-sm font-bold text-green-600 leading-tight">%{{ stats.occupancy?.toFixed(0) }}</p>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';

const props = defineProps({
  targetId: { type: Number, required: true },
  targetType: { type: String, required: true }
});

const stats = ref({
  views: 0,
  likes: 0,
  comments: 0,
  avgRating: 0,
  reservations: 0,
  occupancy: 0
});

const fetchStats = async () => {
  try {
    const response = await axios.get(`http://localhost:8080/entity-stats?targetId=${props.targetId}&targetType=${props.targetType}`);
    stats.value = response.data;
  } catch (error) {
    console.error("İstatistikler yüklenemedi:", error);
  }
};

const logInteraction = async (type) => {
  try {
    const token = localStorage.getItem('userToken');
    const userId = localStorage.getItem('userId');
    
    if (type === 'Like' && !token) {
      alert("Beğenmek için lütfen önce giriş yapın! 👤");
      return;
    }

    const config = token ? { headers: { Authorization: `Bearer ${token}` } } : {};

    await axios.post('http://localhost:8080/log-interaction', {
      userId: userId ? parseInt(userId) : null,
      targetId: props.targetId,
      targetType: props.targetType,
      interactionType: type
    }, config);
    fetchStats();
  } catch (error) {
    console.error("Log hatası:", error);
    if (error.response && error.response.status === 401) {
      alert("Beğenmek için lütfen giriş yapın! 👤");
    }
  }
};

onMounted(() => {
  fetchStats();
  logInteraction('View'); // Sayfa her açıldığında görüntüleme logla
});
</script>
