<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100 p-4">
    <div class="bg-white p-12 rounded-[30px] shadow-2xl w-full max-w-4xl flex flex-col md:flex-row items-center gap-12">
      
      <!-- Sol Panel -->
      <div class="hidden md:flex flex-1 flex-col items-center justify-center border-r border-gray-100 pr-12 text-center">
        <div class="text-[120px] leading-none mb-6 animate-pulse">🖌️</div>
        <h1 class="text-4xl font-extrabold text-galeri-yesil italic tracking-tighter">
          ARAMIZA KATILIN
        </h1>
        <p class="text-gray-400 mt-4 max-w-md">
          Kendi sanat koleksiyonunuzu oluşturmaya başlamak için ilk adımı atın.
        </p>
      </div>

      <!-- Sağ Panel (Form) -->
      <div class="w-full md:w-1/2 space-y-5">
        <div class="text-center md:text-left mb-6">
          <h2 class="text-3xl font-bold text-gray-800">Hesap Oluştur</h2>
          <p class="text-gray-500 mt-2">Lütfen bilgilerinizi eksiksiz doldurun.</p>
        </div>

        <!-- Ad ve Soyad -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <input 
            v-model="user.firstName" 
            type="text" 
            placeholder="Ad" 
            class="w-full p-4 border border-gray-200 rounded-2xl focus:outline-none focus:border-galeri-yesil focus:ring-4 focus:ring-green-50 transition-all" 
          />
          <input 
            v-model="user.lastName" 
            type="text" 
            placeholder="Soyad" 
            class="w-full p-4 border border-gray-200 rounded-2xl focus:outline-none focus:border-galeri-yesil focus:ring-4 focus:ring-green-50 transition-all" 
          />
        </div>

        <!-- E-posta -->
        <div class="relative">
          <input 
            v-model="user.email" 
            type="email" 
            placeholder="E-posta Adresi" 
            class="w-full p-4 pl-12 border border-gray-200 rounded-2xl focus:outline-none focus:border-galeri-yesil focus:ring-4 focus:ring-green-50 transition-all" 
          />
          <span class="absolute left-4 top-4 text-xl opacity-50">✉️</span>
        </div>

        <!-- Şifre -->
        <div class="relative">
          <input 
            v-model="user.password" 
            type="password" 
            placeholder="Şifre" 
            class="w-full p-4 pl-12 border border-gray-200 rounded-2xl focus:outline-none focus:border-galeri-yesil focus:ring-4 focus:ring-green-50 transition-all" 
          />
          <span class="absolute left-4 top-4 text-xl opacity-50">🔑</span>
        </div>

        <!-- Şifre Tekrar -->
        <div class="relative">
          <input 
            v-model="user.confirmPassword" 
            type="password" 
            placeholder="Şifre Tekrar" 
            class="w-full p-4 pl-12 border border-gray-200 rounded-2xl focus:outline-none focus:border-galeri-yesil focus:ring-4 focus:ring-green-50 transition-all" 
          />
          <span class="absolute left-4 top-4 text-xl opacity-50">🔒</span>
        </div>

        <div class="mb-4">
          <label class="block text-gray-700 text-sm font-bold mb-2">Hesap Türü</label>
          <select 
            v-model="user.role" 
            class="w-full p-4 border border-gray-200 rounded-2xl focus:outline-none focus:border-galeri-yesil focus:ring-4 focus:ring-green-50 transition-all appearance-none bg-white"
          >
            <option value="User">Sanatsever (Kullanıcı)</option>
            <option value="Instructor">Eğitmen / Sanatçı</option>
          </select>
        </div>

        <button 
          @click="handleRegister" 
          class="w-full bg-galeri-yesil text-white font-bold py-4 rounded-2xl shadow-lg hover:bg-green-600 transform hover:-translate-y-1 transition-all text-lg mt-4"
        >
          Kayıt Ol
        </button>

        <p class="text-center text-gray-500 pt-4">
          Zaten üye misin? 
          <span @click="router.push('/login')" class="text-galeri-yesil font-bold cursor-pointer hover:underline">
            Giriş Yap
          </span>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios' 

const router = useRouter()
const user = ref({
  firstName: '',
  lastName: '',
  email: '',
  password: '',
  confirmPassword: '',
  role: 'User'
})

const handleRegister = async () => {
  if (!user.value.firstName || !user.value.email || !user.value.password) {
    alert('Lütfen tüm zorunlu alanları doldurun!')
    return
  }

  if (user.value.password !== user.value.confirmPassword) {
    alert('Şifreler uyuşmuyor!')
    return
  }

  try {
    // 8080 portundaki Go sunucumuza veriyi gönderiyoruz[cite: 1]
    const response = await axios.post('http://localhost:8080/register', {
      firstName: user.value.firstName,
      lastName: user.value.lastName,
      email: user.value.email,
      password: user.value.password,
      role: user.value.role 
    })

    if (response.status === 201) {
      alert('Kaydınız başarıyla veritabanına iletildi!')
      router.push('/login')
    }
  } catch (error) {
    console.error('Bağlantı hatası:', error)
    alert('Backend sunucusuna ulaşılamadı. Sunucunun çalıştığından emin olun.')
  }
}
</script>