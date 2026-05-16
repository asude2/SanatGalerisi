<template>
  <div class="min-h-screen bg-gray-50 py-12 px-6">
    <div class="max-w-4xl mx-auto">
      <div class="flex justify-between items-center mb-8">
        <div>
          <span class="text-blue-600 font-bold uppercase text-xs tracking-widest">Destek Sistemi</span>
          <h1 class="text-4xl font-black text-gray-900 mt-1">Talep #{{ route.params.id }}</h1>
        </div>
        <button @click="router.back()" class="bg-white text-gray-700 px-6 py-2 rounded-xl font-bold shadow-sm border border-gray-200 hover:bg-gray-50 transition-all">
          Geri Dön
        </button>
      </div>

      <div class="bg-white rounded-3xl shadow-sm border border-gray-100 p-8 flex flex-col h-[70vh]">
        <!-- Mesajlaşma Alanı -->
        <div class="flex-1 overflow-y-auto space-y-6 mb-6 pr-4 custom-scrollbar">
          <div v-if="messages.length === 0" class="text-center text-gray-400 py-12">
            <span class="text-5xl block mb-4">💬</span>
            Henüz bir mesajlaşma başlatılmamış.
          </div>
          
          <div v-for="msg in messages" :key="msg.messageId" 
               class="flex flex-col" :class="msg.senderRole === 'Admin' || msg.senderRole === 'Instructor' ? 'items-start' : 'items-end'">
            <div class="max-w-[80%]">
              <div class="flex items-center gap-2 mb-1 px-2">
                <span class="text-xs text-gray-400 font-bold">
                  {{ msg.senderName }}
                </span>
                <span v-if="msg.senderRole === 'Admin' || msg.senderRole === 'Instructor'" class="bg-blue-100 text-blue-600 text-[10px] px-2 py-0.5 rounded-full font-black uppercase">
                  Yönetici
                </span>
              </div>
              <div class="p-4 rounded-2xl shadow-sm" 
                   :class="msg.senderRole === 'Admin' || msg.senderRole === 'Instructor' ? 'bg-blue-50 text-blue-900 rounded-tl-sm' : 'bg-galeri-yesil text-white rounded-tr-sm'">
                <p class="whitespace-pre-wrap leading-relaxed">{{ msg.message }}</p>
              </div>
              <span class="text-[10px] text-gray-400 mt-1 px-2 block" :class="msg.senderRole === 'Admin' || msg.senderRole === 'Instructor' ? 'text-left' : 'text-right'">
                {{ new Date(msg.createdAt).toLocaleString('tr-TR') }}
              </span>
            </div>
          </div>
        </div>

        <!-- Mesaj Yazma Alanı -->
        <div class="pt-6 border-t border-gray-100 flex gap-4 items-end">
          <div class="flex-1 relative">
            <textarea 
              v-model="newMessage" 
              rows="2" 
              class="w-full p-4 rounded-2xl border-2 border-gray-50 focus:border-galeri-yesil focus:bg-white outline-none resize-none transition-all bg-gray-50 text-gray-700" 
              placeholder="Mesajınızı buraya yazın..."
              @keyup.enter.ctrl="sendMessage"
            ></textarea>
            <span class="absolute right-4 bottom-2 text-[10px] text-gray-300">Ctrl + Enter ile gönder</span>
          </div>
          <button 
            @click="sendMessage" 
            :disabled="!newMessage.trim()"
            class="bg-blue-600 text-white p-4 rounded-2xl font-bold hover:bg-blue-700 transition-all shadow-lg shadow-blue-100 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span class="text-xl">🚀</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const router = useRouter()
const messages = ref([])
const newMessage = ref('')
let pollInterval = null

const fetchMessages = async () => {
  try {
    const token = localStorage.getItem('userToken')
    const res = await axios.get(`http://localhost:8080/tickets/messages?ticketId=${route.params.id}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    messages.value = res.data || []
    
    // Mesajlar gelince en alta kaydır
    setTimeout(() => {
      const container = document.querySelector('.custom-scrollbar')
      if (container) container.scrollTop = container.scrollHeight
    }, 100)
  } catch (error) {
    console.error("Mesajlar yüklenemedi", error)
  }
}

const sendMessage = async () => {
  if (!newMessage.value.trim()) return

  try {
    const token = localStorage.getItem('userToken')
    await axios.post('http://localhost:8080/tickets/messages/send', {
      ticketId: parseInt(route.params.id),
      message: newMessage.value.trim()
    }, {
      headers: { Authorization: `Bearer ${token}` }
    })
    newMessage.value = ''
    fetchMessages()
  } catch (error) {
    alert("Mesaj gönderilemedi.")
  }
}

onMounted(() => {
  fetchMessages()
  // Canlı destek hissi için her 5 saniyede bir yeni mesajları kontrol et
  pollInterval = setInterval(fetchMessages, 5000)
})

onUnmounted(() => {
  if (pollInterval) clearInterval(pollInterval)
})
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #e5e7eb;
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #d1d5db;
}
</style>