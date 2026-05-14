<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <header class="bg-galeri-yesil text-white shadow-2xl rounded-b-[50px] p-8 pb-12">
      <div class="container mx-auto flex justify-between items-center">
        <div>
          <h1 class="text-4xl font-black italic tracking-tighter">🎧 Destek Merkezi</h1>
          <p class="opacity-90 mt-2 text-lg">Sana yardımcı olmak için buradayız.</p>
        </div>
        <button @click="router.push('/')" class="bg-white text-galeri-yesil px-5 py-3 rounded-2xl font-bold shadow hover:scale-105 transition-transform cursor-pointer">
          ← Ana Sayfa
        </button>
      </div>
    </header>

    <div class="max-w-4xl mx-auto px-6 py-10 space-y-10">

      <!-- Yeni Talep Formu -->
      <div class="bg-white rounded-3xl shadow-xl p-8">
        <h2 class="text-2xl font-black text-gray-800 mb-6">📩 Yeni Destek Talebi</h2>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-bold text-gray-600 mb-2">Konu</label>
            <input
              v-model="newTicket.subject"
              type="text"
              placeholder="Sorununuzu kısaca özetleyin..."
              class="w-full border-2 border-gray-200 rounded-2xl px-4 py-3 focus:outline-none focus:border-galeri-yesil transition-colors"
            />
          </div>
          <div>
            <label class="block text-sm font-bold text-gray-600 mb-2">Mesaj</label>
            <textarea
              v-model="newTicket.message"
              placeholder="Sorununuzu detaylı olarak açıklayın..."
              rows="5"
              class="w-full border-2 border-gray-200 rounded-2xl px-4 py-3 focus:outline-none focus:border-galeri-yesil transition-colors resize-none"
            ></textarea>
          </div>
          <button
            @click="submitTicket"
            :disabled="!newTicket.subject || !newTicket.message"
            class="w-full bg-galeri-yesil text-white font-bold py-4 rounded-2xl shadow-lg hover:opacity-90 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
          >
            {{ sending ? 'Gönderiliyor...' : '📤 Talebi Gönder' }}
          </button>
          <p v-if="successMsg" class="text-center text-green-600 font-bold">{{ successMsg }}</p>
        </div>
      </div>

      <!-- Geçmiş Talepler -->
      <div class="bg-white rounded-3xl shadow-xl p-8">
        <h2 class="text-2xl font-black text-gray-800 mb-6">📋 Geçmiş Taleplerim</h2>
        <div v-if="tickets.length === 0" class="text-center text-gray-400 py-10">
          <p class="text-5xl mb-4">📭</p>
          <p class="text-lg">Henüz bir destek talebiniz yok.</p>
        </div>
        <div v-else class="space-y-4">
          <div
            v-for="ticket in tickets"
            :key="ticket.id"
            class="border-2 border-gray-100 rounded-2xl p-5 hover:border-galeri-yesil transition-colors"
          >
            <div class="flex justify-between items-start mb-2">
              <h3 class="font-bold text-gray-800 text-lg">{{ ticket.subject }}</h3>
              <span
                :class="{
                  'bg-yellow-100 text-yellow-700': ticket.status === 'Beklemede',
                  'bg-blue-100 text-blue-700': ticket.status === 'İşlemde',
                  'bg-green-100 text-green-700': ticket.status === 'Çözüldü'
                }"
                class="px-3 py-1 rounded-full text-sm font-bold"
              >
                {{ ticket.status }}
              </span>
            </div>
            <p class="text-gray-600 text-sm mb-2">{{ ticket.message }}</p>
            <p class="text-gray-400 text-xs">{{ ticket.createdAt }}</p>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';

const router = useRouter();
const email = localStorage.getItem('userEmail');

const newTicket = ref({ subject: '', message: '' });
const tickets = ref([]);
const sending = ref(false);
const successMsg = ref('');

const fetchTickets = async () => {
  try {
    const res = await axios.get(`http://localhost:8080/support/tickets?email=${email}`);
    tickets.value = res.data || [];
  } catch (e) {
    console.error('Talepler yüklenemedi:', e);
  }
};

const submitTicket = async () => {
  if (!newTicket.value.subject || !newTicket.value.message) return;
  sending.value = true;
  try {
    await axios.post('http://localhost:8080/support/create', {
      userEmail: email,
      subject: newTicket.value.subject,
      message: newTicket.value.message,
    });
    successMsg.value = '✅ Talebiniz başarıyla gönderildi!';
    newTicket.value = { subject: '', message: '' };
    await fetchTickets();
    setTimeout(() => (successMsg.value = ''), 3000);
  } catch (e) {
    console.error('Gönderme hatası:', e);
  } finally {
    sending.value = false;
  }
};

onMounted(fetchTickets);
</script>
