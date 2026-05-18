<template>
  <div class="mt-12 space-y-8">
    <div class="flex items-center justify-between border-b border-gray-100 pb-6">
      <h2 class="text-3xl font-black text-gray-900 flex items-center gap-3">
        💬 Değerlendirmeler <span class="text-blue-600 text-lg bg-blue-50 px-3 py-1 rounded-full">{{ comments.length }}</span>
      </h2>
      
      <!-- Filtreleme -->
      <div class="flex items-center gap-4">
        <label class="text-sm font-bold text-gray-400 uppercase tracking-tighter">Sırala:</label>
        <select v-model="sortBy" @change="fetchComments" class="bg-white border border-gray-200 rounded-xl px-4 py-2 text-sm font-semibold focus:ring-2 focus:ring-blue-500 outline-none transition-all">
          <option value="newest">En Yeni</option>
          <option value="highest">En Yüksek Puan</option>
          <option value="most_helpful">En Faydalı</option>
        </select>
      </div>
    </div>

    <!-- Ortalama Puan Özeti -->
    <div v-if="comments.length > 0" class="bg-white p-8 rounded-3xl border border-gray-100 shadow-sm flex flex-col md:flex-row items-center gap-8">
      <div class="text-center md:border-r md:pr-12 border-gray-100">
        <p class="text-6xl font-black text-gray-900">{{ averageRating }}</p>
        <div class="flex gap-1 justify-center my-2">
          <span v-for="i in 5" :key="i" :class="i <= Math.round(averageRating) ? 'text-yellow-400' : 'text-gray-200'" class="text-2xl">★</span>
        </div>
        <p class="text-gray-400 text-sm font-bold uppercase tracking-widest">Genel Puan</p>
      </div>
      <div class="flex-1 w-full space-y-3">
        <div v-for="star in [5, 4, 3, 2, 1]" :key="star" class="flex items-center gap-4">
          <span class="text-sm font-bold text-gray-500 w-4">{{ star }}</span>
          <div class="flex-1 bg-gray-100 h-3 rounded-full overflow-hidden">
            <div 
              class="bg-yellow-400 h-full rounded-full transition-all duration-1000" 
              :style="{ width: `${getStarPercentage(star)}%` }"
            ></div>
          </div>
          <span class="text-sm font-bold text-gray-400 w-10 text-right">{{ getStarCount(star) }}</span>
        </div>
      </div>
    </div>

    <!-- Yorum Yapma Formu -->
    <div class="bg-blue-50/50 p-8 rounded-3xl border-2 border-dashed border-blue-100">
      <h3 class="text-xl font-bold text-gray-800 mb-6 flex items-center gap-2">
        ✨ Senin Deneyimin Nasıl?
      </h3>
      
      <div class="space-y-6">
        <div>
          <p class="text-sm font-bold text-gray-500 mb-3 ml-1 uppercase">Puanın</p>
          <div class="flex gap-3">
            <button 
              v-for="i in 5" 
              :key="i" 
              @click="newComment.rating = i"
              class="text-4xl transition-all hover:scale-125 cursor-pointer"
              :class="i <= newComment.rating ? 'grayscale-0' : 'grayscale opacity-30'"
            >
              {{ i === 1 ? '😡' : i === 2 ? '🙁' : i === 3 ? '😐' : i === 4 ? '😊' : '🤩' }}
            </button>
          </div>
        </div>

        <div class="relative">
          <textarea 
            v-model="newComment.text"
            placeholder="Görüşlerini buraya yazabilirsin..."
            class="w-full p-6 bg-white rounded-2xl border-2 border-transparent focus:border-blue-500 outline-none transition-all shadow-sm min-h-[150px] text-lg"
          ></textarea>
        </div>

        <button 
          @click="submitComment"
          :disabled="!newComment.rating || !newComment.text"
          class="w-full md:w-auto px-10 py-4 bg-blue-600 text-white rounded-2xl font-bold text-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-lg shadow-blue-100"
        >
          Yorumu Paylaş 🚀
        </button>
      </div>
    </div>

    <!-- Yorum Listesi -->
    <div class="space-y-6">
      <div v-if="comments.length === 0" class="text-center py-20 bg-gray-50 rounded-3xl border-2 border-dotted border-gray-200">
        <p class="text-6xl mb-4">🙊</p>
        <p class="text-xl font-bold text-gray-400">Henüz hiç yorum yapılmamış.</p>
        <p class="text-gray-400 mt-2">İlk yorumu sen yaparak diğerlerine rehberlik et!</p>
      </div>

      <div 
        v-for="comment in comments" 
        :key="comment.commentId" 
        class="bg-white p-8 rounded-3xl border border-gray-100 shadow-sm hover:shadow-md transition-all group"
      >
        <div class="flex flex-col md:flex-row justify-between gap-6">
          <div class="flex-1 space-y-4">
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 bg-gradient-to-tr from-blue-500 to-indigo-600 rounded-2xl flex items-center justify-center text-white font-bold text-xl shadow-lg">
                {{ comment.userName.charAt(0) }}
              </div>
              <div>
                <div class="flex items-center gap-2">
                  <p class="font-black text-gray-900">{{ comment.userName }}</p>
                  <span v-if="comment.isVerified" class="bg-green-100 text-green-700 text-[10px] font-black uppercase px-2 py-0.5 rounded-full flex items-center gap-1">
                    ✓ Doğrulanmış {{ targetType === 'Artwork' ? 'Alıcı' : 'Katılımcı' }}
                  </span>
                </div>
                <div class="flex gap-0.5 mt-0.5">
                  <span v-for="i in 5" :key="i" :class="i <= comment.rating ? 'text-yellow-400' : 'text-gray-200'" class="text-xs">★</span>
                  <span class="text-xs text-gray-400 ml-2 font-medium">{{ formatDate(comment.createdAt) }}</span>
                </div>
              </div>
            </div>

            <p class="text-gray-700 leading-relaxed text-lg">{{ comment.commentText }}</p>

            <!-- Admin Yanıtı -->
            <div v-if="comment.adminReply" class="bg-gray-50 p-6 rounded-2xl border-l-4 border-blue-500 mt-4">
              <div class="flex items-center gap-2 mb-2">
                <span class="text-sm font-black text-blue-600 uppercase tracking-tighter">
                  {{ comment.replierRole === 'Admin' ? 'Galeri Yönetimi' : comment.replierRole === 'Instructor' ? 'Eğitmen' : 'Doğrulanmış Yanıt' }}
                </span>
                <span class="text-xs text-gray-400 font-bold">• {{ comment.replierName }}</span>
              </div>
              <p class="text-gray-600 italic">"{{ comment.adminReply }}"</p>
            </div>

            <!-- Admin Yanıt Verme Formu (Sadece Admin/Instructor için) -->
            <div v-if="(userRole === 'Admin' || userRole === 'Instructor') && !comment.adminReply" class="mt-4">
              <button 
                v-if="!replyingTo[comment.commentId]" 
                @click="replyingTo[comment.commentId] = true"
                class="text-sm font-bold text-blue-600 hover:text-blue-800 flex items-center gap-1"
              >
                ↩ Yanıtla
              </button>
              <div v-else class="space-y-3">
                <textarea 
                  v-model="replies[comment.commentId]"
                  placeholder="Yanıtınızı yazın..."
                  class="w-full p-4 bg-gray-50 rounded-xl border border-gray-200 outline-none focus:border-blue-500 transition-all text-sm"
                ></textarea>
                <div class="flex gap-2">
                  <button @click="submitReply(comment.commentId)" class="bg-blue-600 text-white px-4 py-2 rounded-lg text-xs font-bold">Gönder</button>
                  <button @click="replyingTo[comment.commentId] = false" class="bg-gray-200 text-gray-600 px-4 py-2 rounded-lg text-xs font-bold">İptal</button>
                </div>
              </div>
            </div>
          </div>

          <!-- Oy Verme -->
          <div class="flex md:flex-col items-center gap-3 bg-gray-50/50 p-4 rounded-2xl h-fit">
            <button 
              @click="vote(comment.commentId, 'Up')"
              class="p-2 hover:bg-white rounded-xl transition-all hover:scale-110 active:scale-90"
              title="Faydalı"
            >
              👍 <span class="text-xs font-bold text-gray-500">{{ comment.upvotes }}</span>
            </button>
            <div class="w-px h-4 md:w-4 md:h-px bg-gray-200"></div>
            <button 
              @click="vote(comment.commentId, 'Down')"
              class="p-2 hover:bg-white rounded-xl transition-all hover:scale-110 active:scale-90"
              title="Faydalı değil"
            >
              👎 <span class="text-xs font-bold text-gray-500">{{ comment.downvotes }}</span>
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

const props = defineProps({
  targetId: { type: Number, required: true },
  targetType: { type: String, required: true } // 'Artwork' veya 'Workshop'
});

const comments = ref([]);
const sortBy = ref('newest');
const userRole = ref('');
const replyingTo = ref({});
const replies = ref({});

const newComment = ref({
  rating: 0,
  text: ''
});

const averageRating = computed(() => {
  if (comments.value.length === 0) return 0;
  const sum = comments.value.reduce((acc, curr) => acc + curr.rating, 0);
  return (sum / comments.value.length).toFixed(1);
});

const fetchComments = async () => {
  try {
    const response = await axios.get(`http://localhost:8080/comments?targetId=${props.targetId}&targetType=${props.targetType}&sort=${sortBy.value}`);
    comments.value = response.data || [];
  } catch (error) {
    console.error("Yorumlar yüklenemedi:", error);
  }
};

const getStarCount = (star) => {
  return comments.value.filter(c => c.rating === star).length;
};

const getStarPercentage = (star) => {
  if (comments.value.length === 0) return 0;
  return (getStarCount(star) / comments.value.length) * 100;
};

const submitComment = async () => {
  const token = localStorage.getItem('userToken');
  if (!token) {
    alert("Lütfen önce giriş yapın! 👤");
    return;
  }

  try {
    const config = { headers: { Authorization: `Bearer ${token}` } };
    await axios.post('http://localhost:8080/comments/add', {
      targetId: props.targetId,
      targetType: props.targetType,
      commentText: newComment.value.text,
      rating: newComment.value.rating
    }, config);

    newComment.value = { rating: 0, text: '' };
    fetchComments();
    alert("Yorumunuz başarıyla paylaşıldı! ✨");
  } catch (error) {
    alert(error.response?.data || "Yorum gönderilirken bir hata oluştu.");
  }
};

const vote = async (commentId, voteType) => {
  const token = localStorage.getItem('userToken');
  if (!token) {
    alert("Oy vermek için lütfen giriş yapın! 👤");
    return;
  }

  try {
    const config = { headers: { Authorization: `Bearer ${token}` } };
    await axios.post('http://localhost:8080/comments/vote', { commentId, voteType }, config);
    fetchComments();
  } catch (error) {
    alert(error.response?.data || "Oy verme hatası.");
  }
};

const submitReply = async (commentId) => {
  const token = localStorage.getItem('userToken');
  const text = replies.value[commentId];
  if (!text) return;

  try {
    const config = { headers: { Authorization: `Bearer ${token}` } };
    await axios.post('http://localhost:8080/comments/reply', { commentId, replyText: text }, config);
    
    replyingTo.value[commentId] = false;
    replies.value[commentId] = '';
    fetchComments();
  } catch (error) {
    alert(error.response?.data || "Yanıt gönderilirken bir hata oluştu.");
  }
};

const formatDate = (dateString) => {
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('tr-TR', { day: 'numeric', month: 'long', year: 'numeric' }).format(date);
};

onMounted(() => {
  const token = localStorage.getItem('userToken');
  if (token) {
    try {
      const decoded = jwtDecode(token);
      userRole.value = decoded.role;
    } catch (e) {
      console.error("Token decode hatası:", e);
    }
  }
  fetchComments();
});
</script>
