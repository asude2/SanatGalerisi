<template>
  <div class="bg-white rounded-3xl p-6 shadow-sm border border-gray-100 hover:shadow-xl transition-all group">
    <div class="flex flex-col items-center text-center">
      <div class="w-24 h-24 bg-galeri-yesil/10 rounded-full flex items-center justify-center text-4xl mb-4 group-hover:scale-110 transition-transform">
        👤
      </div>
      <h3 class="text-xl font-bold text-gray-900">{{ name }}</h3>
      <p class="text-galeri-yesil font-medium text-sm mb-3">{{ nationality || 'Sanatçı' }}</p>
      <p class="text-gray-500 text-sm line-clamp-2 mb-4">{{ biography || 'Henüz bir biyografi eklenmemiş.' }}</p>
      <div class="grid grid-cols-1 gap-3 w-full">
        <button 
          @click="$emit('view-artworks', name)"
          class="w-full py-3 bg-gray-50 text-gray-600 font-bold rounded-xl hover:bg-galeri-yesil hover:text-white transition-colors"
        >
          Eserlerini Gör
        </button>
        <button 
          @click="$emit('view-workshops', name)"
          class="w-full py-3 bg-galeri-yesil text-white font-bold rounded-xl hover:bg-green-700 transition-colors"
        >
          Atölyelerini Gör
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps(['name', 'biography', 'nationality']);
</script>