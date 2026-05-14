<template>
  <div class="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-4xl mx-auto bg-white rounded-2xl shadow-xl overflow-hidden border border-gray-100">
      
      <div class="bg-gradient-to-r from-blue-600 to-indigo-700 px-8 py-10 text-white">

      <button @click="goBack" class="mb-5 p-2 bg-white/20 hover:bg-white/30 rounded-full transition-all group"title="Geri Git">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-white group-hover:-translate-x-1 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
        </svg>
      </button>

        <div class="flex items-center space-x-4">
          <div class="p-3 bg-white/20 rounded-full">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
          </div>
          <div>
            <h2 class="text-3xl font-extrabold tracking-tight">Profil Bilgilerim</h2>
            <p class="mt-1 text-blue-100 text-lg">Hesap detaylarını buradan görüntüleyebilir ve güncelleyebilirsin.</p>
          </div>
        </div>
      </div>

      <div class="p-10">
        <div v-if="user" class="grid grid-cols-1 md:grid-cols-2 gap-8">
          
          <div class="space-y-2">
            <label class="block text-sm font-semibold text-gray-600 uppercase tracking-wider ml-1">Ad</label>
            <input 
              v-model="user.firstName" 
              type="text" 
              placeholder="Adınızı girin"
              class="block w-full px-5 py-4 text-lg border-2 border-gray-100 rounded-xl bg-gray-50 focus:bg-white focus:border-blue-500 transition-all outline-none"
            >
          </div>

          <div class="space-y-2">
            <label class="block text-sm font-semibold text-gray-600 uppercase tracking-wider ml-1">Soyad</label>
            <input 
              v-model="user.lastName" 
              type="text" 
              placeholder="Soyadınızı girin"
              class="block w-full px-5 py-4 text-lg border-2 border-gray-100 rounded-xl bg-gray-50 focus:bg-white focus:border-blue-500 transition-all outline-none"
            >
          </div>

          <div class="md:col-span-2 space-y-2">
            <label class="block text-sm font-semibold text-gray-600 uppercase tracking-wider ml-1">E-posta Adresi</label>
            <div class="relative">
              <input 
                v-model="user.email" 
                type="email" 
                class="block w-full px-5 py-4 text-lg border-2 border-gray-100 rounded-xl bg-gray-100 text-gray-500 cursor-not-allowed outline-none" 
                readonly
              >
              <div class="absolute right-4 top-1/2 -translate-y-1/2">
                <span class="text-xs font-medium text-gray-400 bg-gray-200 px-2 py-1 rounded">Değiştirilemez</span>
              </div>
            </div>
          </div>

          <div v-if="userRole === 'Instructor'" class="md:col-span-2 space-y-2">
            <label class="block text-sm font-semibold text-gray-600 uppercase tracking-wider ml-1">✨ Biyografi (Sanatçı Bilgisi)</label>
            <textarea 
              v-model="user.biography" 
              rows="5"
              placeholder="Kendinizi tanıtın, sanat stiliniz, ilham kaynaklarınız hakkında yazın..."
              class="block w-full px-5 py-4 text-lg border-2 border-gray-100 rounded-xl bg-gray-50 focus:bg-white focus:border-blue-500 transition-all outline-none resize-none"
            ></textarea>
            <p class="text-xs text-gray-400 ml-1">Bu bilgi ziyaretçiler tarafından görülecektir.</p>
          </div>

          <div class="md:col-span-2 pt-6">
            <button 
              @click="updateProfile" 
              class="w-full md:w-max px-10 py-4 bg-blue-600 text-white text-xl font-bold rounded-xl hover:bg-blue-700 transition-all flex items-center justify-center space-x-2 shadow-blue-200"
            >
              <span>Değişiklikleri Kaydet</span>
            </button>
          </div>

          <div class="md:col-span-2 mt-12 pt-8 border-t-2 border-gray-100">
            <h3 class="text-2xl font-bold text-gray-800 mb-6 flex items-center">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 mr-2 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
              </svg>
              Şifre Değiştir
            </h3>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div class="space-y-2">
                <label class="block text-sm font-semibold text-gray-600 uppercase tracking-wider ml-1">Eski Şifre</label>
                <input v-model="passwords.oldPassword" type="password" placeholder="••••••••" class="block w-full px-5 py-4 border-2 border-gray-100 rounded-xl bg-gray-50 focus:border-blue-500 outline-none">
              </div>

              <div class="space-y-2">
                <label class="block text-sm font-semibold text-gray-600 uppercase tracking-wider ml-1">Yeni Şifre</label>
                <input v-model="passwords.newPassword" type="password" placeholder="••••••••" class="block w-full px-5 py-4 border-2 border-gray-100 rounded-xl bg-gray-50 focus:border-blue-500 outline-none">
              </div>

              <div class="md:col-span-2">
                <button @click="changePassword" class="w-full md:w-max px-10 py-4 bg-gray-800 text-white text-xl font-bold rounded-xl hover:bg-black transition-all">
                  Şifreyi Güncelle
                </button>
              </div>
            </div>
          </div>













              <div class="mt-12 border-t-2 border-gray-100 pt-8">
                <h3 class="text-2xl font-bold text-gray-800 mb-6 flex items-center">
                  🛍️ Satın Aldığım Eserler
                </h3>

                <div v-if="boughtArtworks.length > 0" class="space-y-4">
                  <div v-for="item in boughtArtworks" :key="item.id" 
                      class="bg-white rounded-2xl p-5 border border-gray-100 shadow-sm flex items-center justify-between">
                    
                    <div class="flex items-center space-x-4">
                      <img :src="item.image" class="w-20 h-20 object-cover rounded-xl shadow-sm">
                      <div>
                        <h4 class="text-lg font-bold text-gray-900">{{ item.title }}</h4>
                        <p class="text-blue-600 font-semibold">{{ item.price }} TL</p>
                        <p class="text-xs text-gray-400">{{ new Date(item.date).toLocaleDateString('tr-TR') }}</p>
                      </div>
                    </div>

                    <div class="text-right">
                      <span 
                        :class="{
                          'bg-amber-100 text-amber-700': item.status === 'Hazırlanıyor',
                          'bg-blue-100 text-blue-700': item.status === 'Kargoda',
                          'bg-green-100 text-green-700': item.status === 'Teslim Edildi'
                        }"
                        class="px-4 py-1.5 rounded-full text-xs font-black uppercase tracking-widest"
                      >
                        {{ item.status }}
                      </span>
                    </div>
                  </div>
                </div>

                <div v-else class="text-center py-10 bg-gray-50 rounded-2xl border-2 border-dashed border-gray-200">
                  <p class="text-gray-400 italic">Henüz bir eser satın almadınız.</p>
                </div>
              </div>












          <div class="md:col-span-2 mt-12 pt-8 border-t-2 border-gray-100">
            <h3 class="text-2xl font-bold text-gray-800 mb-8 flex items-center">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 mr-2 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              Atölye Rezervasyonlarım
            </h3>

            <div v-if="enrollments && enrollments.length > 0" class="space-y-4">
              <div v-for="enroll in enrollments" :key="enroll.id" class="bg-gray-50 rounded-2xl p-6 border border-gray-100 flex flex-col md:flex-row md:items-center justify-between gap-4">
                <div>
                  <h4 class="text-xl font-bold text-gray-800">{{ enroll.workshopTitle }}</h4>
                  <div class="flex flex-wrap gap-4 mt-2 text-sm text-gray-600">
                    <span class="flex items-center">📅 {{ enroll.reservedDate }}</span>
                    <span class="flex items-center">👥 {{ enroll.participantCount }} Kişi</span>
                    <span class="flex items-center">📍 {{ enroll.location }}</span>
                  </div>
                </div>
                
                <div class="flex gap-2">
                  <button @click="openEditModal(enroll)" class="px-4 py-2 bg-blue-100 text-blue-700 font-bold rounded-lg hover:bg-blue-200 transition-colors">
                    Düzenle
                  </button>
                  <button @click="cancelEnrollment(enroll.id)" class="px-4 py-2 bg-red-100 text-red-700 font-bold rounded-lg hover:bg-red-200 transition-colors">
                    İptal Et
                  </button>
                </div>
              </div>
            </div>

            <div v-else class="text-center py-10 bg-gray-50 rounded-2xl border-2 border-dashed border-gray-200">
              <p class="text-gray-400 italic">Henüz bir atölye rezervasyonunuz bulunmuyor.</p>
            </div>
          </div>

          <div class="md:col-span-2 mt-12 pt-8 border-t-2 border-gray-100">
            <h3 class="text-2xl font-bold text-gray-800 mb-8 flex items-center">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 mr-2 text-red-500" fill="currentColor" viewBox="0 0 24 24">
                <path d="M11.645 20.91l-.007-.003-.022-.012a15.247 15.247 0 01-.383-.218 25.18 25.18 0 01-4.244-3.17C4.688 15.36 2.25 12.174 2.25 8.25 2.25 5.322 4.714 3 7.688 3c1.74 0 3.285.797 4.312 2.022C13.027 3.797 14.572 3 16.312 3c2.974 0 5.438 2.322 5.438 5.25 0 3.924-2.438 7.111-4.739 9.256a25.175 25.175 0 01-4.244 3.17 15.247 15.247 0 01-.383.219l-.022.012-.007.004-.003.001z" />
              </svg>
              Favori Eserlerim
            </h3>

            <div v-if="favorites && favorites.length > 0" class="grid grid-cols-1 sm:grid-cols-3 gap-6">
              <ArtworkCard 
                v-for="fav in favorites" 
                :key="fav.id"
                :id="fav.id"
                :title="fav.title"
                :artist="fav.artist"
                :price="fav.price"
                :image="fav.imageUrl"
              />
            </div>
            <div v-else class="text-center py-10 bg-gray-50 rounded-2xl border-2 border-dashed border-gray-200">
              <p class="text-gray-400 italic">Henüz bir eseri favorilere eklemediniz.</p>
            </div>
          </div>

          <div v-if="userRole === 'Instructor'"  class="md:col-span-2 mt-12 pt-8 border-t-2 border-gray-100">
            <h3 class="text-2xl font-bold text-gray-800 mb-8 flex items-center">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 mr-2 text-green-600" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z" />
              </svg>
              Eklediğim Eserler
            </h3>

            <div v-if="myArtworks && myArtworks.length > 0" class="space-y-4">
              <div v-for="artwork in myArtworks" :key="artwork.id" 
                  class="bg-white rounded-2xl p-5 border border-gray-100 shadow-sm flex items-center justify-between hover:shadow-md transition-shadow">
                
                <div class="flex items-center space-x-4 flex-1">
                  <img :src="artwork.imageUrl" class="w-20 h-20 object-cover rounded-xl shadow-sm" :alt="artwork.title">
                  <div>
                    <h4 class="text-lg font-bold text-gray-900">{{ artwork.title }}</h4>
                    <p class="text-gray-600 text-sm">{{ artwork.description?.substring(0, 50) }}...</p>
                    <p class="text-green-600 font-semibold mt-1">{{ artwork.price }} TL</p>
                  </div>
                </div>

                <button 
                  @click="deleteMyArtwork(artwork.id)" 
                  class="ml-4 px-4 py-2 bg-red-100 text-red-600 font-bold rounded-lg hover:bg-red-200 transition-colors flex items-center gap-2"
                >
                  🗑️ Sil
                </button>
              </div>
            </div>

            <div v-else class="text-center py-10 bg-gray-50 rounded-2xl border-2 border-dashed border-gray-200">
              <p class="text-gray-400 italic">Henüz eser eklemediniz.</p>
            </div>
          </div>
        </div>

        <div v-else class="flex flex-col items-center justify-center py-20 space-y-4">
          <div class="animate-spin rounded-full h-12 w-12 border-b-4 border-blue-600"></div>
          <p class="text-xl text-gray-500 font-medium">Bilgileriniz getiriliyor...</p>
        </div>


        <div v-if="userRole === 'Instructor'" class="md:col-span-2 mt-12 pt-8 border-t-2 border-gray-100">
          <h3 class="text-2xl font-bold text-gray-800 mb-8 flex items-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 mr-2 text-indigo-600" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 3L1 9l11 6 9-4.91V17h2V9L12 3z" />
            </svg>
            Oluşturduğum Atölyeler
          </h3>

          <div v-if="myWorkshops && myWorkshops.length > 0" class="space-y-4">
            <div v-for="ws in myWorkshops" :key="ws.id" 
                class="bg-white rounded-2xl p-5 border border-gray-100 shadow-sm flex items-center justify-between hover:shadow-md transition-shadow">
              
              <div class="flex items-center space-x-4 flex-1">
                <img :src="ws.imageUrl" class="w-20 h-20 object-cover rounded-xl shadow-sm" :alt="ws.title">
                <div>
                  <h4 class="text-lg font-bold text-gray-900">{{ ws.title }}</h4>
                  <p class="text-gray-500 text-sm">📍 {{ ws.location }} | 📅 {{ ws.availableDates }}</p>
                  <p class="text-indigo-600 font-semibold mt-1">{{ ws.price }} TL <span class="text-gray-400 text-xs">(Kapasite: {{ ws.capacity }})</span></p>
                </div>
              </div>

              <div class="flex gap-2">
                <button 
                  @click="deleteMyWorkshop(ws.id)" 
                  class="px-4 py-2 bg-red-50 text-red-600 font-bold rounded-lg hover:bg-red-100 transition-colors"
                >
                  🗑️ Sil
                </button>
              </div>
            </div>
          </div>

          <div v-else class="text-center py-10 bg-gray-50 rounded-2xl border-2 border-dashed border-gray-200">
            <p class="text-gray-400 italic">Henüz bir atölye oluşturmadınız.</p>
          </div>
        </div>
      </div>
    </div>



    <div v-if="showEditModal && selectedEnrollment" class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-3xl p-8 max-w-md w-full shadow-2xl">
        <h3 class="text-2xl font-bold text-gray-800 mb-6">Rezervasyonu Güncelle</h3>
        <div class="space-y-6">
          <div>
            <label class="block text-sm font-bold text-gray-600 mb-2 uppercase ml-1">Katılımcı Sayısı</label>
            <select v-model.number="selectedEnrollment.participantCount" class="w-full p-4 bg-gray-50 border-2 border-gray-100 rounded-xl outline-none focus:border-blue-500">
              <option v-for="n in 10" :key="n" :value="n">{{ n }} Kişi</option>
            </select>
          </div>
          <div>
          <label class="block text-sm font-bold text-gray-600 mb-2 uppercase ml-1">Yeni Tarih & Saat</label>
          <select 
            v-model="selectedEnrollment.reservedDate" 
            class="w-full p-4 bg-gray-50 border-2 border-gray-100 rounded-xl outline-none focus:border-blue-500"
          >
            <option :value="selectedEnrollment.reservedDate">{{ selectedEnrollment.reservedDate }} (Mevcut)</option>
            
            <option 
              v-for="date in selectedEnrollment.availableDates?.split(',')" 
              :key="date" 
              :value="date.trim()"
            >
              {{ date.trim() }}
            </option>
          </select>
        </div>
          <div class="flex justify-end space-x-3 mt-8">
            <button @click="showEditModal = false" class="px-6 py-3 font-bold text-gray-500 hover:text-gray-700">Vazgeç</button>
            <button @click="updateEnrollment" class="px-8 py-3 bg-blue-600 text-white font-bold rounded-xl hover:shadow-lg transition-all">Güncelle</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import ArtworkCard from '../components/ArtworkCard.vue';

import { useRouter } from 'vue-router';
const router = useRouter();

const goBack = () => {
  router.back(); 
};

const user = ref(null);
const userRole = ref('');
const favorites = ref([]);
const myArtworks = ref([]);
const myWorkshops = ref([]);
const enrollments = ref([]);
const showEditModal = ref(false);
const selectedEnrollment = ref(null);
const boughtArtworks = ref([]);

const passwords = ref({
  oldPassword: '',
  newPassword: ''
});

// Sayfa yüklendiğinde bilgileri getir
onMounted(async () => {
  const userEmail = localStorage.getItem('userEmail');
  userRole.value = localStorage.getItem('userRole') || '';
  if (!userEmail) return;

  try {
    const profileRes = await axios.get(`http://localhost:8080/profile?email=${userEmail}`);
    user.value = profileRes.data;

    const favRes = await axios.get(`http://localhost:8080/favorites/list?email=${userEmail}`);
    favorites.value = favRes.data || [];

    await fetchUserArtworks(userEmail);
    await fetchEnrollments();
    await fetchArtworkPurchases();
  } catch (error) {
    console.error("Veriler yüklenemedi:", error);
  }

});

// Profil Güncelleme
const updateProfile = async () => {
  try {
    await axios.put('http://localhost:8080/profile/update', user.value);
    alert("Profil bilgileriniz başarıyla güncellendi.");
  } catch (error) {
    alert("Güncelleme başarısız.");
  }
};

// Şifre değiştirme fonksiyonu
const changePassword = async () => {
  if (!passwords.value.oldPassword || !passwords.value.newPassword) {
    alert("Lütfen her iki şifre alanını da doldurun.");
    return;
  }

  try {
    const response = await axios.post('http://localhost:8080/profile/change-password', {
      email: user.value.email, 
      oldPassword: passwords.value.oldPassword,
      newPassword: passwords.value.newPassword
    });

    alert(response.data.message);
    passwords.value.oldPassword = '';
    passwords.value.newPassword = '';
  } catch (error) {
    alert(error.response?.data || "Şifre değiştirilemedi.");
  }
};

// Rezervasyonları getir
const fetchEnrollments = async () => {
  const userEmail = localStorage.getItem('userEmail');
  try {
    const res = await axios.get(`http://localhost:8080/user-enrollments?email=${userEmail}`);
    enrollments.value = res.data || [];
  } catch (error) {
    console.error("Rezervasyonlar yüklenemedi:", error);
  }
};

// İptal Etme
const cancelEnrollment = async (id) => {
  if (confirm("Bu rezervasyonu iptal etmek istediğinize emin misiniz?")) {
    try {
      await axios.delete(`http://localhost:8080/delete-enrollment?id=${id}`);
      alert("Rezervasyon başarıyla iptal edildi.");
      fetchEnrollments();
    } catch (error) {
      alert("İptal işlemi sırasında bir hata oluştu.");
    }
  }
};

// Düzenleme Modalını Aç
const openEditModal = (enroll) => {
  selectedEnrollment.value = { ...enroll };
  showEditModal.value = true;
};

// Güncelleme 
const updateEnrollment = async () => {
  try {
    await axios.put(`http://localhost:8080/update-enrollment`, {
      id: selectedEnrollment.value.id,
      participantCount: parseInt(selectedEnrollment.value.participantCount),
      reservedDate: selectedEnrollment.value.reservedDate
    });
    alert("Rezervasyon güncellendi.");
    showEditModal.value = false;
    fetchEnrollments();
  } catch (error) {
    alert("Güncelleme başarısız.");
  }
};

// Satın alınan eserleri getir
const fetchArtworkPurchases = async () => {
  const email = localStorage.getItem('userEmail');
  try {
    const res = await axios.get(`http://localhost:8080/user-purchases?email=${email}`);
    boughtArtworks.value = res.data || [];
  } catch (error) {
    console.error("Siparişler yüklenemedi:", error);
  }
};

// Kullanıcının eklediği eserleri getir
const fetchUserArtworks = async (userEmail) => {
  try {
    const res = await axios.get(`http://localhost:8080/artworks`);
    const allArtworks = res.data || [];
    
    // Kullanıcının artist kaydını al
    const profileRes = await axios.get(`http://localhost:8080/profile?email=${userEmail}`);
    // Backend'den artist bilgisi döneceğini varsayarsak, şimdilik isim üzerinden filtreleyeceğiz
    const userName = profileRes.data.firstName + ' ' + profileRes.data.lastName;
    
    // Sanatçı adı eşleşen eserleri filtrele
    myArtworks.value = allArtworks.filter(art => art.artist === userName);
  } catch (error) {
    console.error("Kendi eserler yüklenemedi:", error);
  }
};

// Eser sil
const deleteMyArtwork = async (artworkId) => {
  if (confirm("Bu eseri silmek istediğinize emin misiniz? Bu işlem geri alınamaz.")) {
    try {
      const userEmail = localStorage.getItem('userEmail');
      await axios.post(`http://localhost:8080/delete-artwork`, {
        artworkId: artworkId,
        email: userEmail
      });
      alert("Eser başarıyla silindi.");
      
      // Listeyi yenile
      const userEmail2 = localStorage.getItem('userEmail');
      await fetchUserArtworks(userEmail2);
    } catch (error) {
      alert(error.response?.data || "Eser silinirken bir hata oluştu.");
    }
  }
};


// Kendi atölyelerini getir
const fetchMyWorkshops = async () => {
  const email = localStorage.getItem('userEmail');
  try {
    const res = await axios.get(`http://localhost:8080/my-workshops?email=${email}`);
    myWorkshops.value = res.data || [];
  } catch (error) {
    console.error("Kendi atölyeleriniz yüklenemedi:", error);
  }
};



const deleteMyWorkshop = async (id) => {
  if (confirm("Bu atölyeyi silmek istediğinize emin misiniz?")) {
    try {
      const userEmail = localStorage.getItem('userEmail');
      await axios.delete(`http://localhost:8080/delete-workshop?id=${id}`, {
        data: {
          workshopId: id,
          email: userEmail
        }
      });
      
      alert("Atölye başarıyla silindi.");
      
      // Listeyi anlık güncellemek için tekrar çekiyoruz
      await fetchMyWorkshops(); 
    } catch (error) {
      console.error("Silme hatası:", error);
      alert("Silme işlemi başarısız: " + (error.response?.data || "Sunucu hatası"));
    }
  }
};



onMounted(async () => {
  await fetchArtworkPurchases();
  if (userRole.value === 'Instructor') {
    await fetchMyWorkshops();
  }
});
</script>