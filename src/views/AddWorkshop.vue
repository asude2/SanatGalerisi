<template>
  <div class="min-h-screen bg-gray-50 p-4 md:p-12">
    <div class="max-w-3xl mx-auto bg-white rounded-[40px] shadow-2xl overflow-hidden border border-gray-100">
      <div class="bg-blue-600 p-10 text-white">
        <button @click="router.push('/')" class="mb-6 opacity-80 hover:opacity-100 flex items-center gap-2">
          ⬅️ Galeriye Dön
        </button>
        <h1 class="text-4xl font-black italic tracking-tighter">Yeni Atölye Oluştur 🏛️</h1>
        <p class="opacity-90 mt-2">Bilgi ve tecrübelerinizi sanatseverlerle paylaşmak için bir workshop planlayın.</p>
      </div>

      <form @submit.prevent="handleSubmit" class="p-10 space-y-6">
        
        <div>
          <label class="block text-gray-700 font-bold mb-2 ml-1">Atölye Başlığı</label>
          <input 
            v-model="workshop.title" 
            type="text" 
            placeholder="Örn: Akrilik Boya Teknikleri Giriş" 
            class="w-full p-4 bg-gray-50 border-none rounded-2xl focus:ring-4 focus:ring-blue-50 transition-all text-lg"
            required
          />
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label class="block text-gray-700 font-bold mb-2 ml-1">Konum / Platform</label>
            <input 
              v-model="workshop.location" 
              type="text" 
              placeholder="Örn: Online (Zoom) veya Atölye Adresi" 
              class="w-full p-4 bg-gray-50 border-none rounded-2xl focus:ring-4 focus:ring-blue-50 transition-all text-lg"
              required
            />
          </div>
          <div>
            <label class="block text-gray-700 font-bold mb-2 ml-1">Kontenjan (Kişi Sayısı)</label>
            <input 
              v-model.number="workshop.capacity" 
              type="number" 
              placeholder="10" 
              class="w-full p-4 bg-gray-50 border-none rounded-2xl focus:ring-4 focus:ring-blue-50 transition-all text-lg"
              required
            />
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label class="block text-gray-700 font-bold mb-2 ml-1">Katılım Ücreti (₺)</label>
            <input 
              v-model.number="workshop.price" 
              type="number" 
              placeholder="0.00" 
              class="w-full p-4 bg-gray-50 border-none rounded-2xl focus:ring-4 focus:ring-blue-50 transition-all text-lg"
              required
            />
          </div>
          <div>
            <label class="block text-gray-700 font-bold mb-2 ml-1">Atölye Tarihi</label>
            <input 
              v-model="workshop.availableDates" 
              type="date" 
              class="w-full p-4 bg-gray-50 border-none rounded-2xl focus:ring-4 focus:ring-blue-50 transition-all text-lg"
              required
            />
          </div>
        </div>

        <div>
          <label class="block text-gray-700 font-bold mb-2 ml-1">Kapak Görseli URL</label>
          <input 
            v-model="workshop.imageUrl" 
            type="text" 
            placeholder="https://.../workshop.jpg" 
            class="w-full p-4 bg-gray-50 border-none rounded-2xl focus:ring-4 focus:ring-blue-50 transition-all text-lg"
            required
          />
        </div>

        <div>
          <label class="block text-gray-700 font-bold mb-2 ml-1">Atölye Detayları ve Program</label>
          <textarea 
            v-model="workshop.description" 
            rows="4" 
            placeholder="Eğitim içeriğinde neler var? Kimler katılmalı?" 
            class="w-full p-4 bg-gray-50 border-none rounded-2xl focus:ring-4 focus:ring-blue-50 transition-all text-lg resize-none"
            required
          ></textarea>
        </div>

        <button 
          type="submit" 
          class="w-full bg-blue-600 text-white font-black py-5 rounded-2xl shadow-xl shadow-blue-100 hover:scale-[1.02] transition-transform text-xl mt-4"
        >
          Atölyeyi Planla ve Duyur 🚀
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()

const workshop = ref({
  title: '',
  description: '',
  location: '',
  capacity: 10,
  price: 0,
  imageUrl: '',
  availableDates: '' // Template'deki type="date" sayesinde burası '2026-05-20' formatında dolacak
})

const handleSubmit = async () => {
  try {
    // 1. Dinamik olarak giriş yapan kullanıcının mailini çekiyoruz
    const userEmail = localStorage.getItem('userEmail');

    if (!userEmail) {
      alert("Oturum süreniz dolmuş veya giriş yapmadınız. Lütfen tekrar giriş yapın. 👤");
      router.push('/login');
      return;
    }

    // 2. Go Backend'e veriyi gönderiyoruz
    // URL'yi senin backend'deki '/add-workshop' rotasıyla tam eşledik
    const response = await axios.post('http://localhost:8080/add-workshop', {
      email: userEmail, 
      title: workshop.value.title,
      description: workshop.value.description,
      location: workshop.value.location,
      capacity: parseInt(workshop.value.capacity),
      price: parseFloat(workshop.value.price),
      imageUrl: workshop.value.imageUrl,
      availableDates: workshop.value.availableDates // Artık tertemiz YYYY-MM-DD gidecek
    });

    if (response.status === 201 || response.status === 200) {
      alert("Atölyeniz başarıyla oluşturuldu ve sisteme kaydedildi! ✨🏛️");
      router.push('/workshops'); 
    }
  } catch (error) {
    console.error("Atölye ekleme hatası:", error);
    
    // Hata mesajını kullanıcıya daha anlamlı gösterelim
    const errorMessage = error.response?.data?.message || error.response?.data || "Sunucuyla bağlantı kurulamadı.";
    alert("Atölye oluşturulurken bir hata meydana geldi: " + errorMessage);
  }
}
</script>