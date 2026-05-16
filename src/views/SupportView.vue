<template>
  <div class="min-h-screen bg-gray-50 py-12 px-6">
    <div class="max-w-5xl mx-auto">
      <div class="flex justify-between items-center mb-8">
        <h1 class="text-4xl font-black text-gray-900">Destek Taleplerim</h1>
        <button @click="router.back()" class="bg-white text-gray-700 px-6 py-2 rounded-xl font-bold shadow-sm border border-gray-200 hover:bg-gray-50">
          Geri Dön
        </button>
      </div>

      <!-- Yeni Talep -->
      <div class="bg-white p-8 rounded-3xl shadow-sm border border-gray-100 mb-8">
        <h3 class="text-2xl font-bold text-gray-800 mb-6">Yeni Destek Talebi Oluştur</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
          <div>
            <label class="block text-sm font-bold text-gray-700 mb-2">Konu</label>
            <input v-model="newTicket.subject" type="text" class="w-full p-3 rounded-xl border border-gray-200 focus:border-galeri-yesil outline-none" />
          </div>
          <div>
            <label class="block text-sm font-bold text-gray-700 mb-2">Destek Türü</label>
            <select v-model="newTicket.supportType" class="w-full p-3 rounded-xl border border-gray-200 focus:border-galeri-yesil outline-none">
              <option value="Satın Alma">Satın Alma İşlemleri</option>
              <option value="Atölye">Atölye & Etkinlik</option>
              <option value="Teknik">Teknik Sorun</option>
              <option value="Diğer">Diğer</option>
            </select>
          </div>
        </div>
        <div class="mb-6">
          <label class="block text-sm font-bold text-gray-700 mb-2">Mesajınız</label>
          <textarea v-model="newTicket.message" rows="4" class="w-full p-3 rounded-xl border border-gray-200 focus:border-galeri-yesil outline-none"></textarea>
        </div>
        <button @click="createTicket" class="bg-galeri-yesil text-white px-8 py-3 rounded-xl font-bold hover:bg-green-700 transition-colors">
          Talebi Gönder
        </button>
      </div>

      <!-- Talepler Listesi -->
      <div class="bg-white rounded-3xl shadow-sm border border-gray-100 overflow-hidden">
        <table class="w-full text-left">
          <thead class="bg-gray-50 border-b border-gray-100">
            <tr>
              <th class="p-6 font-bold text-gray-600">ID</th>
              <th class="p-6 font-bold text-gray-600">Konu</th>
              <th class="p-6 font-bold text-gray-600">Tür</th>
              <th class="p-6 font-bold text-gray-600">Durum</th>
              <th class="p-6 font-bold text-gray-600">İşlem</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="ticket in tickets" :key="ticket.ticketId" class="border-b border-gray-50">
              <td class="p-6">#{{ ticket.ticketId }}</td>
              <td class="p-6 font-bold">{{ ticket.subject }}</td>
              <td class="p-6">{{ ticket.supportType }}</td>
              <td class="p-6">{{ ticket.status }}</td>
              <td class="p-6">
                <button @click="router.push(`/support/${ticket.ticketId}`)" class="text-blue-600 font-bold hover:underline">Detay</button>
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
const tickets = ref([])
const newTicket = ref({ subject: '', supportType: 'Satın Alma', message: '' })

const fetchTickets = async () => {
  try {
    const token = localStorage.getItem('userToken')
    const res = await axios.get('http://localhost:8080/tickets', {
      headers: { Authorization: `Bearer ${token}` }
    })
    tickets.value = res.data || []
  } catch (error) { console.error(error) }
}

const createTicket = async () => {
  try {
    const token = localStorage.getItem('userToken')
    await axios.post('http://localhost:8080/tickets/create', newTicket.value, {
      headers: { Authorization: `Bearer ${token}` }
    })
    newTicket.value.subject = ''; newTicket.value.message = ''; fetchTickets()
  } catch (error) { alert("Hata") }
}

onMounted(() => {
  if (!localStorage.getItem('userToken')) router.push('/login')
  else fetchTickets()
})
</script>