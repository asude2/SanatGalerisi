<template>
  <div class="min-h-screen bg-gray-50 py-12 px-6">
    <div v-if="ticket" class="max-w-4xl mx-auto space-y-8">
      
      <!-- Başlık Kartı -->
      <div class="bg-white p-8 rounded-3xl shadow-sm border border-gray-100 flex flex-col md:flex-row justify-between gap-6">
        <div>
          <div class="flex items-center gap-3 mb-2">
            <span class="px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-[10px] font-black uppercase tracking-tighter">#{{ ticketId }}</span>
            <span class="text-sm font-bold text-gray-400">• {{ ticket.supportType }}</span>
          </div>
          <h1 class="text-3xl font-black text-gray-900">{{ ticket.subject }}</h1>
          <p class="text-gray-500 mt-2 font-medium">Oluşturulma: {{ formatDate(ticket.createdAt) }}</p>
        </div>
        <div class="flex flex-col items-end gap-3">
          <span 
            class="px-6 py-2 rounded-2xl text-xs font-black uppercase tracking-widest shadow-sm"
            :class="{
              'bg-blue-600 text-white': ticket.status === 'Open' || ticket.status === 'Açık',
              'bg-yellow-500 text-white': ticket.status === 'Responded' || ticket.status === 'Beklemede',
              'bg-green-600 text-white': ticket.status === 'Closed' || ticket.status === 'Çözüldü'
            }"
          >
            {{ translateStatus(ticket.status) }}
          </span>
          <!-- Admin Durum Güncelleme -->
          <select 
            v-if="userRole === 'Admin'" 
            @change="updateStatus" 
            v-model="ticket.status"
            class="bg-gray-100 border-none rounded-xl px-4 py-2 text-xs font-bold focus:ring-2 focus:ring-blue-500 outline-none"
          >
            <option value="Açık">Açık Yap</option>
            <option value="Beklemede">Beklemeye Al</option>
            <option value="Çözüldü">Çözüldü İşaretle</option>
          </select>
        </div>
      </div>

      <!-- Mesajlaşma Alanı -->
      <div class="space-y-6">
        <!-- İlk Mesaj -->
        <div class="bg-white p-8 rounded-3xl border border-gray-100 shadow-sm relative overflow-hidden">
          <div class="absolute top-0 left-0 w-1 h-full bg-blue-500"></div>
          <div class="flex items-center gap-4 mb-4">
            <div class="w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center text-blue-600 font-bold">U</div>
            <div>
              <p class="font-bold text-gray-900">Siz (Talep Sahibi)</p>
              <p class="text-xs text-gray-400 font-medium">{{ formatDate(ticket.createdAt) }}</p>
            </div>
          </div>
          <p class="text-gray-700 leading-relaxed text-lg">{{ ticket.message }}</p>
        </div>

        <!-- Yazışmalar -->
        <div 
          v-for="msg in messages" 
          :key="msg.messageId" 
          class="flex flex-col"
          :class="msg.senderRole === 'Admin' ? 'items-start' : 'items-end'"
        >
          <div 
            class="max-w-[80%] p-6 rounded-3xl shadow-sm border border-gray-100 relative"
            :class="msg.senderRole === 'Admin' ? 'bg-indigo-600 text-white rounded-tl-none' : 'bg-white text-gray-700 rounded-tr-none'"
          >
            <div class="flex items-center gap-2 mb-2">
              <span class="text-[10px] font-black uppercase tracking-tighter" :class="msg.senderRole === 'Admin' ? 'text-indigo-200' : 'text-blue-500'">
                {{ msg.senderRole === 'Admin' ? 'Destek Ekibi' : 'Siz' }}
              </span>
              <span class="text-[10px] font-bold opacity-50" :class="msg.senderRole === 'Admin' ? 'text-white' : 'text-gray-400'">
                • {{ formatDate(msg.createdAt) }}
              </span>
            </div>
            <p class="text-base leading-relaxed">{{ msg.message }}</p>
          </div>
        </div>
      </div>

      <!-- Yanıt Yazma Alanı -->
      <div v-if="ticket.status !== 'Closed' && ticket.status !== 'Çözüldü'" class="bg-white p-6 rounded-3xl shadow-xl border border-gray-100 space-y-4">
        <textarea 
          v-model="newMessage"
          placeholder="Mesajınızı buraya yazın..."
          class="w-full p-6 bg-gray-50 rounded-2xl border-2 border-transparent focus:border-blue-500 outline-none transition-all min-h-[120px]"
        ></textarea>
        <div class="flex justify-end">
          <button 
            @click="sendMessage"
            :disabled="!newMessage"
            class="px-10 py-4 bg-blue-600 text-white rounded-2xl font-bold hover:bg-blue-700 disabled:opacity-50 transition-all flex items-center gap-2"
          >
            <span>✈️</span> Gönder
          </button>
        </div>
      </div>
      <div v-else class="text-center py-10 bg-gray-100 rounded-3xl border-2 border-dashed border-gray-200">
        <p class="text-xl font-bold text-gray-400">Bu talep çözüldüğü için yeni mesaj gönderilemez. 🔒</p>
      </div>

    </div>

    <!-- Yükleniyor -->
    <div v-else class="flex flex-col items-center justify-center h-[60vh]">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mb-4"></div>
      <p class="text-gray-500 font-bold">Talep detayları yükleniyor...</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useRoute } from 'vue-router';
import { jwtDecode } from 'jwt-decode';

const route = useRoute();
const ticketId = route.params.id;
const ticket = ref(null);
const messages = ref([]);
const newMessage = ref('');
const userRole = ref('');

const fetchTicketDetails = async () => {
  const token = localStorage.getItem('userToken');
  const config = { headers: { Authorization: `Bearer ${token}` } };
  
  try {
    // Tüm biletleri çekip içinden bulalım (Basitlik için)
    const endpoint = userRole.value === 'Admin' ? 'http://localhost:8080/admin/tickets' : 'http://localhost:8080/tickets';
    const response = await axios.get(endpoint, config);
    ticket.value = response.data.find(t => t.ticketId === parseInt(ticketId));
    
    // Mesajları çek
    const msgResponse = await axios.get(`http://localhost:8080/tickets/messages?ticketId=${ticketId}`, config);
    messages.value = msgResponse.data || [];
  } catch (error) {
    console.error("Detaylar yüklenemedi:", error);
  }
};

const sendMessage = async () => {
  const token = localStorage.getItem('userToken');
  const config = { headers: { Authorization: `Bearer ${token}` } };
  
  try {
    await axios.post('http://localhost:8080/tickets/messages/send', {
      ticketId: parseInt(ticketId),
      message: newMessage.value
    }, config);
    
    newMessage.value = '';
    fetchTicketDetails();
  } catch (error) {
    alert("Mesaj gönderilemedi.");
  }
};

const updateStatus = async () => {
  const token = localStorage.getItem('userToken');
  const config = { headers: { Authorization: `Bearer ${token}` } };
  
  try {
    await axios.post('http://localhost:8080/admin/tickets/update', {
      ticketId: parseInt(ticketId),
      status: ticket.value.status
    }, config);
    alert("Durum güncellendi.");
  } catch (error) {
    alert("Güncelleme hatası.");
  }
};

const translateStatus = (status) => {
  const map = {
    'Open': 'Açık', 'Açık': 'Açık',
    'Responded': 'Yanıtlandı', 'Beklemede': 'Yanıtlandı',
    'Closed': 'Çözüldü', 'Çözüldü': 'Çözüldü'
  };
  return map[status] || status;
};

const formatDate = (dateString) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('tr-TR', { 
    day: 'numeric', month: 'long', year: 'numeric',
    hour: '2-digit', minute: '2-digit' 
  }).format(date);
};

onMounted(() => {
  const token = localStorage.getItem('userToken');
  if (token) {
    const decoded = jwtDecode(token);
    userRole.value = decoded.role;
    fetchTicketDetails();
  }
});
</script>
