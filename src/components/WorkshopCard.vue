<template>
    <div class="p-6 relative group">
        <div class="relative h-52 mb-4 overflow-hidden rounded-2xl">
            <img :src="image" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" />
            <div class="absolute top-4 left-4 bg-white/90 backdrop-blur px-3 py-1 rounded-lg shadow text-sm font-bold text-blue-600">
                {{ price > 0 ? price + ' ₺' : 'Ücretsiz' }}
            </div>
            <!-- Karşılaştırma Butonu -->
            <button 
                @click.stop="addToCompare" 
                class="absolute top-4 right-4 p-2 bg-white/90 backdrop-blur rounded-full shadow hover:bg-blue-600 hover:text-white transition-all text-gray-400 cursor-pointer"
                title="Karşılaştır"
            >
                ⚖️
            </button>
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

const addToCompare = () => {
  const currentList = JSON.parse(localStorage.getItem('compareList') || '[]');
  const currentType = localStorage.getItem('compareType') || 'Workshop';

  if (currentType !== 'Workshop') {
    if (!confirm('Karşılaştırma listenizdeki eserler temizlenecek. Devam edilsin mi?')) return;
    localStorage.setItem('compareList', JSON.stringify([props.id]));
    localStorage.setItem('compareType', 'Workshop');
  } else {
    if (!currentList.includes(props.id)) {
      if (currentList.length >= 4) return alert('En fazla 4 ürünü karşılaştırabilirsiniz!');
      currentList.push(props.id);
      localStorage.setItem('compareList', JSON.stringify(currentList));
      localStorage.setItem('compareType', 'Workshop');
    }
  }
  router.push('/compare');
};
</script>