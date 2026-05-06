<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100 p-4">
    <div class="bg-white p-12 rounded-[30px] shadow-2xl w-full max-w-md">
      <div class="text-center mb-10">
        <h2 class="text-3xl font-bold text-gray-800">Hoş Geldiniz</h2>
        <p class="text-gray-500 mt-2">Koleksiyonunuza erişmek için giriş yapın.</p>
      </div>

      <div class="space-y-6">
        <div class="relative">
          <input v-model="loginData.email" type="email" placeholder="E-posta" 
            class="w-full p-4 pl-12 border border-gray-200 rounded-2xl focus:outline-none focus:border-galeri-yesil focus:ring-4 focus:ring-green-50 transition-all" />
          <span class="absolute left-4 top-4 opacity-50">✉️</span>
        </div>

        <div class="relative">
          <input v-model="loginData.password" type="password" placeholder="Şifre" 
            class="w-full p-4 pl-12 border border-gray-200 rounded-2xl focus:outline-none focus:border-galeri-yesil focus:ring-4 focus:ring-green-50 transition-all" />
          <span class="absolute left-4 top-4 opacity-50">🔑</span>
        </div>

        <button @click="handleLogin" 
          class="w-full bg-galeri-yesil text-white font-bold py-4 rounded-2xl shadow-lg hover:bg-green-600 transition-all text-lg">
          Giriş Yap
        </button>
      </div>

      <p class="text-center text-gray-500 mt-8">
        Hesabınız yok mu? 
        <span @click="router.push('/register')" class="text-galeri-yesil font-bold cursor-pointer hover:underline">Kaydol</span>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()
const loginData = ref({ email: '', password: '' })

const handleLogin = async () => {
  try {
    const response = await axios.post('http://localhost:8080/login', loginData.value)
    
    if (response.data && response.data.token) {
      localStorage.setItem('userToken', response.data.token) 
      
      localStorage.setItem('userEmail', loginData.value.email) 
      
      alert('Giriş başarılı!')
      router.push('/')
    }
  } catch (error) {
    alert('Hata: ' + (error.response?.data || 'Giriş yapılamadı'))
  }
}
</script>