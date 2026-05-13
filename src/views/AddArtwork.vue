<template>
  <div class="min-h-screen bg-gray-50 p-4 md:p-12">
    <div class="max-w-3xl mx-auto bg-white rounded-[40px] shadow-2xl overflow-hidden border border-gray-100">
      <div class="bg-galeri-yesil p-10 text-white">
        <button @click="router.push('/')" class="mb-6 opacity-80 hover:opacity-100 flex items-center gap-2">
          ⬅️ Galeriye Dön
        </button>
        <h1 class="text-4xl font-black italic tracking-tighter">Yeni Eser Sergile 🎨</h1>
        <p class="opacity-90 mt-2">Eserinizin detaylarını girerek sanatseverlerle buluşturun.</p>
      </div>

      <form @submit.prevent="handleSubmit" class="p-10 space-y-6">
        <div>
          <label class="block text-gray-700 font-bold mb-2 ml-1">Eserin Adı</label>
          <input 
            v-model="artwork.title" 
            type="text" 
            placeholder="Örn: Yıldızlı Gece" 
            class="w-full p-4 bg-gray-50 border-none rounded-2xl focus:ring-4 focus:ring-green-50 transition-all text-lg"
            required
          />
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label class="block text-gray-700 font-bold mb-2 ml-1">Fiyat (₺)</label>
            <input 
              v-model="artwork.price" 
              type="number" 
              placeholder="0.00" 
              class="w-full p-4 bg-gray-50 border-none rounded-2xl focus:ring-4 focus:ring-green-50 transition-all text-lg"
              required
            />
          </div>

          <div>
            <label class="block text-gray-700 font-bold mb-2 ml-1">Kategori</label>
            <select 
              v-model="artwork.category" 
              class="w-full p-4 bg-gray-50 border-none rounded-2xl focus:ring-4 focus:ring-green-50 transition-all text-lg appearance-none"
            >
            <option value="Manzara">Manzara</option>
            <option value="Rönesans">Rönesans</option>
            <option value="Realizm">Realizm</option>
            <option value="Natürmort">Natürmort</option>
            <option value="Modern Sanat">Modern Sanat</option>
            </select>
          </div>
        </div>

        <div>
          <label class="block text-gray-700 font-bold mb-2 ml-1">Eser Görsel URL (Link)</label>
          <input 
            v-model="artwork.imageUrl" 
            type="text" 
            placeholder="https://.../eser.jpg" 
            class="w-full p-4 bg-gray-50 border-none rounded-2xl focus:ring-4 focus:ring-green-50 transition-all text-lg"
            required
          />
        </div>

        <div>
          <label class="block text-gray-700 font-bold mb-2 ml-1">Eserin Hikayesi / Açıklama</label>
          <textarea 
            v-model="artwork.description" 
            rows="4" 
            placeholder="Bu eser hakkında kısa bir bilgi verin..." 
            class="w-full p-4 bg-gray-50 border-none rounded-2xl focus:ring-4 focus:ring-green-50 transition-all text-lg resize-none"
          ></textarea>
        </div>

        <button 
          type="submit" 
          class="w-full bg-galeri-yesil text-white font-black py-5 rounded-2xl shadow-xl shadow-green-100 hover:scale-[1.02] transition-transform text-xl mt-4"
        >
          Esere Hayat Ver ve Yayınla ✨
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

const artwork = ref({
  title: '',
  price: 0,
  category: 'Tablo',
  imageUrl: '',
  description: ''
})

const handleSubmit = async () => {
  try {
    const userEmail = localStorage.getItem('userEmail')
    const response = await axios.post('http://localhost:8080/add-artwork', {
      title: artwork.value.title,
      price: parseFloat(artwork.value.price), 
      category: artwork.value.category,
      imageUrl: artwork.value.imageUrl,
      description: artwork.value.description,
      email: userEmail
    });

    if (response.status === 201) {
      alert("Eseriniz başarıyla yayınlandı! 🎨");
      router.push('/');
    }
  } catch (error) {
    console.error("Eser ekleme hatası:", error);
    alert("Eser eklenirken bir hata oluştu. Lütfen tekrar deneyin.");
  }
}
</script>