<template>
  <div class="mt-10">
    <h3 class="text-2xl font-black text-gray-800 mb-6">💬 Yorumlar</h3>

    <!-- Yorum Ekleme -->
    <div class="bg-white rounded-2xl shadow p-6 mb-6">
      <textarea
        v-model="newComment"
        placeholder="Yorumunuzu yazın..."
        rows="3"
        class="w-full border-2 border-gray-200 rounded-xl px-4 py-3 focus:outline-none focus:border-galeri-yesil transition-colors resize-none"
      ></textarea>
      <button
        @click="postComment"
        :disabled="!newComment.trim()"
        class="mt-3 bg-galeri-yesil text-white px-6 py-3 rounded-xl font-bold hover:opacity-90 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
      >
        Yorum Ekle
      </button>
    </div>

    <!-- Yorum Listesi -->
    <div v-if="comments.length === 0" class="text-center text-gray-400 py-8">
      <p class="text-4xl mb-3">💭</p>
      <p>Henüz yorum yok. İlk yorumu sen yap!</p>
    </div>

    <div v-else class="space-y-4">
      <div
        v-for="comment in comments"
        :key="comment.id"
        class="bg-white rounded-2xl shadow p-5"
      >
        <div class="flex justify-between items-start mb-2">
          <div>
            <span class="font-bold text-gray-800">{{ comment.userName }}</span>
            <span class="text-gray-400 text-sm ml-2">{{ formatDate(comment.createdAt) }}</span>
          </div>
        </div>
        <p class="text-gray-700 mb-3">{{ comment.content }}</p>
        <button
          @click="toggleHelpful(comment)"
          :class="helpfulVoted.has(comment.id)
            ? 'bg-green-100 text-green-700 border-green-300'
            : 'bg-gray-100 text-gray-600 border-gray-200'"
          class="flex items-center gap-2 px-4 py-2 rounded-xl text-sm font-bold border-2 hover:scale-105 transition-all cursor-pointer"
        >
          👍 Faydalı ({{ comment.helpfulCount }})
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';

const props = defineProps({
  targetType: { type: String, required: true },
  targetId: { type: [Number, String], required: true },
});

const comments = ref([]);
const newComment = ref('');
const helpfulVoted = ref(new Set());
const userEmail = localStorage.getItem('userEmail') || '';
const userName = localStorage.getItem('userFirstName') || 'Ziyaretçi';

const fetchComments = async () => {
  try {
    const res = await axios.get(`http://localhost:8080/comments/list?targetType=${props.targetType}&targetId=${props.targetId}`);
    comments.value = res.data || [];
  } catch (e) {
    console.error('Yorumlar yüklenemedi:', e);
  }
};

const postComment = async () => {
  if (!newComment.value.trim()) return;
  try {
    await axios.post('http://localhost:8080/comments/add', {
      targetType: props.targetType,
      targetId: Number(props.targetId),
      userEmail,
      userName,
      content: newComment.value.trim(),
    });
    newComment.value = '';
    await fetchComments();
  } catch (e) {
    console.error('Yorum eklenemedi:', e);
  }
};

const toggleHelpful = async (comment) => {
  try {
    await axios.post('http://localhost:8080/comments/helpful', {
      commentId: comment.id,
      userEmail,
    });
    // Toggle local state
    if (helpfulVoted.value.has(comment.id)) {
      helpfulVoted.value.delete(comment.id);
      comment.helpfulCount--;
    } else {
      helpfulVoted.value.add(comment.id);
      comment.helpfulCount++;
    }
    // Force reactivity update
    helpfulVoted.value = new Set(helpfulVoted.value);
  } catch (e) {
    console.error('Faydalı oyu verilemedi:', e);
  }
};

const formatDate = (dateStr) => {
  if (!dateStr) return '';
  const d = new Date(dateStr);
  return d.toLocaleDateString('tr-TR', { day: 'numeric', month: 'long', year: 'numeric' });
};

onMounted(fetchComments);
</script>
