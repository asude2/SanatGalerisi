<template>
  <div class="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div v-if="user" class="max-w-4xl mx-auto">
      
      <!-- Profil Kartı -->
      <div class="bg-white rounded-3xl shadow-xl overflow-hidden mb-8 border border-gray-100">
        <div class="h-32 bg-gradient-to-r from-blue-600 to-indigo-700"></div>
        <div class="px-8 pb-8">
          <div class="relative flex justify-between items-end -mt-12 mb-6">
            <div class="w-32 h-32 bg-white p-2 rounded-3xl shadow-lg">
              <div class="w-full h-full bg-blue-100 rounded-2xl flex items-center justify-center text-4xl">
                {{ user.firstName.charAt(0) }}
              </div>
            </div>
            <div class="flex gap-3">
              <button 
                @click="activeTab = 'edit'"
                class="px-6 py-2.5 bg-gray-900 text-white rounded-xl font-bold text-sm hover:bg-gray-800 transition-all"
              >
                Profili Düzenle
              </button>
              <button 
                @click="logout"
                class="px-6 py-2.5 bg-red-50 text-red-600 rounded-xl font-bold text-sm hover:bg-red-100 transition-all"
              >
                Çıkış Yap
              </button>
            </div>
          </div>
          
          <div>
            <h1 class="text-3xl font-black text-gray-900">{{ user.firstName }} {{ user.lastName }}</h1>
            <p class="text-gray-500 font-medium">{{ user.email }} • <span class="text-blue-600">{{ user.role }}</span></p>
            <p v-if="user.biography" class="mt-4 text-gray-600 leading-relaxed max-w-2xl">{{ user.biography }}</p>
          </div>

          <!-- Hızlı İşlemler (Sadece Eğitmenler İçin) -->
          <div v-if="user.role === 'Instructor'" class="mt-8 flex gap-4">
             <router-link to="/add-artwork" class="px-6 py-3 bg-green-600 text-white rounded-2xl font-bold text-sm shadow-lg shadow-green-100 hover:bg-green-700 transition-all">+ Eser Ekle</router-link>
             <router-link to="/add-workshop" class="px-6 py-3 bg-blue-600 text-white rounded-2xl font-bold text-sm shadow-lg shadow-blue-100 hover:bg-blue-700 transition-all">+ Atölye Oluştur</router-link>
          </div>
        </div>
      </div>

      <!-- Navigasyon Sekmeleri -->
      <div class="flex gap-2 mb-8 bg-gray-200/50 p-1.5 rounded-2xl w-fit overflow-x-auto">
        <button 
          v-for="tab in visibleTabs" 
          :key="tab.id"
          @click="activeTab = tab.id"
          class="px-6 py-2.5 rounded-xl font-bold text-sm transition-all whitespace-nowrap"
          :class="activeTab === tab.id ? 'bg-white text-blue-600 shadow-sm' : 'text-gray-500 hover:text-gray-700'"
        >
          {{ tab.label }}
        </button>
      </div>

      <!-- Sekme İçerikleri -->
      <div class="space-y-6">
        
        <!-- Favoriler -->
        <div v-if="activeTab === 'favorites'" class="grid grid-cols-1 sm:grid-cols-2 gap-6">
          <div v-for="artwork in favorites" :key="artwork.id" class="bg-white p-4 rounded-3xl border border-gray-100 shadow-sm group">
            <div class="aspect-square rounded-2xl overflow-hidden mb-4">
              <img :src="artwork.imageUrl" class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500" />
            </div>
            <h3 class="font-bold text-gray-900">{{ artwork.title }}</h3>
            <p class="text-sm text-gray-500 mb-4">{{ artwork.artist }}</p>
            <router-link :to="`/artwork/${artwork.id}`" class="block w-full py-3 text-center bg-gray-50 text-blue-600 rounded-xl font-bold text-sm hover:bg-blue-600 hover:text-white transition-all">İncele</router-link>
          </div>
          <div v-if="favorites.length === 0" class="col-span-full py-20 text-center bg-white rounded-3xl border-2 border-dashed border-gray-200">
             <p class="text-4xl mb-4">❤️</p>
             <p class="text-gray-400 font-bold">Henüz favori eseriniz yok.</p>
          </div>
        </div>

        <!-- Siparişler -->
        <div v-if="activeTab === 'orders'" class="space-y-4">
          <div v-for="order in purchases" :key="order.id" class="bg-white p-6 rounded-3xl border border-gray-100 shadow-sm flex items-center gap-6">
            <img :src="order.image" class="w-20 h-20 rounded-2xl object-cover" />
            <div class="flex-1">
              <h3 class="font-bold text-gray-900">{{ order.title }}</h3>
              <p class="text-sm text-gray-500">{{ formatDate(order.date) }}</p>
            </div>
            <div class="text-right">
              <p class="font-black text-gray-900">{{ order.price }} ₺</p>
              <span class="text-[10px] font-black uppercase px-2 py-1 bg-green-100 text-green-700 rounded-lg">Tamamlandı</span>
            </div>
          </div>
          <div v-if="purchases.length === 0" class="py-20 text-center bg-white rounded-3xl border-2 border-dashed border-gray-200">
             <p class="text-4xl mb-4">🛍️</p>
             <p class="text-gray-400 font-bold">Henüz bir eser satın almadınız.</p>
          </div>
        </div>

        <!-- Atölye Kayıtlarım -->
        <div v-if="activeTab === 'enrollments'" class="space-y-4">
          <div v-for="enrollment in enrollments" :key="enrollment.id" class="bg-white p-6 rounded-3xl border border-gray-100 shadow-sm flex items-center gap-6">
            <div class="w-16 h-16 bg-blue-100 rounded-2xl flex items-center justify-center text-2xl">🎨</div>
            <div class="flex-1">
              <h3 class="font-bold text-gray-900">{{ enrollment.workshopTitle }}</h3>
              <p class="text-sm text-gray-500">{{ enrollment.location }} • {{ enrollment.reservedDate }}</p>
              <p class="text-xs text-blue-600 font-bold mt-1">{{ enrollment.participantCount }} Katılımcı</p>
            </div>
            <div class="text-right">
              <span class="text-[10px] font-black uppercase px-2 py-1 bg-blue-50 text-blue-700 rounded-lg">Kayıtlı</span>
            </div>
          </div>
          <div v-if="enrollments.length === 0" class="py-20 text-center bg-white rounded-3xl border-2 border-dashed border-gray-200">
             <p class="text-4xl mb-4">📅</p>
             <p class="text-gray-400 font-bold">Henüz bir atölyeye kayıt olmadınız.</p>
          </div>
        </div>

        <!-- Atölyelerim (Sadece Eğitmenler) -->
        <div v-if="activeTab === 'my_workshops'" class="space-y-4">
          <div v-for="ws in myWorkshops" :key="ws.id" class="bg-white p-6 rounded-3xl border border-gray-100 shadow-sm flex items-center gap-6">
            <img :src="ws.image" class="w-20 h-20 rounded-2xl object-cover" />
            <div class="flex-1">
              <h3 class="font-bold text-gray-900">{{ ws.title }}</h3>
              <p class="text-sm text-gray-500">{{ ws.location }} • {{ ws.price }} ₺</p>
              <p class="text-xs text-blue-600 font-bold mt-1">Kapasite: {{ ws.capacity }} Kişi</p>
            </div>
            <div class="text-right">
              <button @click="deleteWorkshop(ws.id)" class="p-2 text-red-500 hover:bg-red-50 rounded-xl transition-colors">🗑️</button>
            </div>
          </div>
          <div v-if="myWorkshops.length === 0" class="py-20 text-center bg-white rounded-3xl border-2 border-dashed border-gray-200">
             <p class="text-4xl mb-4">🎨</p>
             <p class="text-gray-400 font-bold">Henüz bir atölye oluşturmadınız.</p>
          </div>
        </div>

        <!-- Karşılaştırma Analizlerim -->
        <div v-if="activeTab === 'comparisons'" class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div v-for="comp in comparisons" :key="comp.comparisonId" class="bg-white p-8 rounded-[32px] border border-gray-100 shadow-sm hover:shadow-md transition-all group">
            <div class="flex justify-between items-start mb-6">
              <div class="flex-1">
                <div v-if="editingId === comp.comparisonId" class="flex gap-2">
                  <input v-model="editTitle" type="text" class="flex-1 p-2 border-2 border-blue-500 rounded-xl outline-none font-bold" />
                  <button @click="updateTitle(comp.comparisonId)" class="p-2 bg-green-500 text-white rounded-xl">✓</button>
                  <button @click="editingId = null" class="p-2 bg-gray-200 text-gray-600 rounded-xl">✕</button>
                </div>
                <div v-else class="flex items-center gap-2 group/title">
                  <h3 class="text-xl font-black text-gray-900">{{ comp.title || 'İsimsiz Analiz' }}</h3>
                  <button @click="startEdit(comp)" class="opacity-0 group-hover/title:opacity-100 p-1 text-gray-400 hover:text-blue-600 transition-all">✏️</button>
                </div>
                <p class="text-xs font-bold text-gray-400 uppercase tracking-widest mt-1">{{ comp.targetType === 'Artwork' ? 'Sanat Eserleri' : 'Atölyeler' }} • {{ formatDate(comp.createdAt) }}</p>
              </div>
              <button @click="deleteComparison(comp.comparisonId)" class="p-2 text-gray-300 hover:text-red-500 transition-colors">🗑️</button>
            </div>
            <button 
              @click="viewComparison(comp)"
              class="w-full py-4 bg-blue-50 text-blue-600 rounded-2xl font-bold hover:bg-blue-600 hover:text-white transition-all flex items-center justify-center gap-2"
            >
              Analizi Görüntüle ⚖️
            </button>
          </div>
          <div v-if="comparisons.length === 0" class="col-span-full py-20 text-center bg-white rounded-[40px] border-2 border-dashed border-gray-200">
             <p class="text-6xl mb-4">📊</p>
             <p class="text-gray-400 font-bold text-xl">Henüz kayıtlı bir analiziniz yok.</p>
             <router-link to="/compare" class="text-blue-600 font-bold mt-2 hover:underline block">Hemen ilk analizini yap →</router-link>
          </div>
        </div>

        <!-- Destek Taleplerim -->
        <div v-if="activeTab === 'tickets'" class="space-y-4">
          <div v-for="ticket in tickets" :key="ticket.ticketId" class="bg-white p-6 rounded-3xl border border-gray-100 shadow-sm flex items-center gap-6">
            <div class="w-12 h-12 bg-indigo-100 rounded-2xl flex items-center justify-center text-xl">🎫</div>
            <div class="flex-1">
              <h3 class="font-bold text-gray-900">{{ ticket.subject }}</h3>
              <p class="text-xs text-gray-400">{{ formatDate(ticket.createdAt) }} • {{ ticket.supportType }}</p>
            </div>
            <div class="text-right">
              <span class="px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-[10px] font-black uppercase">{{ translateStatus(ticket.status) }}</span>
              <router-link :to="`/support/${ticket.ticketId}`" class="block text-xs font-bold text-blue-600 mt-2 hover:underline">Detaylar →</router-link>
            </div>
          </div>
          <div v-if="tickets.length === 0" class="py-20 text-center bg-white rounded-3xl border-2 border-dashed border-gray-200">
             <p class="text-4xl mb-4">🎧</p>
             <p class="text-gray-400 font-bold">Henüz bir destek talebiniz yok.</p>
             <router-link to="/support" class="text-blue-600 font-bold mt-2 hover:underline block">Yeni talep oluştur →</router-link>
          </div>
        </div>

        <!-- Profil Düzenle -->
        <div v-if="activeTab === 'edit'" class="bg-white p-8 rounded-3xl border border-gray-100 shadow-sm max-w-2xl">
          <h2 class="text-2xl font-black text-gray-900 mb-8">Profil Bilgilerini Güncelle</h2>
          <form @submit.prevent="updateProfile" class="space-y-6">
            <div class="grid grid-cols-2 gap-4">
              <div class="space-y-2">
                <label class="text-xs font-black text-gray-400 uppercase ml-1">Ad</label>
                <input v-model="editUser.firstName" type="text" class="w-full p-4 bg-gray-50 rounded-2xl border-2 border-transparent focus:border-blue-500 outline-none transition-all" />
              </div>
              <div class="space-y-2">
                <label class="text-xs font-black text-gray-400 uppercase ml-1">Soyad</label>
                <input v-model="editUser.lastName" type="text" class="w-full p-4 bg-gray-50 rounded-2xl border-2 border-transparent focus:border-blue-500 outline-none transition-all" />
              </div>
            </div>
            <div class="space-y-2">
              <label class="text-xs font-black text-gray-400 uppercase ml-1">Biyografi (Sanatçılar İçin)</label>
              <textarea v-model="editUser.biography" rows="4" class="w-full p-4 bg-gray-50 rounded-2xl border-2 border-transparent focus:border-blue-500 outline-none transition-all resize-none"></textarea>
            </div>
            <button type="submit" class="w-full py-4 bg-blue-600 text-white rounded-2xl font-bold hover:bg-blue-700 transition-all shadow-lg shadow-blue-100">Değişiklikleri Kaydet</button>
          </form>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';

const router = useRouter();
const user = ref(null);
const favorites = ref([]);
const purchases = ref([]);
const comparisons = ref([]);
const enrollments = ref([]);
const tickets = ref([]);
const myWorkshops = ref([]);
const activeTab = ref('favorites');
const editingId = ref(null);
const editTitle = ref('');

const tabs = [
  { id: 'favorites', label: 'Favorilerim' },
  { id: 'orders', label: 'Siparişlerim' },
  { id: 'enrollments', label: 'Atölye Kayıtlarım' },
  { id: 'my_workshops', label: 'Atölyelerim', instructorOnly: true },
  { id: 'comparisons', label: 'Analizlerim' },
  { id: 'tickets', label: 'Destek Taleplerim' },
  { id: 'edit', label: 'Düzenle' }
];

const visibleTabs = computed(() => {
  return tabs.filter(tab => !tab.instructorOnly || (user.value && user.value.role === 'Instructor'));
});

const editUser = ref({ firstName: '', lastName: '', biography: '' });

const fetchUser = async () => {
  const email = localStorage.getItem('userEmail');
  try {
    const res = await axios.get(`http://localhost:8080/profile?email=${email}`);
    user.value = res.data;
    editUser.value = { ...res.data };
  } catch (e) { console.error(e); }
};

const fetchFavorites = async () => {
  const email = localStorage.getItem('userEmail');
  try {
    const res = await axios.get(`http://localhost:8080/favorites/list?email=${email}`);
    favorites.value = res.data || [];
  } catch (e) { console.error(e); }
};

const fetchPurchases = async () => {
  const email = localStorage.getItem('userEmail');
  try {
    const res = await axios.get(`http://localhost:8080/user-purchases?email=${email}`);
    purchases.value = res.data || [];
  } catch (e) { console.error(e); }
};

const fetchEnrollments = async () => {
  const email = localStorage.getItem('userEmail');
  try {
    const res = await axios.get(`http://localhost:8080/user-enrollments?email=${email}`);
    enrollments.value = res.data || [];
  } catch (e) { console.error(e); }
};

const fetchTickets = async () => {
  const token = localStorage.getItem('userToken');
  try {
    const res = await axios.get('http://localhost:8080/tickets', { headers: { Authorization: `Bearer ${token}` } });
    tickets.value = res.data || [];
  } catch (e) { console.error(e); }
};

const fetchMyWorkshops = async () => {
  const email = localStorage.getItem('userEmail');
  try {
    const res = await axios.get(`http://localhost:8080/my-workshops?email=${email}`);
    myWorkshops.value = res.data || [];
  } catch (e) { console.error(e); }
};

const fetchComparisons = async () => {
  const token = localStorage.getItem('userToken');
  try {
    const res = await axios.get('http://localhost:8080/comparisons', { headers: { Authorization: `Bearer ${token}` } });
    comparisons.value = res.data || [];
  } catch (e) { console.error(e); }
};

const updateProfile = async () => {
  try {
    await axios.post('http://localhost:8080/profile/update', { ...editUser.value, email: user.value.email });
    alert("Profil güncellendi!");
    fetchUser();
  } catch (e) { alert("Hata oluştu."); }
};

const startEdit = (comp) => {
  editingId.value = comp.comparisonId;
  editTitle.value = comp.title;
};

const updateTitle = async (id) => {
  const token = localStorage.getItem('userToken');
  try {
    await axios.post('http://localhost:8080/comparisons/update', { comparisonId: id, title: editTitle.value }, { headers: { Authorization: `Bearer ${token}` } });
    editingId.value = null;
    fetchComparisons();
  } catch (e) { alert("Güncellenemedi."); }
};

const deleteComparison = async (id) => {
  if (!confirm("Bu analizi silmek istediğinize emin misiniz?")) return;
  const token = localStorage.getItem('userToken');
  try {
    await axios.delete(`http://localhost:8080/comparisons/delete?id=${id}`, { headers: { Authorization: `Bearer ${token}` } });
    fetchComparisons();
  } catch (e) { alert("Silinemedi."); }
};

const deleteWorkshop = async (id) => {
  if (!confirm("Bu atölyeyi silmek istediğinize emin misiniz?")) return;
  try {
    await axios.delete(`http://localhost:8080/delete-workshop?id=${id}`);
    fetchMyWorkshops();
  } catch (e) { alert("Silinemedi."); }
};

const viewComparison = (comp) => {
  localStorage.setItem('compareType', comp.targetType);
  const ids = comp.targetIds.split(',').map(id => parseInt(id));
  localStorage.setItem('compareList', JSON.stringify(ids));
  router.push('/compare');
};

const logout = () => {
  localStorage.clear();
  router.push('/login');
};

const translateStatus = (status) => {
  const map = { 'Open': 'Açık', 'Açık': 'Açık', 'Responded': 'Yanıtlandı', 'Beklemede': 'Yanıtlandı', 'Closed': 'Çözüldü', 'Çözüldü': 'Çözüldü' };
  return map[status] || status;
};

const formatDate = (d) => {
  if (!d) return '';
  return new Intl.DateTimeFormat('tr-TR', { day: 'numeric', month: 'long' }).format(new Date(d));
};

onMounted(() => {
  fetchUser();
  fetchFavorites();
  fetchPurchases();
  fetchEnrollments();
  fetchComparisons();
  fetchTickets();
  fetchMyWorkshops();
});
</script>
