<template>
  <div class="mt-16 bg-white rounded-3xl shadow-xl border border-gray-100 overflow-hidden">
    <!-- Başlık -->
    <div class="bg-gray-50 border-b border-gray-100 p-8 flex justify-between items-center">
      <div>
        <h2 class="text-2xl font-black text-gray-800">Değerlendirmeler ve Yorumlar</h2>
        <p class="text-gray-500 text-sm mt-1">Bu {{ targetType === 'Artwork' ? 'eseri' : 'atölyeyi' }} deneyimleyenlerin fikirleri</p>
      </div>
      <div class="flex items-center gap-2 bg-white px-4 py-2 rounded-xl shadow-sm border border-gray-100" v-if="comments.length > 0">
        <span class="text-yellow-400 text-xl">⭐</span>
        <span class="font-bold text-gray-800 text-lg">{{ averageRating }}</span>
        <span class="text-gray-400 text-sm">({{ comments.length }} Yorum)</span>
      </div>
    </div>

    <div class="p-8">
      <!-- Yorum Ekleme Formu -->
      <div class="mb-10 bg-blue-50/50 p-6 rounded-2xl border border-blue-100">
        <h3 class="text-lg font-bold text-blue-900 mb-4">Senin Düşüncen Nedir?</h3>
        
        <form @submit.prevent="submitComment" class="space-y-4">
          <!-- Yıldız Seçimi -->
          <div class="flex items-center gap-2">
            <span class="text-sm font-semibold text-gray-600 mr-2">Puanın:</span>
            <div class="flex gap-1 cursor-pointer">
              <span 
                v-for="star in 5" 
                :key="star"
                @click="newComment.rating = star"
                @mouseenter="hoverRating = star"
                @mouseleave="hoverRating = 0"
                class="text-3xl transition-transform hover:scale-110"
                :class="(hoverRating ? star <= hoverRating : star <= newComment.rating) ? 'text-yellow-400 grayscale-0' : 'text-gray-300 grayscale'"
              >
                ⭐
              </span>
            </div>
          </div>

          <!-- Yorum Alanı -->
          <textarea 
            v-model="newComment.content" 
            placeholder="Harika bir deneyimdi, herkese tavsiye ederim..." 
            class="w-full p-4 bg-white rounded-xl border border-blue-200 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all min-h-[120px] resize-none"
            required
          ></textarea>

          <div class="flex justify-end">
            <button 
              type="submit" 
              :disabled="isSubmitting || newComment.rating === 0"
              class="bg-blue-600 text-white px-8 py-3 rounded-xl font-bold hover:bg-blue-700 transition-all disabled:opacity-50 disabled:cursor-not-allowed shadow-md shadow-blue-200 flex items-center gap-2"
            >
              <span v-if="isSubmitting" class="animate-spin text-xl">⏳</span>
              <span>{{ isSubmitting ? 'Gönderiliyor...' : 'Yorumu Paylaş' }}</span>
            </button>
          </div>
        </form>
      </div>

      <!-- Yorum Listesi -->
      <div v-if="loading" class="flex justify-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-4 border-gray-400"></div>
      </div>

      <div v-else-if="comments.length === 0" class="text-center py-12 bg-gray-50 rounded-2xl border border-dashed border-gray-200">
        <span class="text-5xl block mb-3 opacity-50">💭</span>
        <p class="text-gray-500 font-medium">İlk yorumu sen yap! Henüz kimse deneyimini paylaşmamış.</p>
      </div>

      <div v-else class="space-y-6">
        <div 
          v-for="comment in comments" 
          :key="comment.id" 
          class="bg-white p-6 rounded-2xl border border-gray-100 hover:shadow-md transition-shadow group"
        >
          <div class="flex justify-between items-start mb-3">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 bg-gradient-to-br from-blue-100 to-galeri-yesil/20 rounded-full flex items-center justify-center font-bold text-gray-600 border border-gray-200">
                {{ comment.userEmail.charAt(0).toUpperCase() }}
              </div>
              <div>
                <h4 class="font-bold text-gray-800">{{ maskEmail(comment.userEmail) }}</h4>
                <div class="flex items-center gap-1 mt-0.5">
                  <span v-for="s in 5" :key="s" class="text-sm" :class="s <= comment.rating ? 'text-yellow-400' : 'text-gray-200'">⭐</span>
                </div>
              </div>
            </div>
            <span class="text-xs text-gray-400 font-medium bg-gray-100 px-2 py-1 rounded-md">
              {{ new Date(comment.createdAt).toLocaleDateString('tr-TR') }}
            </span>
          </div>
          
          <p class="text-gray-600 leading-relaxed mt-4 text-[15px]">
            {{ comment.content }}
          </p>

          <div class="mt-5 flex items-center justify-end border-t border-gray-50 pt-4">
            <button 
              @click="markHelpful(comment.id)" 
              class="flex items-center gap-2 text-sm font-semibold transition-colors"
              :class="helpfulVotes[comment.id] ? 'text-galeri-yesil' : 'text-gray-400 hover:text-gray-600'"
            >
              <span class="text-lg transition-transform" :class="helpfulVotes[comment.id] ? 'scale-125' : 'group-hover:scale-110'">👍</span>
              <span>Faydalı ({{ comment.helpfulCount }})</span>
            </button>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { jwtDecode } from 'jwt-decode';
import axios from 'axios';

const props = defineProps({
  targetType: {
    type: String,
    required: true, // "Artwork" veya "Workshop"
  },
  targetId: {
    type: Number,
    required: true
  }
});

const comments = ref([]);
const loading = ref(true);
const isSubmitting = ref(false);
const hoverRating = ref(0);
const helpfulVotes = ref({}); // { commentId: true/false }

const newComment = ref({
  rating: 0,
  content: ''
});

// E-posta gizleme (örnek: ali@gmail.com -> a***@gmail.com)
const maskEmail = (email) => {
  if (!email) return 'Kullanıcı';
  const parts = email.split('@');
  if (parts.length !== 2) return email;
  return `${parts[0].charAt(0)}***@${parts[1]}`;
};

// Ortalama puan hesaplama
const averageRating = computed(() => {
  if (comments.value.length === 0) return "0.0";
  const sum = comments.value.reduce((acc, curr) => acc + curr.rating, 0);
  return (sum / comments.value.length).toFixed(1);
});

// Yorumları getir
const fetchComments = async () => {
  try {
    loading.value = true;
    const response = await axios.get(`http://localhost:8080/comments/list?type=${props.targetType}&id=${props.targetId}`);
    comments.value = response.data || [];
  } catch (error) {
    console.error("Yorumlar yüklenemedi:", error);
  } finally {
    loading.value = false;
  }
};

// Yorum gönder
const submitComment = async () => {
  const token = localStorage.getItem('userToken');
  if (!token) {
    alert("Yorum yapabilmek için giriş yapmalısınız! 👤");
    return;
  }

  if (newComment.value.rating === 0) {
    alert("Lütfen bir yıldız puanı verin! ⭐");
    return;
  }

  isSubmitting.value = true;
  try {
    const decoded = jwtDecode(token);
    await axios.post('http://localhost:8080/comments/add', {
      userEmail: decoded.email,
      targetType: props.targetType,
      targetId: props.targetId,
      content: newComment.value.content,
      rating: newComment.value.rating
    });

    // Başarılı
    newComment.value.content = '';
    newComment.value.rating = 0;
    await fetchComments(); // Listeyi yenile
    alert("Yorumunuz başarıyla eklendi! 🎉");
    
  } catch (error) {
    if (error.response && error.response.status === 403) {
      alert(`Bu ${props.targetType === 'Artwork' ? 'eseri satın almadığınız' : 'atölyeye kayıt olmadığınız'} için yorum yapamazsınız! 🚫`);
    } else {
      alert("Yorum eklenirken bir hata oluştu.");
      console.error(error);
    }
  } finally {
    isSubmitting.value = false;
  }
};

// Faydalı bul
const markHelpful = async (commentId) => {
  const token = localStorage.getItem('userToken');
  if (!token) {
    alert("Değerlendirme yapmak için giriş yapmalısınız!");
    return;
  }
  
  if (helpfulVotes.value[commentId]) {
    return; // Zaten basılmış
  }

  try {
    const decoded = jwtDecode(token);
    await axios.post('http://localhost:8080/comments/helpful', {
      userEmail: decoded.email,
      commentId: commentId
    });

    helpfulVotes.value[commentId] = true;
    
    // UI'da optimistic update yapalım (anında artsın)
    const comment = comments.value.find(c => c.id === commentId);
    if (comment) {
      comment.helpfulCount++;
    }

  } catch (error) {
    if (error.response && error.response.status === 409) {
      alert("Bu yorumu zaten faydalı buldunuz!");
      helpfulVotes.value[commentId] = true; // Lokal durumu güncelle
    } else {
      console.error("Faydalı oyu verilemedi:", error);
    }
  }
};

onMounted(() => {
  fetchComments();
});
</script>
