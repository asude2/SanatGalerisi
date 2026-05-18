<template>
  <div class="min-h-screen bg-gray-50 py-12 px-6">
    <div class="max-w-6xl mx-auto">
      
      <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-6 mb-12">
        <div>
          <h1 class="text-4xl font-black text-gray-900 tracking-tight">Müşteri Destek Merkezi 🎧</h1>
          <p class="text-gray-500 font-medium mt-2">Sorularınız için buradayız. Taleplerinizi iletebilir veya mevcutları takip edebilirsiniz.</p>
        </div>
        <button 
          v-if="userRole !== 'Admin'"
          @click="showCreateModal = true"
          class="px-8 py-4 bg-blue-600 text-white rounded-2xl font-bold shadow-xl shadow-blue-100 hover:bg-blue-700 transition-all flex items-center gap-2"
        >
          <span>➕</span> Yeni Destek Talebi
        </button>
      </div>

      <!-- İstatistikler / Özet (Opsiyonel) -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
        <div class="bg-white p-6 rounded-3xl border border-gray-100 shadow-sm">
          <p class="text-sm font-bold text-gray-400 uppercase tracking-widest mb-1">Aktif Talepler</p>
          <p class="text-4xl font-black text-blue-600">{{ activeTicketsCount }}</p>
        </div>
        <div class="bg-white p-6 rounded-3xl border border-gray-100 shadow-sm">
          <p class="text-sm font-bold text-gray-400 uppercase tracking-widest mb-1">Çözülenler</p>
          <p class="text-4xl font-black text-green-600">{{ resolvedTicketsCount }}</p>
        </div>
        <div class="bg-white p-6 rounded-3xl border border-gray-100 shadow-sm">
          <p class="text-sm font-bold text-gray-400 uppercase tracking-widest mb-1">Ort. Yanıt Süresi</p>
          <p class="text-4xl font-black text-indigo-600">~2 Saat</p>
        </div>
      </div>

      <!-- Bilet Listesi -->
      <div class="bg-white rounded-3xl shadow-sm border border-gray-100 overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full text-left">
            <thead class="bg-gray-50/50 border-b border-gray-100">
              <tr>
                <th class="px-8 py-5 text-xs font-black text-gray-400 uppercase tracking-tighter">Konu</th>
                <th class="px-8 py-5 text-xs font-black text-gray-400 uppercase tracking-tighter">Tür</th>
                <th class="px-8 py-5 text-xs font-black text-gray-400 uppercase tracking-tighter">Durum</th>
                <th class="px-8 py-5 text-xs font-black text-gray-400 uppercase tracking-tighter">Son Güncelleme</th>
                <th class="px-8 py-5 text-xs font-black text-gray-400 uppercase tracking-tighter">İşlem</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-50">
              <tr v-for="ticket in tickets" :key="ticket.ticketId" class="hover:bg-gray-50 transition-all group">
                <td class="px-8 py-6">
                  <p class="font-bold text-gray-900 group-hover:text-blue-600 transition-colors">{{ ticket.subject }}</p>
                  <p class="text-xs text-gray-400 mt-1 line-clamp-1">{{ ticket.message }}</p>
                </td>
                <td class="px-8 py-6">
                  <span class="px-3 py-1 bg-gray-100 rounded-full text-xs font-bold text-gray-600">{{ ticket.supportType }}</span>
                </td>
                <td class="px-8 py-6">
                  <span 
                    class="px-3 py-1 rounded-full text-[10px] font-black uppercase tracking-tighter"
                    :class="{
                      'bg-blue-100 text-blue-700': ticket.status === 'Açık' || ticket.status === 'Open',
                      'bg-yellow-100 text-yellow-700': ticket.status === 'Beklemede' || ticket.status === 'Responded',
                      'bg-green-100 text-green-700': ticket.status === 'Çözüldü' || ticket.status === 'Closed'
                    }"
                  >
                    {{ translateStatus(ticket.status) }}
                  </span>
                </td>
                <td class="px-8 py-6 text-sm text-gray-500 font-medium">
                  {{ formatDate(ticket.updatedAt) }}
                </td>
                <td class="px-8 py-6">
                  <router-link 
                    :to="`/support/${ticket.ticketId}`"
                    class="text-blue-600 font-bold hover:underline underline-offset-4"
                  >
                    Detaylar →
                  </router-link>
                </td>
              </tr>
              <tr v-if="tickets.length === 0">
                <td colspan="5" class="px-8 py-20 text-center">
                  <p class="text-6xl mb-4">📭</p>
                  <p class="text-xl font-bold text-gray-400">Henüz hiç destek talebiniz bulunmuyor.</p>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Yeni Talep Oluşturma Modalı -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-3xl p-8 max-w-xl w-full shadow-2xl">
        <div class="flex items-center justify-between mb-8">
          <h3 class="text-2xl font-black text-gray-900">Yeni Destek Talebi ✨</h3>
          <button @click="showCreateModal = false" class="text-2xl text-gray-400 hover:text-gray-600">✕</button>
        </div>

        <div class="space-y-6">
          <div class="space-y-2">
            <label class="text-sm font-bold text-gray-500 ml-1 uppercase">Konu</label>
            <input 
              v-model="newTicket.subject"
              type="text" 
              placeholder="Örn: Satın alma hatası"
              class="w-full p-4 bg-gray-50 rounded-2xl border-2 border-transparent focus:border-blue-500 outline-none transition-all"
            />
          </div>

          <div class="space-y-2">
            <label class="text-sm font-bold text-gray-500 ml-1 uppercase">Destek Türü</label>
            <select v-model="newTicket.supportType" class="w-full p-4 bg-gray-50 rounded-2xl border-2 border-transparent focus:border-blue-500 outline-none transition-all">
              <option value="Teknik Destek">Teknik Destek</option>
              <option value="Satın Alma / Ödeme">Satın Alma / Ödeme</option>
              <option value="Eser / Sanatçı Sorgulama">Eser / Sanatçı Sorgulama</option>
              <option value="Diğer">Diğer</option>
            </select>
          </div>

          <div class="space-y-2">
            <label class="text-sm font-bold text-gray-500 ml-1 uppercase">Mesajınız</label>
            <textarea 
              v-model="newTicket.message"
              placeholder="Sorununuzu detaylıca açıklayın..."
              class="w-full p-4 bg-gray-50 rounded-2xl border-2 border-transparent focus:border-blue-500 outline-none transition-all min-h-[150px]"
            ></textarea>
          </div>

          <div class="flex gap-4 pt-4">
            <button 
              @click="submitTicket"
              :disabled="!newTicket.subject || !newTicket.message"
              class="flex-1 py-4 bg-blue-600 text-white rounded-2xl font-bold hover:bg-blue-700 disabled:opacity-50 transition-all"
            >
              Talebi Oluştur 🚀
            </button>
            <button 
              @click="showCreateModal = false"
              class="flex-1 py-4 bg-gray-100 text-gray-600 rounded-2xl font-bold hover:bg-200 transition-all"
            >
              İptal
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import axios from 'axios';
import { jwtDecode } from 'jwt-decode';

const tickets = ref([]);
const showCreateModal = ref(false);
const userRole = ref('');

const newTicket = ref({
  subject: '',
  supportType: 'Teknik Destek',
  message: ''
});

const activeTicketsCount = computed(() => tickets.value.filter(t => t.status !== 'Closed' && t.status !== 'Çözüldü').length);
const resolvedTicketsCount = computed(() => tickets.value.filter(t => t.status === 'Closed' || t.status === 'Çözüldü').length);

const fetchTickets = async () => {
  const token = localStorage.getItem('userToken');
  if (!token) return;

  try {
    const config = { headers: { Authorization: `Bearer ${token}` } };
    const endpoint = userRole.value === 'Admin' ? 'http://localhost:8080/admin/tickets' : 'http://localhost:8080/tickets';
    const response = await axios.get(endpoint, config);
    tickets.value = response.data || [];
  } catch (error) {
    console.error("Biletler yüklenemedi:", error);
  }
};

const submitTicket = async () => {
  const token = localStorage.getItem('userToken');
  try {
    const config = { headers: { Authorization: `Bearer ${token}` } };
    await axios.post('http://localhost:8080/tickets/create', newTicket.value, config);
    showCreateModal.value = false;
    newTicket.value = { subject: '', supportType: 'Teknik Destek', message: '' };
    fetchTickets();
    alert("Destek talebiniz başarıyla oluşturuldu! ✨");
  } catch (error) {
    alert("Talep oluşturulurken bir hata oluştu.");
  }
};

const translateStatus = (status) => {
  const map = {
    'Open': 'Açık',
    'Açık': 'Açık',
    'Responded': 'Yanıtlandı',
    'Beklemede': 'Yanıtlandı',
    'Closed': 'Çözüldü',
    'Çözüldü': 'Çözüldü'
  };
  return map[status] || status;
};

const formatDate = (dateString) => {
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('tr-TR', { 
    day: 'numeric', 
    month: 'short', 
    hour: '2-digit', 
    minute: '2-digit' 
  }).format(date);
};

onMounted(() => {
  const token = localStorage.getItem('userToken');
  if (token) {
    const decoded = jwtDecode(token);
    userRole.value = decoded.role;
    fetchTickets();
  }
});
</script>
