<template>
  <div class="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <!-- Ana Kart Konteynırı-->
    <div class="max-w-4xl mx-auto bg-white rounded-2xl shadow-xl overflow-hidden border border-gray-100">
      
      <!-- Üst Başlık Bölümü-->
      <div class="bg-gradient-to-r from-blue-600 to-indigo-700 px-8 py-10 text-white">
        <div class="flex items-center space-x-4">
          <div class="p-3 bg-white/20 rounded-full">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
          </div>
          <div>
            <h2 class="text-3xl font-extrabold tracking-tight">Profil Bilgilerim</h2>
            <p class="mt-1 text-blue-100 text-lg">Hesap detaylarını buradan görüntüleyebilir ve güncelleyebilirsin.</p>
          </div>
        </div>
      </div>

      <!-- Form Alanı -->
      <div class="p-10">
        <div v-if="user" class="grid grid-cols-1 md:grid-cols-2 gap-8">
          
          <!-- Ad Alanı -->
          <div class="space-y-2">
            <label class="block text-sm font-semibold text-gray-600 uppercase tracking-wider ml-1">Ad</label>
            <input 
              v-model="user.firstName" 
              type="text" 
              placeholder="Adınızı girin"
              class="block w-full px-5 py-4 text-lg border-2 border-gray-100 rounded-xl bg-gray-50 focus:bg-white focus:border-blue-500 focus:ring-4 focus:ring-blue-100 transition-all duration-200 outline-none"
            >
          </div>

          <!-- Soyad Alanı -->
          <div class="space-y-2">
            <label class="block text-sm font-semibold text-gray-600 uppercase tracking-wider ml-1">Soyad</label>
            <input 
              v-model="user.lastName" 
              type="text" 
              placeholder="Soyadınızı girin"
              class="block w-full px-5 py-4 text-lg border-2 border-gray-100 rounded-xl bg-gray-50 focus:bg-white focus:border-blue-500 focus:ring-4 focus:ring-blue-100 transition-all duration-200 outline-none"
            >
          </div>

          <!-- E-posta Alanı-->
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

          <!-- Buton Alanı -->
          <div class="md:col-span-2 pt-6">
            <button 
              @click="updateProfile" 
              class="w-full md:w-max px-10 py-4 bg-blue-600 text-white text-xl font-bold rounded-xl hover:bg-blue-700 hover:shadow-lg active:transform active:scale-[0.98] transition-all duration-200 flex items-center justify-center space-x-2 shadow-blue-200"
            >
              <span>Değişiklikleri Kaydet</span>
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </button>
          </div>

              <!-- Şifre Değiştirme Bölümü -->
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
                    <input v-model="passwords.oldPassword" type="password" placeholder="••••••••" class="block w-full px-5 py-4 text-lg border-2 border-gray-100 rounded-xl bg-gray-50 focus:bg-white focus:border-blue-500 focus:ring-4 focus:ring-blue-100 transition-all duration-200 outline-none">
                    </div>

                    <div class="space-y-2">
                    <label class="block text-sm font-semibold text-gray-600 uppercase tracking-wider ml-1">Yeni Şifre</label>
                    <input v-model="passwords.newPassword" type="password" placeholder="••••••••" class="block w-full px-5 py-4 text-lg border-2 border-gray-100 rounded-xl bg-gray-50 focus:bg-white focus:border-blue-500 focus:ring-4 focus:ring-blue-100 transition-all duration-200 outline-none">
                    </div>

                    <div class="md:col-span-2">
                    <button @click="changePassword" class="w-full md:w-max px-10 py-4 bg-gray-800 text-white text-xl font-bold rounded-xl hover:bg-black hover:shadow-lg active:transform active:scale-[0.98] transition-all duration-200 shadow-gray-200">
                        Şifreyi Güncelle
                    </button>
                    </div>
                </div>
                </div>
          
        </div>

        <div v-else class="flex flex-col items-center justify-center py-20 space-y-4">
          <div class="animate-spin rounded-full h-12 w-12 border-b-4 border-blue-600"></div>
          <p class="text-xl text-gray-500 font-medium">Bilgileriniz getiriliyor...</p>
        </div>
      </div>
    </div>

      <div class="md:col-span-2 mt-12 pt-8 border-t-2 border-gray-100">
      <h3 class="text-2xl font-bold text-gray-800 mb-8 flex items-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 mr-2 text-red-500" fill="currentColor" viewBox="0 0 24 24">
          <path d="M11.645 20.91l-.007-.003-.022-.012a15.247 15.247 0 01-.383-.218 25.18 25.18 0 01-4.244-3.17C4.688 15.36 2.25 12.174 2.25 8.25 2.25 5.322 4.714 3 7.688 3c1.74 0 3.285.797 4.312 2.022C13.027 3.797 14.572 3 16.312 3c2.974 0 5.438 2.322 5.438 5.25 0 3.924-2.438 7.111-4.739 9.256a25.175 25.175 0 01-4.244 3.17 15.247 15.247 0 01-.383.219l-.022.012-.007.004-.003.001z" />
        </svg>
        Favori Eserlerim
      </h3>

      <div v-if="favorites.length > 0" class="grid grid-cols-1 sm:grid-cols-3 gap-6">
        <ArtworkCard 
          v-for="fav in favorites" 
          :key="fav.id" 
          v-bind="fav"
        />
      </div>

      <div v-else class="text-center py-10 bg-gray-50 rounded-2xl border-2 border-dashed border-gray-200">
        <p class="text-gray-400 italic">Henüz bir eseri favorilere eklemediniz.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import ArtworkCard from '../components/ArtworkCard.vue';

const user = ref(null);
const favorites = ref([]);

const passwords = ref({
  oldPassword: '',
  newPassword: ''
});

// Sayfa yüklendiğinde bilgileri getir
onMounted(async () => {
  const userEmail = localStorage.getItem('userEmail'); 
  if (!userEmail) return;

  try {
    // 1. Profil bilgilerini getir
    const profileRes = await axios.get(`http://localhost:8080/profile?email=${userEmail}`);
    user.value = profileRes.data;

    // 2. Favorileri getir (Yeni Eklediğimiz Kısım)
    const favRes = await axios.get(`http://localhost:8080/favorites/list?email=${userEmail}`);
    favorites.value = favRes.data || [];
    
  } catch (error) {
    console.error("Veriler yüklenemedi:", error);
  }
});



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
    // Başarılı olduktan sonra kutuları temizleyelim
    passwords.value.oldPassword = '';
    passwords.value.newPassword = '';
  } catch (error) {
    console.error("Şifre değiştirme hatası:", error);
    alert(error.response?.data || "Şifre değiştirilemedi. Eski şifrenizi kontrol edin.");
  }
};
</script>