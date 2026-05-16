<template>
  <div class="min-h-screen bg-gray-50 py-12 px-6">
    <div class="max-w-6xl mx-auto">
      <div class="flex justify-between items-center mb-8">
        <h1 class="text-4xl font-black text-gray-900">Yönetici Paneli</h1>
        <button @click="router.back()" class="bg-white text-gray-700 px-6 py-2 rounded-xl font-bold shadow-sm border border-gray-200 hover:bg-gray-50">
          Geri Dön
        </button>
      </div>

      <!-- İstatistikler -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-10">
        <div class="bg-white p-6 rounded-3xl shadow-sm border border-gray-100 flex items-center gap-4">
          <div class="text-4xl">🎨</div>
          <div>
            <p class="text-gray-500 text-sm font-bold">Toplam Eser</p>
            <p class="text-3xl font-black text-gray-800">{{ stats.totalArtworks }}</p>
          </div>
        </div>
        <div class="bg-white p-6 rounded-3xl shadow-sm border border-gray-100 flex items-center gap-4">
          <div class="text-4xl">🎭</div>
          <div>
            <p class="text-gray-500 text-sm font-bold">Atölyeler</p>
            <p class="text-3xl font-black text-gray-800">{{ stats.totalWorkshops }}</p>
          </div>
        </div>
        <div class="bg-white p-6 rounded-3xl shadow-sm border border-gray-100 flex items-center gap-4">
          <div class="text-4xl">👥</div>
          <div>
            <p class="text-gray-500 text-sm font-bold">Kullanıcılar</p>
            <p class="text-3xl font-black text-gray-800">{{ stats.totalUsers }}</p>
          </div>
        </div>
        <div class="bg-white p-6 rounded-3xl shadow-sm border border-blue-200 flex items-center gap-4">
          <div class="text-4xl">🎫</div>
          <div>
            <p class="text-blue-500 text-sm font-bold">Açık Talepler</p>
            <p class="text-3xl font-black text-blue-800">{{ stats.activeTickets }}</p>
          </div>
        </div>
      </div>

      <!-- Tüm Destek Talepleri -->
      <div class="bg-white rounded-3xl shadow-sm border border-gray-100 overflow-hidden">
        <div class="p-6 border-b border-gray-100">
          <h3 class="text-2xl font-bold text-gray-800">Tüm Destek Talepleri</h3>
        </div>
        <table class="w-full text-left">
          <thead class="bg-gray-50">
            <tr>
              <th class="p-6 font-bold text-gray-600">ID</th>
              <th class="p-6 font-bold text-gray-600">Kullanıcı ID</th>
              <th class="p-6 font-bold text-gray-600">Konu</th>
              <th class="p-6 font-bold text-gray-600">Durum</th>
              <th class="p-6 font-bold text-gray-600">İşlem</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="ticket in tickets" :key="ticket.ticketId" class="border-b border-gray-50">
              <td class="p-6">#{{ ticket.ticketId }}</td>
              <td class="p-6">{{ ticket.userId }}</td>
              <td class="p-6 font-bold">{{ ticket.subject }}</td>
              <td class="p-6">
                <select v-model="ticket.status" @change="updateTicketStatus(ticket.ticketId, ticket.status)" class="p-2 border border-gray-200 rounded-lg outline-none">
                  <option value="Açık">Açık</option>
                  <option value="Beklemede">Beklemede</option>
                  <option value="Çözüldü">Çözüldü</option>
                </select>
              </td>
              <td class="p-6">
                <button @click="router.push(`/support/${ticket.ticketId}`)" class="text-blue-600 font-bold hover:underline">Mesajlar</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()
const stats = ref({})
const tickets = ref([])

const fetchDashboardData = async () => {
  try {
    const token = localStorage.getItem('userToken')
    const config = { headers: { Authorization: `Bearer ${token}` } }
    
    const [statsRes, ticketsRes] = await Promise.all([
      axios.get('http://localhost:8080/admin/dashboard-stats', config),
      axios.get('http://localhost:8080/admin/tickets', config)
    ])
    
    stats.value = statsRes.data
    tickets.value = ticketsRes.data || []
  } catch (error) {
    console.error("Dashboard verisi yüklenemedi", error)
  }
}

const updateTicketStatus = async (ticketId, status) => {
  try {
    const token = localStorage.getItem('userToken')
    await axios.post('http://localhost:8080/admin/tickets/update', { ticketId, status }, {
      headers: { Authorization: `Bearer ${token}` }
    })
    fetchDashboardData()
  } catch (error) {
    alert("Durum güncellenemedi.")
  }
}

onMounted(() => {
  if (!localStorage.getItem('userToken')) {
    router.push('/login')
  } else {
    fetchDashboardData()
  }
})
</script>