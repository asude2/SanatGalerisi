<template>
  <div class="min-h-screen bg-gray-50 py-12 px-6 pb-24">
    <div class="max-w-7xl mx-auto">
      
      <div class="flex flex-col md:flex-row justify-between items-start md:items-center mb-12 gap-6">
        <div>
          <h1 class="text-4xl font-black text-gray-900 tracking-tight">Karşılaştırma Merkezi ⚖️</h1>
          <p class="text-gray-500 font-medium mt-2">
            Seçtiğiniz {{ targetType === 'Artwork' ? 'eserleri' : 'atölyeleri' }} yan yana analiz edin. 
            <span class="text-blue-600 cursor-pointer hover:underline" @click="toggleType">
              ({{ targetType === 'Artwork' ? 'Atölyelere Geç' : 'Eserlere Geç' }})
            </span>
          </p>
        </div>
        <div class="flex gap-4">
          <button 
            @click="prepareSave"
            class="px-6 py-3 bg-blue-600 text-white rounded-xl font-bold shadow-lg shadow-blue-100 hover:bg-blue-700 transition-all"
          >
            Analizi Kaydet 💾
          </button>
          <button 
            @click="clearComparison"
            class="px-6 py-3 bg-gray-200 text-gray-600 rounded-xl font-bold hover:bg-gray-300 transition-all"
          >
            Temizle
          </button>
        </div>
      </div>

      <div v-if="items.length > 0" class="space-y-12">
        <div class="overflow-x-auto pb-8">
          <div class="flex gap-6 min-w-max">
            <!-- Karşılaştırma Kartları -->
            <div v-for="item in items" :key="item.id" class="w-80 bg-white rounded-3xl shadow-sm border border-gray-100 overflow-hidden flex flex-col hover:shadow-md transition-all">
              <div class="h-48 relative">
                <img :src="item.imageUrl || item.image" class="w-full h-full object-cover" />
                <button @click="removeItem(item.id)" class="absolute top-4 right-4 bg-white/90 p-2 rounded-full shadow-md text-red-500 hover:bg-red-50 transition-all">✕</button>
              </div>
              
              <div class="p-6 space-y-6 flex-1">
                <div>
                  <h3 class="text-xl font-black text-gray-900 line-clamp-1">{{ item.title }}</h3>
                  <p class="text-blue-600 font-bold text-sm">{{ item.artist || item.instructorName }}</p>
                </div>

                <div class="space-y-4">
                  <!-- Eser Özellikleri -->
                  <template v-if="targetType === 'Artwork'">
                    <div class="p-4 bg-gray-50 rounded-2xl">
                      <p class="text-[10px] font-black text-gray-400 uppercase mb-1">Kategori</p>
                      <p class="font-bold text-gray-700">{{ item.category || 'Belirtilmemiş' }}</p>
                    </div>
                    <div class="p-4 bg-blue-50 rounded-2xl border border-blue-100">
                      <p class="text-[10px] font-black text-blue-400 uppercase mb-1">Fiyat</p>
                      <p class="text-2xl font-black text-blue-600">{{ item.price }} ₺</p>
                    </div>
                  </template>

                  <!-- Workshop Özellikleri -->
                  <template v-else>
                    <div class="p-4 bg-gray-50 rounded-2xl">
                      <p class="text-[10px] font-black text-gray-400 uppercase mb-1">Tarihler</p>
                      <p class="text-sm font-bold text-gray-700">{{ item.availableDates }}</p>
                    </div>
                    <div class="p-4 bg-blue-50 rounded-2xl border border-blue-100">
                      <p class="text-[10px] font-black text-blue-400 uppercase mb-1">Ücret</p>
                      <p class="text-xl font-black text-blue-600">{{ item.price }} ₺</p>
                    </div>
                  </template>

                  <!-- Ortak İstatistikler -->
                  <div v-if="item.stats" class="grid grid-cols-2 gap-2 border-t border-gray-50 pt-4">
                    <div class="p-3 bg-gray-50 rounded-xl border border-gray-100 text-center">
                      <p class="text-[9px] font-black text-gray-400 uppercase mb-0.5">👁️ İzlenme</p>
                      <p class="text-sm font-black text-gray-700">{{ item.stats.views }}</p>
                    </div>
                    <div class="p-3 bg-pink-50 rounded-xl border border-pink-100 text-center">
                      <p class="text-[9px] font-black text-pink-400 uppercase mb-0.5">❤️ Beğeni</p>
                      <p class="text-sm font-black text-pink-600">{{ item.stats.likes }}</p>
                    </div>
                    <div class="p-3 bg-blue-50 rounded-xl border border-blue-100 text-center">
                      <p class="text-[9px] font-black text-blue-400 uppercase mb-0.5">💬 Yorum</p>
                      <p class="text-sm font-black text-blue-600">{{ item.stats.comments }}</p>
                    </div>
                    <div class="p-3 bg-yellow-50 rounded-xl border border-yellow-100 text-center">
                      <p class="text-[9px] font-black text-yellow-500 uppercase mb-0.5">⭐ Puan</p>
                      <p class="text-sm font-black text-yellow-600">{{ item.stats.avgRating || '0' }}</p>
                    </div>
                  </div>
                </div>

                <router-link 
                  :to="targetType === 'Artwork' ? `/artwork/${item.id}` : `/workshops/${item.id}`"
                  class="block w-full py-4 text-center bg-gray-900 text-white rounded-2xl font-bold hover:bg-gray-800 transition-all mt-auto"
                >
                  Detaya Git
                </router-link>
              </div>
            </div>

            <!-- Yeni Ekle Slotu -->
            <div 
              v-if="items.length < 4" 
              @click="showSelector = true"
              class="w-80 bg-gray-100/50 border-4 border-dashed border-gray-200 rounded-3xl flex flex-col items-center justify-center p-12 text-center group cursor-pointer hover:bg-gray-100 transition-all"
            >
              <div class="w-16 h-16 bg-white rounded-full flex items-center justify-center text-3xl shadow-sm mb-4 group-hover:scale-110 transition-transform">➕</div>
              <p class="text-gray-400 font-bold">Yeni {{ targetType === 'Artwork' ? 'Eser' : 'Atölye' }} Ekle</p>
              <p class="text-xs text-gray-300 mt-2 italic">Listeden seçim yapmak için tıklayın.</p>
            </div>
          </div>
        </div>

        <!-- Analiz Özeti -->
        <div class="bg-white rounded-[40px] p-12 border border-gray-100 shadow-sm">
          <h2 class="text-2xl font-black text-gray-900 mb-8 flex items-center gap-3">📊 Karşılaştırmalı Analiz Raporu</h2>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
            <div class="space-y-3" v-if="analysis.cheapest">
              <p class="text-xs font-black text-gray-400 uppercase tracking-widest">💰 En Ekonomik</p>
              <div class="bg-green-50 p-4 rounded-2xl border border-green-100">
                <p class="font-bold text-green-700">{{ analysis.cheapest.title }}</p>
                <p class="text-2xl font-black text-green-600 mt-1">{{ analysis.cheapest.price }} ₺</p>
              </div>
            </div>
            <div class="space-y-3" v-if="analysis.mostPopular">
              <p class="text-xs font-black text-gray-400 uppercase tracking-widest">🔥 En Popüler</p>
              <div class="bg-blue-50 p-4 rounded-2xl border border-blue-100">
                <p class="font-bold text-blue-700">{{ analysis.mostPopular.title }}</p>
                <p class="text-2xl font-black text-blue-600 mt-1">{{ analysis.mostPopular.stats.views }} İzlenme</p>
              </div>
            </div>
            <div class="space-y-3" v-if="analysis.mostLiked">
              <p class="text-xs font-black text-gray-400 uppercase tracking-widest">❤️ En Çok Beğeni</p>
              <div class="bg-pink-50 p-4 rounded-2xl border border-pink-100">
                <p class="font-bold text-pink-700">{{ analysis.mostLiked.title }}</p>
                <p class="text-2xl font-black text-pink-600 mt-1">{{ analysis.mostLiked.stats.likes }} Beğeni</p>
              </div>
            </div>
            <div class="space-y-3" v-if="analysis.highestRated">
              <p class="text-xs font-black text-gray-400 uppercase tracking-widest">⭐ En Yüksek Puan</p>
              <div class="bg-yellow-50 p-4 rounded-2xl border border-yellow-100">
                <p class="font-bold text-yellow-700">{{ analysis.highestRated.title }}</p>
                <p class="text-2xl font-black text-yellow-600 mt-1">{{ analysis.highestRated.stats.avgRating }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="text-center py-32 bg-white rounded-[40px] shadow-sm border border-gray-100">
        <p class="text-8xl mb-8">🔭</p>
        <h3 class="text-3xl font-black text-gray-900">Karşılaştırılacak bir şey yok!</h3>
        <p class="text-gray-400 mt-4 max-w-md mx-auto text-lg font-medium leading-relaxed">
          Yukarıdaki "+" butonuna tıklayarak veya gözatarak karşılaştırma listene ürün ekleyebilirsin.
        </p>
        <div class="flex justify-center gap-4 mt-12">
          <button @click="showSelector = true" class="px-8 py-4 bg-blue-600 text-white rounded-2xl font-bold shadow-lg shadow-blue-100">Hemen Ürün Ekle</button>
        </div>
      </div>
    </div>

    <!-- Seçim Modalı -->
    <div v-if="showSelector" class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-[40px] p-10 max-w-4xl w-full shadow-2xl max-h-[85vh] flex flex-col">
        <div class="flex items-center justify-between mb-8">
          <div>
            <h3 class="text-3xl font-black text-gray-900">{{ targetType === 'Artwork' ? 'Eser Seç' : 'Atölye Seç' }} 🎨</h3>
            <p class="text-gray-400 font-medium">Karşılaştırmak istediğiniz öğeyi seçin.</p>
          </div>
          <button @click="showSelector = false" class="text-3xl text-gray-300 hover:text-gray-600 transition-colors">✕</button>
        </div>

        <div class="flex-1 overflow-y-auto pr-2 space-y-4">
          <div 
            v-for="available in availableItems" 
            :key="available.id"
            @click="addItem(available)"
            class="flex items-center gap-6 p-4 rounded-3xl border border-gray-100 hover:bg-blue-50 hover:border-blue-200 cursor-pointer transition-all group"
          >
            <img :src="available.imageUrl || available.image" class="w-20 h-20 rounded-2xl object-cover" />
            <div class="flex-1">
              <p class="font-black text-gray-900 group-hover:text-blue-700">{{ available.title }}</p>
              <p class="text-sm text-gray-500 font-bold">{{ available.artist || available.instructorName }}</p>
            </div>
            <div class="text-right">
              <p class="font-black text-blue-600 text-lg">{{ available.price }} ₺</p>
              <p class="text-[10px] font-black text-gray-300 uppercase">{{ available.category || 'Eser' }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Kaydetme Modalı -->
    <div v-if="showSaveModal" class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-[60] p-4">
      <div class="bg-white rounded-[40px] p-10 max-w-md w-full shadow-2xl">
        <h3 class="text-2xl font-black text-gray-900 mb-2">Analizi Kaydet 💾</h3>
        <p class="text-gray-400 mb-8 font-medium">Bu karşılaştırmaya bir isim vererek profilinizden ulaşabilirsiniz.</p>
        
        <input 
          v-model="comparisonTitle" 
          type="text" 
          placeholder="Örn: Salon için Tablolar"
          class="w-full p-4 bg-gray-50 rounded-2xl border-2 border-transparent focus:border-blue-500 outline-none mb-6 transition-all font-bold"
        />

        <div class="flex gap-4">
          <button @click="saveComparison" class="flex-1 py-4 bg-blue-600 text-white rounded-2xl font-bold hover:bg-blue-700 transition-all">Kaydet</button>
          <button @click="showSaveModal = false" class="flex-1 py-4 bg-gray-100 text-gray-600 rounded-2xl font-bold hover:bg-gray-200 transition-all">İptal</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import axios from 'axios';

const items = ref([]);
const targetType = ref('Artwork'); 
const allItems = ref([]);
const showSelector = ref(false);
const showSaveModal = ref(false);
const comparisonTitle = ref('');

const availableItems = computed(() => {
  const currentIds = items.value.map(i => i.id);
  return allItems.value.filter(a => !currentIds.includes(a.id));
});

const analysis = computed(() => {
  if (items.value.length === 0) return {};
  const sortedByPrice = [...items.value].sort((a, b) => a.price - b.price);
  const sortedByViews = [...items.value].sort((a, b) => (b.stats?.views || 0) - (a.stats?.views || 0));
  const sortedByLikes = [...items.value].sort((a, b) => (b.stats?.likes || 0) - (a.stats?.likes || 0));
  const sortedByRating = [...items.value].sort((a, b) => (b.stats?.avgRating || 0) - (a.stats?.avgRating || 0));
  return { cheapest: sortedByPrice[0], mostPopular: sortedByViews[0], mostLiked: sortedByLikes[0], highestRated: sortedByRating[0] };
});

const toggleType = () => {
  targetType.value = targetType.value === 'Artwork' ? 'Workshop' : 'Artwork';
  localStorage.setItem('compareType', targetType.value);
  items.value = [];
  fetchAll();
};

const fetchAll = async () => {
  try {
    const endpoint = targetType.value === 'Artwork' ? 'http://localhost:8080/artworks' : 'http://localhost:8080/workshops';
    const response = await axios.get(endpoint);
    allItems.value = response.data || [];
    const stored = localStorage.getItem('compareList');
    if (stored) {
      const ids = JSON.parse(stored);
      const filtered = allItems.value.filter(a => ids.includes(a.id));
      items.value = await Promise.all(filtered.map(async (item) => {
        try {
          const stats = await axios.get(`http://localhost:8080/entity-stats?targetId=${item.id}&targetType=${targetType.value}`);
          return { ...item, stats: stats.data };
        } catch (e) {
          return { ...item, stats: { views: 0, likes: 0, comments: 0, avgRating: 0 } };
        }
      }));
    }
  } catch (error) { console.error(error); }
};

const addItem = async (item) => {
  try {
    const stats = await axios.get(`http://localhost:8080/entity-stats?targetId=${item.id}&targetType=${targetType.value}`);
    items.value.push({ ...item, stats: stats.data });
  } catch (e) {
    items.value.push({ ...item, stats: { views: 0, likes: 0, comments: 0, avgRating: 0 } });
  }
  showSelector.value = false;
  localStorage.setItem('compareList', JSON.stringify(items.value.map(i => i.id)));
};

const removeItem = (id) => {
  items.value = items.value.filter(i => i.id !== id);
  localStorage.setItem('compareList', JSON.stringify(items.value.map(i => i.id)));
};

const clearComparison = () => {
  items.value = [];
  localStorage.removeItem('compareList');
};

const prepareSave = () => {
  const token = localStorage.getItem('userToken');
  if (!token) return alert("Lütfen giriş yapın!");
  if (items.value.length === 0) return alert("Karşılaştırılacak ürün yok!");
  showSaveModal.value = true;
};

const saveComparison = async () => {
  const token = localStorage.getItem('userToken') || localStorage.getItem('token');
  if (!comparisonTitle.value) return alert("Lütfen bir isim verin!");

  try {
    const ids = items.value.map(i => i.id).join(',');
    
    const payload = {
      title: comparisonTitle.value,
      targetType: targetType.value,
      targetIds: ids
    };

    const response = await axios.post('http://localhost:8080/comparisons/save', 
      payload, 
      { headers: { Authorization: `Bearer ${token}` } }
    );
    
    alert("Karşılaştırma başarıyla kaydedildi! Profilinden ulaşabilirsin kanka. 💾⚖️");
    showSaveModal.value = false;
    comparisonTitle.value = '';
    
    // Kayıt biter bitmez profile yönlendir
    router.push('/profile');

  } catch (e) { 
    console.error("Analiz veritabanına yazılırken backend hata döndü:", e);
    // 🔥 KESİN ÇÖZÜM: Backend'den gelen asıl SQL hatasını Alert ile ekrana basıyoruz!
    const errMsg = e.response?.data?.message || "Bilinmeyen Sunucu Hatası";
  }
};

onMounted(() => {
  targetType.value = localStorage.getItem('compareType') || 'Artwork';
  fetchAll();
});
</script>
