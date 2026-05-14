<template>
  <div class="mt-16 bg-white rounded-3xl shadow-xl border border-gray-100 overflow-hidden">

    <!-- Başlık + Ortalama Puan -->
    <div class="bg-gray-50 border-b border-gray-100 p-8 flex flex-wrap gap-4 justify-between items-center">
      <div>
        <h2 class="text-2xl font-black text-gray-800">Değerlendirmeler ve Yorumlar</h2>
        <p class="text-gray-500 text-sm mt-1">
          Bu {{ targetType === 'Artwork' ? 'eseri' : 'atölyeyi' }} deneyimleyenlerin fikirleri
        </p>
      </div>

      <div v-if="avgData.count > 0" class="flex items-center gap-3 bg-white px-5 py-3 rounded-xl shadow-sm border border-yellow-100">
        <div class="text-center">
          <div class="text-3xl font-black text-yellow-500">{{ avgData.average.toFixed(1) }}</div>
          <div class="flex gap-0.5 justify-center mt-1">
            <svg v-for="s in 5" :key="s" class="w-4 h-4" viewBox="0 0 24 24">
              <polygon points="12,2 15.09,8.26 22,9.27 17,14.14 18.18,21.02 12,17.77 5.82,21.02 7,14.14 2,9.27 8.91,8.26"
                :fill="s <= Math.round(avgData.average) ? '#FBBF24' : 'none'"
                :stroke="s <= Math.round(avgData.average) ? '#FBBF24' : '#D1D5DB'"
                stroke-width="1.5"/>
            </svg>
          </div>
        </div>
        <div class="text-left border-l border-gray-100 pl-3">
          <div class="text-sm font-bold text-gray-700">{{ avgData.count }} Değerlendirme</div>
          <div class="text-xs text-gray-400 mt-0.5">Doğrulanmış puanlar</div>
        </div>
      </div>
    </div>

    <div class="p-8">

      <!-- Yorum Ekleme Formu veya Uyarı Bloğu -->
      <div class="mb-10">

        <!-- Giriş yapılmamış -->
        <div v-if="!isLoggedIn" class="flex items-center gap-4 bg-amber-50 border border-amber-200 rounded-2xl p-6">
          <span class="text-3xl">🔒</span>
          <div>
            <p class="font-bold text-amber-800">Yorum yapmak için giriş yapmalısınız</p>
            <p class="text-sm text-amber-600 mt-0.5">Hesabınız varsa giriş yapın, yoksa kayıt olun.</p>
          </div>
        </div>

        <!-- Giriş yapılmış ama hak yok -->
        <div v-else-if="!isEligible" class="flex items-center gap-4 bg-blue-50 border border-blue-200 rounded-2xl p-6">
          <span class="text-3xl">{{ targetType === 'Artwork' ? '🛒' : '📋' }}</span>
          <div>
            <p class="font-bold text-blue-800">
              {{ targetType === 'Artwork'
                ? 'Bu eseri satın aldıktan sonra değerlendirme yapabilirsiniz'
                : 'Bu atölyeye kayıt yaptırdıktan sonra yorum yapabilirsiniz' }}
            </p>
            <p class="text-sm text-blue-600 mt-0.5">
              {{ targetType === 'Artwork'
                ? 'Yalnızca doğrulanmış alıcılar yorum yapabilir.'
                : 'Yalnızca kayıtlı katılımcılar yorum yapabilir.' }}
            </p>
          </div>
        </div>

        <!-- Yorum Formu (hak varsa) -->
        <div v-else class="bg-blue-50/50 p-6 rounded-2xl border border-blue-100">
          <h3 class="text-lg font-bold text-blue-900 mb-4">Senin Düşüncen Nedir?</h3>

          <form @submit.prevent="submitComment" class="space-y-4">
            <!-- Yıldız Seçimi (SVG, boş→dolu) -->
            <div class="flex items-center gap-3">
              <span class="text-sm font-semibold text-gray-600">Puanın:</span>
              <div class="flex gap-1">
                <button
                  v-for="star in 5"
                  :key="star"
                  type="button"
                  @click="newComment.rating = star"
                  @mouseenter="hoverRating = star"
                  @mouseleave="hoverRating = 0"
                  class="transition-transform hover:scale-125 focus:outline-none"
                >
                  <svg class="w-8 h-8" viewBox="0 0 24 24">
                    <polygon
                      points="12,2 15.09,8.26 22,9.27 17,14.14 18.18,21.02 12,17.77 5.82,21.02 7,14.14 2,9.27 8.91,8.26"
                      :fill="(hoverRating ? star <= hoverRating : star <= newComment.rating) ? '#FBBF24' : 'none'"
                      :stroke="(hoverRating ? star <= hoverRating : star <= newComment.rating) ? '#FBBF24' : '#9CA3AF'"
                      stroke-width="1.5"
                      stroke-linejoin="round"
                    />
                  </svg>
                </button>
              </div>
              <span v-if="newComment.rating > 0 || hoverRating > 0"
                class="text-sm font-semibold text-blue-700 bg-blue-100 px-3 py-1 rounded-full transition-all">
                {{ ratingLabels[(hoverRating || newComment.rating) - 1] }}
              </span>
            </div>

            <!-- Yorum Alanı -->
            <textarea
              v-model="newComment.content"
              placeholder="Deneyiminizi paylaşın..."
              class="w-full p-4 bg-white rounded-xl border border-blue-200 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all min-h-[110px] resize-none"
              required
            ></textarea>

            <div class="flex justify-end">
              <button
                type="submit"
                :disabled="isSubmitting || newComment.rating === 0"
                class="bg-blue-600 text-white px-8 py-3 rounded-xl font-bold hover:bg-blue-700 transition-all disabled:opacity-50 disabled:cursor-not-allowed shadow-md shadow-blue-200 flex items-center gap-2"
              >
                <span v-if="isSubmitting" class="animate-spin">⏳</span>
                <span>{{ isSubmitting ? 'Gönderiliyor...' : '💬 Yorumu Paylaş' }}</span>
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- Sıralama Butonları -->
      <div v-if="comments.length > 0" class="flex items-center gap-3 mb-6 flex-wrap">
        <span class="text-sm font-semibold text-gray-400">Sırala:</span>
        <button
          v-for="opt in sortOptions"
          :key="opt.value"
          @click="changeSortBy(opt.value)"
          class="text-sm px-4 py-1.5 rounded-full font-semibold transition-all border"
          :class="sortBy === opt.value
            ? 'bg-gray-800 text-white border-gray-800'
            : 'bg-white text-gray-500 border-gray-200 hover:border-gray-400'"
        >{{ opt.label }}</button>
      </div>

      <!-- Yükleniyor -->
      <div v-if="loading" class="flex justify-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-4 border-gray-300"></div>
      </div>

      <!-- Boş durum -->
      <div v-else-if="comments.length === 0" class="text-center py-12 bg-gray-50 rounded-2xl border border-dashed border-gray-200">
        <span class="text-5xl block mb-3 opacity-50">💭</span>
        <p class="text-gray-500 font-medium">İlk yorumu sen yap! Henüz kimse deneyimini paylaşmamış.</p>
      </div>

      <!-- Yorum Listesi -->
      <div v-else class="space-y-5">
        <div
          v-for="comment in comments"
          :key="comment.id"
          class="bg-white p-6 rounded-2xl border border-gray-100 hover:shadow-md transition-shadow group"
        >
          <!-- Üst: Avatar + Bilgi + Tarih -->
          <div class="flex justify-between items-start mb-3">
            <div class="flex items-center gap-3">
              <div
                class="w-10 h-10 rounded-full flex items-center justify-center font-bold text-white text-lg flex-shrink-0"
                :style="{ background: avatarColor(comment.userEmail) }"
              >
                {{ comment.userEmail.charAt(0).toUpperCase() }}
              </div>
              <div>
                <div class="flex items-center gap-2">
                  <h4 class="font-bold text-gray-800 text-sm">{{ maskEmail(comment.userEmail) }}</h4>
                  
                  <!-- Sanatçı / Eğitmen Rozeti -->
                  <span v-if="comment.userRole === 'Instructor'" 
                    class="text-[10px] font-black px-2 py-0.5 rounded-md bg-amber-500 text-white shadow-sm flex items-center gap-1">
                    <span>🎨</span> SANATÇI
                  </span>

                  <!-- Doğrulanmış rozeti -->
                  <span v-if="comment.isVerified"
                    class="text-[10px] font-bold px-2 py-0.5 rounded-md flex items-center gap-1"
                    :class="targetType === 'Artwork'
                      ? 'bg-green-50 text-green-700 border border-green-100'
                      : 'bg-blue-50 text-blue-700 border border-blue-100'"
                  >
                    <span>{{ targetType === 'Artwork' ? '✓ Doğrulanmış Alıcı' : '✓ Kayıtlı Katılımcı' }}</span>
                  </span>
                </div>
                <!-- Yıldız gösterimi (SVG) -->
                <div class="flex gap-0.5 mt-1" v-if="comment.rating > 0">
                  <svg v-for="s in 5" :key="s" class="w-3.5 h-3.5" viewBox="0 0 24 24">
                    <polygon
                      points="12,2 15.09,8.26 22,9.27 17,14.14 18.18,21.02 12,17.77 5.82,21.02 7,14.14 2,9.27 8.91,8.26"
                      :fill="s <= comment.rating ? '#FBBF24' : 'none'"
                      :stroke="s <= comment.rating ? '#FBBF24' : '#D1D5DB'"
                      stroke-width="1.5"
                    />
                  </svg>
                  <span class="text-xs text-gray-400 ml-1">{{ comment.rating }}/5</span>
                </div>
              </div>
            </div>
            <span class="text-xs text-gray-400 font-medium bg-gray-100 px-2 py-1 rounded-md flex-shrink-0">
              {{ formatDate(comment.createdAt) }}
            </span>
          </div>

          <!-- Yorum İçeriği -->
          <p class="text-gray-600 leading-relaxed text-[15px] mt-3 ml-13">{{ comment.content }}</p>

          <!-- Faydalı Oy -->
          <div class="mt-4 flex items-center justify-end border-t border-gray-50 pt-4">
            <button
              @click="markHelpful(comment)"
              class="flex items-center gap-2 text-sm font-semibold transition-all px-3 py-1.5 rounded-lg"
              :class="comment.userVoted
                ? 'text-green-600 bg-green-50 hover:bg-green-100'
                : 'text-gray-400 hover:text-gray-600 hover:bg-gray-50'"
            >
              <span class="text-base transition-transform" :class="comment.userVoted ? 'scale-125' : ''">👍</span>
              <span>{{ comment.userVoted ? 'Faydalı Buldum ✓' : 'Faydalı' }} ({{ comment.helpfulCount }})</span>
            </button>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { jwtDecode } from 'jwt-decode';
import axios from 'axios';

const props = defineProps({
  targetType: { type: String, required: true },
  targetId:   { type: Number, required: true }
});

const comments     = ref([]);
const loading      = ref(true);
const isSubmitting = ref(false);
const hoverRating  = ref(0);
const sortBy       = ref('newest');
const avgData      = ref({ average: 0, count: 0 });
const isEligible   = ref(false); // satın almış / kayıt olmuş mu?

const newComment = ref({ rating: 0, content: '' });

const ratingLabels = ['Berbat 😞', 'Kötü 😕', 'Orta 😐', 'İyi 😊', 'Mükemmel 🤩'];
const sortOptions  = [
  { value: 'newest',  label: '🕐 En Yeni'       },
  { value: 'helpful', label: '👍 En Faydalı'    },
  { value: 'rating',  label: '⭐ En Yüksek Puan' },
];

// Giriş yapılmış mı?
const isLoggedIn = computed(() => !!localStorage.getItem('userToken'));

const userEmail = computed(() => localStorage.getItem('userEmail') || '');

// Tarih formatlama
const formatDate = (str) => {
  if (!str) return '';
  return new Date(str).toLocaleDateString('tr-TR', { day: '2-digit', month: 'long', year: 'numeric' });
};

// E-posta gizleme
const maskEmail = (email) => {
  if (!email) return 'Kullanıcı';
  const [local, domain] = email.split('@');
  return `${local.charAt(0)}***@${domain}`;
};

// Avatar renk (deterministik)
const avatarColor = (email) => {
  const colors = ['#6366f1','#f59e0b','#10b981','#ef4444','#8b5cf6','#06b6d4','#f97316'];
  return colors[email.charCodeAt(0) % colors.length];
};

// Kullanıcının yorum yapma hakkı var mı? (satın alma / katılım)
const checkEligibility = async () => {
  if (!isLoggedIn.value) { isEligible.value = false; return; }
  try {
    if (props.targetType === 'Artwork') {
      const res = await axios.get('http://localhost:8080/user-purchases', {
        params: { email: userEmail.value }
      });
      const purchases = res.data || [];
      isEligible.value = purchases.some(p => p.artworkId === props.targetId || p.id === props.targetId);
    } else if (props.targetType === 'Workshop') {
      const res = await axios.get('http://localhost:8080/user-enrollments', {
        params: { email: userEmail.value }
      });
      const enrollments = res.data || [];
      isEligible.value = enrollments.some(e => e.workshopId === props.targetId || e.id === props.targetId);
    }
  } catch {
    isEligible.value = false;
  }
};

// Yorumları getir
const fetchComments = async () => {
  try {
    loading.value = true;
    const res = await axios.get('http://localhost:8080/comments/list', {
      params: { targetType: props.targetType, targetId: props.targetId, sortBy: sortBy.value, userEmail: userEmail.value }
    });
    comments.value = res.data || [];
  } catch (e) {
    console.error('Yorumlar yüklenemedi:', e);
  } finally {
    loading.value = false;
  }
};

// Ortalama puanı getir
const fetchAverage = async () => {
  try {
    const res = await axios.get('http://localhost:8080/comments/average', {
      params: { targetType: props.targetType, targetId: props.targetId }
    });
    avgData.value = res.data;
  } catch { /* sessizce geç */ }
};

// Sıralama değiştir
const changeSortBy = async (val) => {
  sortBy.value = val;
  await fetchComments();
};

// Yorum gönder
const submitComment = async () => {
  if (!isLoggedIn.value) { alert('Giriş yapmalısınız!'); return; }
  if (newComment.value.rating === 0) { alert('Lütfen bir puan verin! ⭐'); return; }

  isSubmitting.value = true;
  try {
    const token = localStorage.getItem('userToken');
    const decoded = jwtDecode(token);
    await axios.post('http://localhost:8080/comments/add', {
      userEmail:  decoded.email,
      userName:   `${decoded.firstName || ''}`.trim(),
      targetType: props.targetType,
      targetId:   props.targetId,
      content:    newComment.value.content,
      rating:     newComment.value.rating
    }, {
      headers: { Authorization: `Bearer ${token}` }
    });
    newComment.value = { rating: 0, content: '' };
    hoverRating.value = 0;
    await fetchComments();
    await fetchAverage();
  } catch (e) {
    if (e.response?.status === 403) {
      alert(`🚫 ${e.response.data}`);
    } else if (e.response?.status === 401) {
      alert('Yorum yapabilmek için giriş yapmalısınız! 🔒');
    } else {
      alert('Yorum eklenirken bir hata oluştu.');
    }
    console.error(e);
  } finally {
    isSubmitting.value = false;
  }
};

// Faydalı oy toggle
const markHelpful = async (comment) => {
  if (!isLoggedIn.value) { alert('Değerlendirme yapmak için giriş yapmalısınız!'); return; }
  try {
    const res = await axios.post('http://localhost:8080/comments/helpful', {
      commentId: comment.id,
      userEmail: userEmail.value
    });
    comment.userVoted   = res.data.voted;
    comment.helpfulCount = res.data.voted
      ? comment.helpfulCount + 1
      : Math.max(0, comment.helpfulCount - 1);
  } catch (e) {
    console.error('Faydalı oy verilemedi:', e);
  }
};

onMounted(async () => {
  await Promise.all([fetchComments(), fetchAverage(), checkEligibility()]);
});
</script>
