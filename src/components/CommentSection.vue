<template>
  <div class="mt-12 bg-white p-8 rounded-3xl shadow-sm border border-gray-100">
    <div class="flex justify-between items-center mb-8">
      <h3 class="text-2xl font-bold text-gray-800">Yorumlar & Değerlendirmeler</h3>
      
      <div class="flex gap-4">
        <select v-model="sortBy" @change="fetchComments" class="p-2 bg-gray-50 border border-gray-200 rounded-xl outline-none focus:border-galeri-yesil">
          <option value="newest">En Yeni</option>
          <option value="highest">En Yüksek Puan</option>
          <option value="most_helpful">En Faydalı</option>
        </select>
      </div>
    </div>

    <!-- Yorum Ekleme Formu -->
    <div v-if="isLoggedIn" class="mb-10 bg-gray-50 p-6 rounded-2xl">
      <h4 class="font-bold text-gray-700 mb-4">Değerlendirme Yazın</h4>
      <div class="flex items-center gap-2 mb-4">
        <span class="text-gray-600 font-medium">Puanınız:</span>
        <div class="flex gap-1">
          <button 
            v-for="star in 5" :key="star" 
            @click="newComment.rating = star"
            class="text-2xl hover:scale-110 transition-transform cursor-pointer"
            :class="star <= newComment.rating ? 'text-yellow-400' : 'text-gray-300'"
          >
            ★
          </button>
        </div>
      </div>
      <textarea 
        v-model="newComment.commentText" 
        rows="3" 
        class="w-full p-4 rounded-xl border border-gray-200 focus:border-galeri-yesil outline-none mb-4"
        placeholder="Deneyiminizi paylaşın..."
      ></textarea>
      <button 
        @click="submitComment" 
        class="bg-galeri-yesil text-white px-6 py-2 rounded-xl font-bold hover:bg-green-700 transition-colors"
        :disabled="!newComment.commentText || newComment.rating === 0"
      >
        Gönder
      </button>
      <p v-if="errorMessage" class="text-red-500 mt-2">{{ errorMessage }}</p>
    </div>
    <div v-else class="mb-10 p-6 bg-yellow-50 rounded-2xl border border-yellow-100 text-yellow-800 flex items-center justify-between">
      <span>Yorum yapabilmek için giriş yapmalısınız.</span>
      <router-link to="/login" class="bg-yellow-500 text-white px-4 py-2 rounded-xl font-bold hover:bg-yellow-600">Giriş Yap</router-link>
    </div>

    <!-- Yorum Listesi -->
    <div class="space-y-6">
      <div v-if="comments.length === 0" class="text-center text-gray-400 py-8">
        Henüz yorum yapılmamış. İlk yorumu siz yapın!
      </div>
      
      <div v-for="comment in comments" :key="comment.commentId" class="border border-gray-100 p-6 rounded-2xl">
        <div class="flex justify-between items-start mb-2">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-gray-200 rounded-full flex items-center justify-center font-bold text-gray-500">
              {{ comment.userName.charAt(0) }}
            </div>
            <div>
              <div class="flex items-center gap-2">
                <span class="font-bold text-gray-800">{{ comment.userName }}</span>
                <span v-if="comment.isVerified" class="text-xs bg-green-100 text-green-700 px-2 py-1 rounded-full font-bold">
                  ✓ Doğrulanmış {{ targetType === 'Artwork' ? 'Alıcı' : 'Katılımcı' }}
                </span>
              </div>
              <div class="text-yellow-400 text-sm">
                {{ '★'.repeat(comment.rating) }}{{ '☆'.repeat(5 - comment.rating) }}
              </div>
            </div>
          </div>
          <span class="text-gray-400 text-sm">{{ new Date(comment.createdAt).toLocaleDateString() }}</span>
        </div>
        
        <p class="text-gray-700 mt-3">{{ comment.commentText }}</p>
        
        <div class="mt-4 flex items-center gap-4">
          <div class="flex items-center bg-gray-50 rounded-lg p-1 border border-gray-100">
            <button 
              @click="voteComment(comment.commentId, 'Up')" 
              class="text-xs flex items-center gap-1 px-2 py-1 text-gray-500 hover:text-blue-600 hover:bg-white rounded-md transition-all"
              title="Faydalı"
            >
              👍 {{ comment.upvotes }}
            </button>
            <div class="w-px h-4 bg-gray-200 mx-1"></div>
            <button 
              @click="voteComment(comment.commentId, 'Down')" 
              class="text-xs flex items-center gap-1 px-2 py-1 text-gray-500 hover:text-red-600 hover:bg-white rounded-md transition-all"
              title="Faydalı Değil"
            >
              👎 {{ comment.downvotes }}
            </button>
          </div>
          
          <button 
            v-if="isLoggedIn && !comment.adminReply" 
            @click="replyingTo = comment.commentId"
            class="text-sm text-galeri-yesil font-medium hover:underline"
          >
            Yanıtla
          </button>
        </div>

        <!-- Yanıt -->
        <div v-if="comment.adminReply" class="mt-4 ml-8 p-4 rounded-xl border" :class="comment.replierRole === 'User' ? 'bg-gray-50 border-gray-200' : 'bg-blue-50 border-blue-100'">
          <div class="flex items-center gap-2 mb-1">
            <span class="text-xl" v-if="comment.replierRole !== 'User'">🛡️</span>
            <span class="text-xl" v-else>👤</span>
            <span class="font-bold" :class="comment.replierRole === 'User' ? 'text-gray-800' : 'text-blue-800'">
              {{ comment.replierName || 'Yönetici' }} {{ comment.replierRole === 'User' ? '(Doğrulanmış Alıcı)' : '(Yönetici)' }}
            </span>
          </div>
          <p :class="comment.replierRole === 'User' ? 'text-gray-700' : 'text-blue-900'">{{ comment.adminReply }}</p>
        </div>

        <!-- Yanıt Formu -->
        <div v-if="replyingTo === comment.commentId" class="mt-4 ml-8">
          <textarea 
            v-model="replyText" 
            rows="2" 
            class="w-full p-3 rounded-xl border border-gray-200 focus:border-blue-400 outline-none mb-2"
            placeholder="Yanıtınız..."
          ></textarea>
          <div class="flex gap-2">
            <button @click="submitReply(comment.commentId)" class="bg-blue-600 text-white px-4 py-1 rounded-lg font-bold text-sm hover:bg-blue-700">Gönder</button>
            <button @click="replyingTo = null" class="bg-gray-200 text-gray-700 px-4 py-1 rounded-lg font-bold text-sm hover:bg-gray-300">İptal</button>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const props = defineProps({
  targetId: { type: Number, required: true },
  targetType: { type: String, required: true } // 'Artwork' veya 'Workshop'
})

const comments = ref([])
const sortBy = ref('newest')
const isLoggedIn = ref(false)
const userRole = ref('')
const errorMessage = ref('')

const newComment = ref({
  targetId: props.targetId,
  targetType: props.targetType,
  commentText: '',
  rating: 0
})

const replyingTo = ref(null)
const replyText = ref('')

const fetchComments = async () => {
  try {
    const res = await axios.get(`http://localhost:8080/comments?targetId=${props.targetId}&targetType=${props.targetType}&sort=${sortBy.value}`)
    comments.value = res.data || []
  } catch (error) {
    console.error("Yorumlar yüklenemedi", error)
  }
}

const submitComment = async () => {
  try {
    errorMessage.value = ''
    const token = localStorage.getItem('userToken')
    await axios.post('http://localhost:8080/comments/add', newComment.value, {
      headers: { Authorization: `Bearer ${token}` }
    })
    newComment.value.commentText = ''
    newComment.value.rating = 0
    fetchComments()
  } catch (error) {
    errorMessage.value = "Yorum eklenirken bir hata oluştu."
  }
}

const voteComment = async (commentId, voteType) => {
  if (!isLoggedIn.value) return alert("Oy vermek için giriş yapmalısınız.")
  try {
    const token = localStorage.getItem('userToken')
    await axios.post('http://localhost:8080/comments/vote', { commentId, voteType }, {
      headers: { Authorization: `Bearer ${token}` }
    })
    fetchComments()
  } catch (error) {
    console.error(error)
  }
}

const submitReply = async (commentId) => {
  try {
    const token = localStorage.getItem('userToken')
    await axios.post('http://localhost:8080/comments/reply', { commentId, replyText: replyText.value }, {
      headers: { Authorization: `Bearer ${token}` }
    })
    replyText.value = ''
    replyingTo.value = null
    fetchComments()
  } catch (error) {
    alert(error.response?.data || "Yanıt gönderilemedi.")
  }
}

onMounted(() => {
  isLoggedIn.value = !!localStorage.getItem('userToken')
  userRole.value = localStorage.getItem('userRole')
  fetchComments()
})
</script>