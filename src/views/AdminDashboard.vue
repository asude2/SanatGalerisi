<template>
  <div class="min-h-screen bg-gray-50 py-12 px-6">
    <div class="max-w-7xl mx-auto space-y-8">
      
      <div class="flex justify-between items-end">
        <div>
          <h1 class="text-4xl font-black text-gray-900 tracking-tight">Yönetici Paneli 🏛️</h1>
          <p class="text-gray-500 font-medium mt-2">Galerinizin genel durumunu ve istatistiklerini buradan takip edin.</p>
        </div>
        <div class="text-right">
          <p class="text-xs font-black text-gray-400 uppercase tracking-widest">Son Güncelleme</p>
          <p class="text-sm font-bold text-gray-600">{{ currentTime }}</p>
        </div>
      </div>

      <!-- Özet Kartları -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div class="bg-white p-8 rounded-3xl shadow-sm border border-gray-100 group hover:shadow-xl hover:shadow-blue-500/5 transition-all">
          <div class="w-12 h-12 bg-blue-100 rounded-2xl flex items-center justify-center text-blue-600 text-2xl mb-4 group-hover:scale-110 transition-transform">🖼️</div>
          <p class="text-sm font-bold text-gray-400 uppercase tracking-widest">Toplam Eser</p>
          <p class="text-4xl font-black text-gray-900 mt-1">{{ stats.totalArtworks }}</p>
        </div>
        
        <div class="bg-white p-8 rounded-3xl shadow-sm border border-gray-100 group hover:shadow-xl hover:shadow-indigo-500/5 transition-all">
          <div class="w-12 h-12 bg-indigo-100 rounded-2xl flex items-center justify-center text-indigo-600 text-2xl mb-4 group-hover:scale-110 transition-transform">🎨</div>
          <p class="text-sm font-bold text-gray-400 uppercase tracking-widest">Toplam Atölye</p>
          <p class="text-4xl font-black text-gray-900 mt-1">{{ stats.totalWorkshops }}</p>
        </div>

        <div class="bg-white p-8 rounded-3xl shadow-sm border border-gray-100 group hover:shadow-xl hover:shadow-green-500/5 transition-all">
          <div class="w-12 h-12 bg-green-100 rounded-2xl flex items-center justify-center text-green-600 text-2xl mb-4 group-hover:scale-110 transition-transform">👥</div>
          <p class="text-sm font-bold text-gray-400 uppercase tracking-widest">Kayıtlı Kullanıcı</p>
          <p class="text-4xl font-black text-gray-900 mt-1">{{ stats.totalUsers }}</p>
        </div>

        <div class="bg-white p-8 rounded-3xl shadow-sm border border-gray-100 group hover:shadow-xl hover:shadow-pink-500/5 transition-all">
          <div class="w-12 h-12 bg-pink-100 rounded-2xl flex items-center justify-center text-pink-600 text-2xl mb-4 group-hover:scale-110 transition-transform">🎧</div>
          <p class="text-sm font-bold text-gray-400 uppercase tracking-widest">Bekleyen Destek</p>
          <p class="text-4xl font-black text-gray-900 mt-1">{{ stats.activeTickets }}</p>
        </div>
      </div>

      <!-- Detaylı Analizler -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- En Çok İlgi Gören Eserler -->
        <div class="bg-white rounded-3xl shadow-sm border border-gray-100 overflow-hidden">
          <div class="p-6 border-b border-gray-50 flex justify-between items-center">
            <h3 class="text-xl font-bold text-gray-900 flex items-center gap-2">🔥 Popüler Eserler</h3>
            <router-link to="/" class="text-xs font-bold text-blue-600 hover:underline">Tümünü Gör</router-link>
          </div>
          <div class="p-6">
            <div class="space-y-6">
              <div v-for="artwork in popularArtworks" :key="artwork.id" class="flex items-center gap-4 group">
                <div class="w-12 h-12 bg-gray-100 rounded-xl overflow-hidden">
                  <img :src="artwork.imageUrl || 'https://images.unsplash.com/photo-1579783902614-a3fb3927b6a5?w=100'" class="w-full h-full object-cover" />
                </div>
                <div class="flex-1">
                  <p class="font-bold text-gray-900 text-sm group-hover:text-blue-600 transition-colors">{{ artwork.title }}</p>
                  <p class="text-xs text-gray-400">{{ artwork.artist }} • {{ artwork.category }}</p>
                </div>
                <div class="text-right">
                  <p class="text-sm font-black text-gray-900">{{ artwork.views }}</p>
                  <p class="text-[10px] font-bold text-gray-400 uppercase">Görüntüleme</p>
                </div>
              </div>
              <div v-if="popularArtworks.length === 0" class="text-center py-10">
                <p class="text-gray-400 font-bold">Henüz görüntüleme verisi yok. ✨</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Son Destek Talepleri -->
        <div class="bg-white rounded-3xl shadow-sm border border-gray-100 overflow-hidden">
          <div class="p-6 border-b border-gray-50 flex justify-between items-center">
            <h3 class="text-xl font-bold text-gray-900 flex items-center gap-2">📧 Son Talepler</h3>
            <router-link to="/support" class="text-xs font-bold text-blue-600 hover:underline">Tümünü Yönet</router-link>
          </div>
          <div class="p-6">
            <div class="space-y-6">
              <div v-for="ticket in recentTickets" :key="ticket.ticketId" class="flex items-center gap-4">
                <div class="w-10 h-10 bg-gray-50 rounded-full flex items-center justify-center text-lg">🎫</div>
                <div class="flex-1">
                  <p class="font-bold text-gray-900 text-sm">{{ ticket.subject }}</p>
                  <p class="text-xs text-gray-400">ID: #{{ ticket.ticketId }} • {{ formatDate(ticket.createdAt) }}</p>
                </div>
                <span class="px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-[10px] font-black uppercase tracking-tighter">{{ translateStatus(ticket.status) }}</span>
              </div>
              <div v-if="recentTickets.length === 0" class="text-center py-10">
                <p class="text-gray-400 font-bold">Bekleyen talep yok. ✨</p>
              </div>
            </div>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';

const stats = ref({
  totalArtworks: 0,
  totalWorkshops: 0,
  totalUsers: 0,
  activeTickets: 0
});

const recentTickets = ref([]);
const popularArtworks = ref([]);
const currentTime = ref(new Date().toLocaleString('tr-TR'));

const fetchStats = async () => {
  const token = localStorage.getItem('userToken');
  const config = { headers: { Authorization: `Bearer ${token}` } };
  
  try {
    const response = await axios.get('http://localhost:8080/admin/dashboard-stats', config);
    stats.value = response.data;
    
    // Son biletleri çek
    const ticketsResponse = await axios.get('http://localhost:8080/admin/tickets', config);
    recentTickets.value = (ticketsResponse.data || []).slice(0, 5);

    // Popüler eserleri çek
    const popularResponse = await axios.get('http://localhost:8080/admin/popular-artworks', config);
    popularArtworks.value = popularResponse.data || [];
  } catch (error) {
    console.error("Dashboard verileri çekilemedi:", error);
  }
};

const translateStatus = (status) => {
  const map = { 'Open': 'Açık', 'Açık': 'Açık', 'Responded': 'Yanıtlandı', 'Beklemede': 'Yanıtlandı', 'Closed': 'Çözüldü', 'Çözüldü': 'Çözüldü' };
  return map[status] || status;
};

const formatDate = (dateString) => {
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('tr-TR', { day: 'numeric', month: 'short' }).format(date);
};

onMounted(() => {
  fetchStats();
  setInterval(() => {
    currentTime.value = new Date().toLocaleString('tr-TR');
  }, 1000);
});
</script>
