<template>
  <div class="min-h-screen bg-gray-50 py-12 px-6">
    <div v-if="workshop" class="max-w-4xl mx-auto bg-white rounded-3xl shadow-2xl overflow-hidden">
      <div class="relative h-96">
        <img :src="workshop.image" class="w-full h-full object-cover" />
        <button @click="router.back()" class="absolute top-6 left-6 bg-white/90 p-3 rounded-full shadow-lg cursor-pointer">
          ⬅️ Geri Dön
        </button>
      </div>

      <div class="p-10">
        <div class="flex justify-between items-start mb-6">
          <div>
            <h1 class="text-4xl font-black text-gray-900">{{ workshop.title }}</h1>
            <p class="text-blue-600 font-bold text-lg mt-2 italic">Eğitmen: {{ workshop.instructorName }}</p>
          </div>
          <div class="bg-blue-50 text-blue-700 px-6 py-3 rounded-2xl font-black text-2xl">
            {{ workshop.price > 0 ? workshop.price + ' ₺' : 'Ücretsiz' }}
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-10 border-y py-8 border-gray-100">
          <div class="flex items-center gap-3">
            <span class="text-3xl">📅</span>
            <div>
              <p class="text-gray-400 text-xs font-bold uppercase">Mevcut Tarihler</p>
              <p class="font-bold text-gray-800">{{ workshop.availableDates }}</p>
            </div>
          </div>
          <div class="flex items-center gap-3">
            <span class="text-3xl">📍</span>
            <div>
              <p class="text-gray-400 text-xs font-bold uppercase">Konum</p>
              <p class="font-bold text-gray-800">{{ workshop.location }}</p>
            </div>
          </div>
          <div class="flex items-center gap-3">
            <span class="text-3xl">👥</span>
            <div>
              <p class="text-gray-400 text-xs font-bold uppercase">Toplam Kontenjan</p>
              <p class="font-bold text-gray-800">{{ workshop.capacity }} Kişi</p>
            </div>
          </div>
        </div>

        <div class="mb-10">
          <h3 class="text-xl font-bold text-gray-800 mb-4">Etkinlik Hakkında</h3>
          <p class="text-gray-600 leading-relaxed text-lg whitespace-pre-line text-justify">
            {{ workshop.description }}
          </p>
        </div>

        <div class="bg-blue-50 p-6 rounded-2xl mb-8 border border-blue-100 space-y-6">
          <h3 class="text-lg font-bold text-blue-900 flex items-center gap-2">
            ✨ Rezervasyon Detaylarını Belirleyin
          </h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="space-y-2">
              <label class="text-sm font-bold text-gray-600 ml-1">Katılımcı Sayısı</label>
              <select v-model.number="reservation.participantCount" class="w-full p-4 bg-white rounded-xl border-2 border-transparent focus:border-blue-500 outline-none transition-all shadow-sm">
                <option v-for="n in 10" :key="n" :value="n">{{ n }} Kişi</option>
              </select>
            </div>
            <div class="space-y-2">
              <label class="text-sm font-bold text-gray-600 ml-1">Tercih Edilen Tarih</label>
              <select 
                v-model="reservation.date" 
                class="w-full p-4 bg-white rounded-xl border-2 border-transparent focus:border-blue-500 outline-none transition-all shadow-sm"
              >
                <option value="" disabled selected>Lütfen bir tarih seçin...</option>
                <option v-for="date in workshop.availableDates?.split(',')" :key="date" :value="date.trim()">
                  {{ date.trim() }}
                </option>
              </select>
            </div>
          </div>

          <div class="space-y-2 pt-2">
            <label class="text-sm font-bold text-gray-600 ml-1">Ödeme Yöntemi</label>
            <select v-model="selectedPaymentMethod" class="w-full p-4 bg-white rounded-xl border-2 border-transparent focus:border-blue-500 outline-none transition-all shadow-sm">
              <option value="Kredi Kartı">Kredi Kartı</option>
              <option value="Uygulama Bakiyesi">Uygulama Bakiyesi</option>
            </select>
          </div>

          <div class="p-4 bg-white rounded-xl border border-blue-200 flex justify-between items-center mt-4 shadow-sm">
            <div>
              <span class="text-xs text-blue-600 font-bold uppercase block">Toplam Ödenecek Tutar</span>
              <span class="text-3xl font-black text-blue-900">{{ totalWorkshopPrice.toLocaleString() }} ₺</span>
            </div>
            <span class="text-xs text-gray-400 font-medium">({{ reservation.participantCount }} x {{ workshop.price }} ₺)</span>
          </div>
        </div>

        <button 
          @click="enroll"
          class="w-full py-5 bg-blue-600 text-white rounded-2xl font-black text-xl hover:bg-blue-700 transition-all shadow-xl shadow-blue-100 cursor-pointer active:scale-[0.98]"
        >
          Ödeme Yap ve Rezervasyon Oluştur 💳
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import { jwtDecode } from 'jwt-decode' 

const route = useRoute()
const router = useRouter()
const workshop = ref(null)

const reservation = ref({
  participantCount: 1,
  date: ''
})

// 🚀 YENİ: Ödeme yöntemi reactive statetimiz (Varsayılan: Kredi Kartı)
const selectedPaymentMethod = ref('Kredi Kartı')

// Katılımcı sayısı ile ham fiyatı çarpan computed property
const totalWorkshopPrice = computed(() => {
  if (!workshop.value) return 0
  return workshop.value.price * reservation.value.participantCount
})

onMounted(async () => {
  try {
    const response = await axios.get('http://localhost:8080/workshops')
    const all = response.data
    workshop.value = all.find(w => w.id === parseInt(route.params.id))
  } catch (error) {
    console.error("Detay yüklenemedi:", error)
  }
})

const enroll = async () => {
  const token = localStorage.getItem('userToken')
  if (!token) {
    alert("Rezervasyon yapmak için lütfen giriş yapın! 👤")
    return
  }

  if (!reservation.value.date) {
    alert("Lütfen bir tarih seçiniz! 📅")
    return
  }

  const confirmMessage = `${workshop.value.title} atölyesine ${reservation.value.participantCount} kişi için [${selectedPaymentMethod.value}] kullanılarak toplam ${totalWorkshopPrice.value.toLocaleString()} ₺ ödeme yapılacaktır. Onaylıyor musunuz?`
  if (!confirm(confirmMessage)) return

  try {
    const decoded = jwtDecode(token)
    // 🚀 Geliştirilen payload: Seçilen ödeme yöntemini de backend'e paslıyoruz
    const response = await axios.post('http://localhost:8080/enroll-workshop-payment', {
      email: decoded.email,
      workshopId: workshop.value.id,
      participantCount: parseInt(reservation.value.participantCount),
      reservedDate: reservation.value.date,
      paymentMethod: selectedPaymentMethod.value
    })

    alert(response.data.message || "Rezervasyon ve ödeme başarılı! 🎉")
    router.push('/profile') 
  } catch (error) {
    alert(error.response?.data?.message || "Rezervasyon sırasında bir hata oluştu veya bakiyeniz yetersiz! 💸")
  }
}
</script>