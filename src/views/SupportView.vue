<template>
  <div class="min-h-screen bg-gray-50 pt-24 pb-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-6xl mx-auto space-y-8">
      
      <!-- Header -->
      <div class="text-center mb-12">
        <h1 class="text-4xl font-black text-gray-900 mb-4">Destek Merkezi 🎧</h1>
        <p class="text-xl text-gray-500 max-w-2xl mx-auto">Sanat Galerisi ile ilgili bir sorununuz mu var? Bize bir bilet oluşturun veya anında canlı desteğe bağlanın.</p>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        
        <!-- Sol Taraf: Bilet Oluşturma ve Geçmiş Biletler -->
        <div class="lg:col-span-2 space-y-8">
          
          <!-- Yeni Destek Talebi Formu -->
          <div class="bg-white p-8 rounded-3xl shadow-xl border border-gray-100 relative overflow-hidden">
            <div class="absolute top-0 right-0 w-32 h-32 bg-blue-500/10 rounded-bl-[100px] -z-10"></div>
            <h2 class="text-2xl font-bold text-gray-800 mb-6 flex items-center gap-2">
              <span class="text-3xl">🎫</span> Yeni Destek Talebi
            </h2>
            
            <form @submit.prevent="submitTicket" class="space-y-5">
              <div>
                <label class="block text-sm font-bold text-gray-700 mb-2">Konu</label>
                <select v-model="newTicket.subject" required class="w-full p-4 bg-gray-50 rounded-xl border border-gray-200 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all">
                  <option value="" disabled>Lütfen bir konu seçin</option>
                  <option value="Sipariş & Teslimat">Sipariş & Teslimat</option>
                  <option value="Eser İadesi">Eser İadesi</option>
                  <option value="Atölye İptali">Atölye İptali</option>
                  <option value="Teknik Sorun">Teknik Sorun</option>
                  <option value="Diğer">Diğer</option>
                </select>
              </div>

              <div>
                <label class="block text-sm font-bold text-gray-700 mb-2">Mesajınız</label>
                <textarea 
                  v-model="newTicket.message" 
                  required
                  placeholder="Yaşadığınız sorunu detaylıca anlatın..."
                  class="w-full p-4 bg-gray-50 rounded-xl border border-gray-200 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all min-h-[150px] resize-none"
                ></textarea>
              </div>

              <div class="flex justify-end">
                <button 
                  type="submit" 
                  :disabled="isSubmitting"
                  class="bg-blue-600 text-white px-8 py-4 rounded-xl font-bold hover:bg-blue-700 transition-all shadow-lg shadow-blue-200 flex items-center gap-2 disabled:opacity-50"
                >
                  <span v-if="isSubmitting" class="animate-spin text-xl">⏳</span>
                  <span>{{ isSubmitting ? 'Gönderiliyor...' : 'Talebi Oluştur' }}</span>
                </button>
              </div>
            </form>
          </div>

          <!-- Geçmiş Talepler -->
          <div class="bg-white p-8 rounded-3xl shadow-lg border border-gray-100">
            <h2 class="text-xl font-bold text-gray-800 mb-6 border-b pb-4">Talepleriniz</h2>
            
            <div v-if="loading" class="flex justify-center py-6">
              <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            </div>

            <div v-else-if="tickets.length === 0" class="text-center py-8 text-gray-400 font-medium">
              Henüz oluşturulmuş bir destek talebiniz bulunmamaktadır.
            </div>

            <div v-else class="space-y-4">
              <div v-for="ticket in tickets" :key="ticket.id" class="p-5 rounded-2xl border border-gray-100 hover:bg-gray-50 transition-colors flex justify-between items-start">
                <div>
                  <h3 class="font-bold text-gray-800">{{ ticket.subject }}</h3>
                  <p class="text-gray-500 text-sm mt-1 line-clamp-2 max-w-lg">{{ ticket.message }}</p>
                  <span class="text-xs text-gray-400 mt-2 block">{{ new Date(ticket.createdAt).toLocaleString('tr-TR') }}</span>
                </div>
                <div class="px-3 py-1 rounded-lg text-sm font-bold" 
                     :class="ticket.status === 'Open' ? 'bg-yellow-100 text-yellow-700' : 'bg-green-100 text-green-700'">
                  {{ ticket.status === 'Open' ? 'Açık' : 'Çözüldü' }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Sağ Taraf: Canlı Destek Chatbot -->
        <div class="lg:col-span-1 h-[600px] bg-white rounded-3xl shadow-xl border border-gray-100 flex flex-col overflow-hidden sticky top-24">
          <div class="bg-galeri-yesil text-white p-6 flex items-center gap-3 shadow-md z-10">
            <div class="relative">
              <div class="w-12 h-12 bg-white rounded-full flex items-center justify-center text-2xl">🤖</div>
              <div class="absolute bottom-0 right-0 w-3 h-3 bg-green-300 rounded-full border-2 border-galeri-yesil"></div>
            </div>
            <div>
              <h3 class="font-black text-lg">SanatBot</h3>
              <p class="text-xs text-green-100">Çevrimiçi | Hızlı Yanıt</p>
            </div>
          </div>

          <!-- Chat Messages -->
          <div class="flex-1 p-6 overflow-y-auto bg-gray-50 space-y-4" id="chatbox">
            <div v-for="(msg, idx) in chatMessages" :key="idx" class="flex" :class="msg.sender === 'user' ? 'justify-end' : 'justify-start'">
              <div 
                class="max-w-[80%] p-4 rounded-2xl text-sm"
                :class="msg.sender === 'user' ? 'bg-blue-600 text-white rounded-tr-none' : 'bg-white text-gray-700 border border-gray-200 rounded-tl-none shadow-sm'"
              >
                {{ msg.text }}
              </div>
            </div>
          </div>

          <!-- Chat Input -->
          <div class="p-4 bg-white border-t border-gray-100">
            <form @submit.prevent="sendChatMessage" class="flex gap-2">
              <input 
                v-model="currentChatMsg"
                type="text" 
                placeholder="Bir mesaj yazın..." 
                class="flex-1 bg-gray-50 border border-gray-200 rounded-xl px-4 py-3 focus:outline-none focus:ring-2 focus:ring-galeri-yesil text-sm"
              />
              <button type="submit" :disabled="!currentChatMsg.trim()" class="bg-galeri-yesil text-white px-4 rounded-xl hover:bg-green-600 transition-colors disabled:opacity-50 flex items-center justify-center">
                ➤
              </button>
            </form>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue';
import { jwtDecode } from 'jwt-decode';
import { useRouter } from 'vue-router';
import axios from 'axios';

const router = useRouter();
const tickets = ref([]);
const loading = ref(false);
const isSubmitting = ref(false);

const newTicket = ref({
  subject: '',
  message: ''
});

// Chatbot State
const currentChatMsg = ref('');
const chatMessages = ref([
  { sender: 'bot', text: 'Merhaba! 👋 Ben SanatBot. Eserler, atölyeler veya kargo durumu hakkında sorularınızı yanıtlayabilirim.' }
]);

const fetchTickets = async () => {
  const token = localStorage.getItem('userToken');
  if (!token) return;

  try {
    loading.value = true;
    const decoded = jwtDecode(token);
    const response = await axios.get(`http://localhost:8080/support/tickets?email=${decoded.email}`);
    tickets.value = response.data || [];
  } catch (error) {
    console.error("Biletler yüklenemedi:", error);
  } finally {
    loading.value = false;
  }
};

const submitTicket = async () => {
  const token = localStorage.getItem('userToken');
  if (!token) {
    alert("Destek talebi oluşturmak için giriş yapmalısınız.");
    router.push('/login');
    return;
  }

  try {
    isSubmitting.value = true;
    const decoded = jwtDecode(token);
    
    await axios.post('http://localhost:8080/support/create', {
      userEmail: decoded.email,
      subject: newTicket.value.subject,
      message: newTicket.value.message
    });

    alert("Destek talebiniz başarıyla oluşturuldu! Müşteri temsilcimiz en kısa sürede size dönüş yapacaktır.");
    newTicket.value.subject = '';
    newTicket.value.message = '';
    fetchTickets(); // Listeyi güncelle
    
  } catch (error) {
    alert("Talep oluşturulurken bir hata meydana geldi.");
    console.error(error);
  } finally {
    isSubmitting.value = false;
  }
};

// Basit Chatbot Mantığı
const sendChatMessage = async () => {
  if (!currentChatMsg.value.trim()) return;

  const userText = currentChatMsg.value;
  chatMessages.value.push({ sender: 'user', text: userText });
  currentChatMsg.value = '';

  scrollToBottom();

  // Bot Typing Delay Simulation
  setTimeout(() => {
    let botReply = "Üzgünüm, sizi tam olarak anlayamadım. Sol taraftan bir destek bileti oluşturursanız ekibimiz detaylı ilgilenecektir.";
    const lowerText = userText.toLowerCase();

    if (lowerText.includes("kargo") || lowerText.includes("teslimat") || lowerText.includes("nerede")) {
      botReply = "Siparişleriniz genellikle 3-5 iş günü içinde kargoya teslim edilmektedir. Sipariş numaranızı yazar mısınız?";
    } else if (lowerText.includes("iade") || lowerText.includes("iptal")) {
      botReply = "Eserlerimizde 14 gün içerisinde iade hakkınız bulunmaktadır. İade kodu almak için bir destek talebi oluşturabilirsiniz.";
    } else if (lowerText.includes("merhaba") || lowerText.includes("selam")) {
      botReply = "Merhaba! Size nasıl yardımcı olabilirim? 😊";
    }

    chatMessages.value.push({ sender: 'bot', text: botReply });
    scrollToBottom();
  }, 1000);
};

const scrollToBottom = () => {
  nextTick(() => {
    const chatbox = document.getElementById('chatbox');
    if (chatbox) {
      chatbox.scrollTop = chatbox.scrollHeight;
    }
  });
};

onMounted(() => {
  const token = localStorage.getItem('userToken');
  if (!token) {
    alert("Bu sayfayı görüntülemek için giriş yapmalısınız.");
    router.push('/login');
    return;
  }
  fetchTickets();
});
</script>
