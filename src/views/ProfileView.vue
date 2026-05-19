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
      <div class="bg-gradient-to-r from-blue-600 to-indigo-700 px-8 py-10 text-white">
        <button @click="goBack" class="mb-5 p-2 bg-white/20 hover:bg-white/30 rounded-full transition-all group" title="Geri Git">
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
            <h1 class="text-3xl font-black text-gray-900">{{ user.firstName }} {{ user.lastName }}</h1>
            <p class="text-gray-500 font-medium">{{ user.email }} • <span class="text-blue-600">{{ user.role }}</span></p>
            <p v-if="user.biography" class="mt-4 text-gray-600 leading-relaxed max-w-2xl">{{ user.biography }}</p>
          </div>
        </div>
      </div>

      <!-- Navigasyon Sekmeleri -->
      <div class="flex gap-2 mb-8 bg-gray-200/50 p-1.5 rounded-2xl w-fit overflow-x-auto">
        <button 
          v-for="tab in tabs" 
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


            <br><br><br><br><br><br>

            <div class="bg-blue-50 border-2 border-blue-100 rounded-2xl p-6 mb-8 flex items-center justify-between">
              <div class="flex items-center space-x-4">
                <div class="p-3 bg-blue-600 rounded-xl text-white shadow-lg text-2xl">
                  💰
                </div>
                <div>
                  <p class="text-sm font-bold text-blue-600 uppercase tracking-wider">Hesap Bakiyeniz</p>
                  <h3 class="text-3xl font-black text-blue-900">
                    {{ user?.balance ? user.balance.toLocaleString() : '0' }} TL
                  </h3>
                </div>
              </div>
              
              <button 
                @click="addBalance" 
                class="px-6 py-3 bg-white border-2 border-blue-600 text-blue-600 font-bold rounded-xl hover:bg-blue-600 hover:text-white transition-all shadow-sm"
              >
                + Bakiye Yükle
              </button>
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
                  <h4 class="text-xl font-bold text-gray-800 mb-2">{{ enroll.workshopTitle }}</h4>
                  <div class="flex flex-wrap gap-4 text-sm text-gray-600 mb-4">
                    <span class="flex items-center">📅 {{ enroll.reservedDate }}</span>
                    <span class="flex items-center">👥 {{ enroll.participantCount }} Kişi</span>
                    <span class="flex items-center">📍 {{ enroll.location }}</span>
                  </div>
                  
                  <span 
                    :class="enroll.status === 'Onaylandı' ? 'bg-green-100 text-green-700 border-green-200' : 'bg-amber-100 text-amber-700 border-amber-200'"
                    class="px-4 py-1.5 rounded-full text-xs font-black uppercase tracking-widest border"
                  >
                    {{ enroll.status }}
                  </span>
                </div>
                
                <div class="flex gap-2">
                  <button 
                    v-if="enroll.status !== 'Onaylandı'"
                    @click="openEditModal(enroll)" 
                    class="px-4 py-2 bg-blue-100 text-blue-700 font-bold rounded-lg hover:bg-blue-200 transition-colors"
                  >
                    Düzenle
                  </button>
                  <button 
                    v-if="enroll.status !== 'Onaylandı'"
                    @click="cancelEnrollment(enroll.id)" 
                    class="px-4 py-2 bg-red-100 text-red-700 font-bold rounded-lg hover:bg-red-200 transition-colors"
                  >
                    İptal Et
                  </button>
                  
                  <span 
                    v-else 
                    class="text-green-600 font-bold text-sm flex items-center gap-1 bg-green-50 px-4 py-2 rounded-xl border border-green-100"
                  >
                    ✨ Rezervasyon Kesinleşti
                  </span>
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
                  :key="fav.id || fav.Id" 
                  :artwork="fav" 
                />
              </div>

              <div v-else class="text-center py-10 bg-gray-50 rounded-2xl border-2 border-dashed border-gray-200">
                <p class="text-gray-400 italic">Henüz bir eseri favorilere eklemediniz.</p>
              </div>
            </div>
            
            <div v-if="userRole === 'Instructor'" class="md:col-span-2 mt-12 pt-8 border-t-2 border-gray-100">
              <h3 class="text-2xl font-bold text-gray-800 mb-8 flex items-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 mr-2 text-orange-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
                </svg>
                🎨 Gelen Siparişler (Satış Onayı)
              </h3>

              <div v-if="sellerOrders && sellerOrders.length > 0" class="overflow-hidden bg-white border border-gray-100 rounded-2xl shadow-sm">
                <table class="min-w-full divide-y divide-gray-200 text-left">
                  <thead class="bg-gray-50">
                    <tr>
                      <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Eser Bilgisi</th>
                      <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Alıcı</th>
                      <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Fiyat</th>
                      <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Durum</th>
                      <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider text-right">İşlem</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-gray-200">
                    <tr v-for="order in sellerOrders" :key="order.id" class="hover:bg-gray-50 transition-colors">
                      <td class="px-6 py-4 font-bold text-gray-900">{{ order.title }}</td>
                      <td class="px-6 py-4 text-sm text-gray-600">{{ order.email }}</td>
                      <td class="px-6 py-4 font-semibold text-blue-600">{{ order.price }} TL</td>
                      <td class="px-6 py-4">
                        <span :class="order.status === 'Onaylandı' ? 'bg-green-100 text-green-700' : 'bg-orange-100 text-orange-700'"
                              class="px-3 py-1 rounded-full text-xs font-bold">
                          {{ order.status }}
                        </span>
                      </td>
                      <td class="px-6 py-4 text-right">
                        <button 
                          v-if="order.status === 'Hazırlanıyor'"
                          @click="confirmSale(order.id)" 
                          class="bg-blue-600 text-white px-4 py-2 rounded-xl text-sm font-bold hover:bg-blue-700 shadow-lg shadow-blue-100 transition-all"
                        >
                          Onayla
                        </button>
                        <span v-else class="text-gray-400 text-sm italic">Tamamlandı</span>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <div v-else class="text-center py-10 bg-gray-50 rounded-2xl border-2 border-dashed border-gray-200">
                <p class="text-gray-400 italic">Henüz bir satış talebi almadınız.</p>
              </div>
            </div>

            <div v-if="userRole === 'Instructor'" class="md:col-span-2 mt-12 pt-8 border-t-2 border-gray-100">
              <h3 class="text-2xl font-bold text-gray-800 mb-8 flex items-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 mr-2 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
                🎫 Gelen Atölye Başvuruları (Rezervasyon Onayı)
              </h3>

              <div v-if="workshopOrders && workshopOrders.length > 0" class="overflow-hidden bg-white border border-gray-100 rounded-2xl shadow-sm">
                <table class="min-w-full divide-y divide-gray-200 text-left">
                  <thead class="bg-gray-50">
                    <tr>
                      <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Atölye Adı</th>
                      <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Katılımcı E-posta</th>
                      <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Kişi Sayısı</th>
                      <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Seçilen Tarih</th>
                      <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Durum</th>
                      <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider text-right">İşlem</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-gray-200">
                    <tr v-for="order in workshopOrders" :key="order.id" class="hover:bg-gray-50 transition-colors">
                      <td class="px-6 py-4 font-bold text-gray-900">{{ order.title }}</td>
                      <td class="px-6 py-4 text-sm text-gray-600">{{ order.email }}</td>
                      <td class="px-6 py-4 font-medium text-gray-700">{{ order.participantCount }} Kişi</td>
                      <td class="px-6 py-4 text-sm text-gray-500">{{ order.reservedDate }}</td>
                      <td class="px-6 py-4">
                        <span :class="order.status === 'Onaylandı' ? 'bg-green-100 text-green-700' : 'bg-amber-100 text-amber-700'"
                              class="px-3 py-1 rounded-full text-xs font-bold">
                          {{ order.status }}
                        </span>
                      </td>
                      <td class="px-6 py-4 text-right">
                        <button 
                          v-if="order.status !== 'Onaylandı'"
                          @click="confirmWorkshopOrder(order.id)" 
                          class="bg-indigo-600 text-white px-4 py-2 rounded-xl text-sm font-bold hover:bg-indigo-700 shadow-lg shadow-indigo-100 transition-all"
                        >
                          Başvuruyu Onayla ✅
                        </button>
                        <span v-else class="text-green-600 font-bold text-sm italic">Onaylandı</span>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <div v-else class="text-center py-10 bg-gray-50 rounded-2xl border-2 border-dashed border-gray-200">
                <p class="text-gray-400 italic">Henüz bir atölye rezervasyon başvurusu almadınız.</p>
              </div>
            </div>

            <div v-if="userRole === 'Instructor'" class="md:col-span-2 mt-12 pt-8 border-t-2 border-gray-100">
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
            <select v-model="selectedEnrollment.reservedDate" class="w-full p-4 bg-gray-50 border-2 border-gray-100 rounded-xl outline-none focus:border-blue-500">
              <option :value="selectedEnrollment.reservedDate">{{ selectedEnrollment.reservedDate }} (Mevcut)</option>
              <option v-for="date in selectedEnrollment.availableDates?.split(',')" :key="date" :value="date.trim()">
                {{ date.trim() }}
              </option>
            </select>
          </div>

          <div v-if="selectedEnrollment.participantCount > originalParticipantCount" class="space-y-4 pt-4 border-t border-gray-100">
            <div class="p-4 bg-amber-50 rounded-xl border border-amber-200 text-sm text-amber-800 font-medium">
              ⚠️ Katılımcı sayısını artırdığınız için <strong>{{ (selectedEnrollment.participantCount - originalParticipantCount) * selectedEnrollment.workshopPrice }} TL</strong> ekstra ödeme yapmanız gerekmektedir.
            </div>
            
            <div class="space-y-2">
              <label class="block text-sm font-bold text-gray-600 ml-1 uppercase">Ödeme Yöntemi</label>
              <select v-model="editModalPaymentMethod" class="w-full p-4 bg-white border-2 border-gray-200 rounded-xl outline-none focus:border-blue-500 shadow-sm">
                <option value="Kredi Kartı">Kredi Kartı</option>
                <option value="Uygulama Bakiyesi">Uygulama Bakiyesi</option>
              </select>
            </div>
          </div>

          <div class="flex justify-end space-x-3 mt-8">
            <button @click="showEditModal = false" class="px-6 py-3 font-bold text-gray-500 hover:text-gray-700">Vazgeç</button>
            <button @click="updateEnrollment" class="px-8 py-3 bg-blue-600 text-white font-bold rounded-xl hover:shadow-lg transition-all">
              {{ selectedEnrollment.participantCount > originalParticipantCount ? 'Öde ve Güncelle' : 'Güncelle' }}
            </button>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import { ref, onMounted, computed, watch, nextTick } from 'vue';
import axios from 'axios';
import ArtworkCard from '../components/ArtworkCard.vue';
import WorkshopCard from '../components/WorkshopCard.vue';
import { useRouter } from 'vue-router';

const router = useRouter();
const user = ref(null);
const favorites = ref([]);
const purchases = ref([]);
const comparisons = ref([]);
const activeTab = ref('favorites');
const editingId = ref(null);
const editTitle = ref('');

const tabs = [
  { id: 'favorites', label: 'Favorilerim' },
  { id: 'orders', label: 'Siparişlerim' },
  { id: 'comparisons', label: 'Analizlerim' },
  { id: 'edit', label: 'Düzenle' }
];
// 🚀 YENİ: Ekstra ödeme durumlarını izlemek için gereken modal state'leri
const originalParticipantCount = ref(0);
const editModalPaymentMethod = ref('Kredi Kartı');

const passwords = ref({
  oldPassword: '',
  newPassword: ''
});
const sellerOrders = ref([]);
const workshopOrders = ref([]); // Gelen atölye başvurularını tutacak state


// 🛍️ ESER SİPARİŞLERİNİ ÇEKER
// 🛍️ 1. ESER SİPARİŞLERİNİ E-POSTA İLE ÇEKER (KUSURSUZ SÜRÜM)
const fetchSellerOrders = async () => {
  try {
    // 🚀 Doğrudan tıkır tıkır çalışan e-postayı çekiyoruz kanka:
    const userEmail = localStorage.getItem('userEmail'); 
    
    // Eski 'userId' kontrolünü ve o erkenden bitiren sinsi 'return'ü tamamen sildik!
    if (!userEmail) {
      console.warn("⚠️ Giriş yapmış kullanıcı e-postası bulunamadı!");
      return;
    }

    // Backend'e tertemiz e-posta parametresiyle gidiyoruz:
    const response = await axios.get(`http://localhost:8080/seller-orders?email=${userEmail}`);
    sellerOrders.value = response.data || [];
  } catch (error) {
    console.error("Siparişler çekilemedi:", error);
  }
};

// 🏫 2. ATÖLYE BAŞVURULARINI ÇEKER
// 🏫 ATÖLYE REZERVASYON BAŞVURULARINI E-POSTA İLE ÇEKER
const fetchWorkshopOrders = async () => {
  try {
    // Hafızadan gıcır gıcır çalışan aktif e-postayı çekiyoruz:
    const userEmail = localStorage.getItem('userEmail');
    
    if (!userEmail) {
      console.warn("⚠️ Giriş yapmış kullanıcı e-postası bulunamadı!");
      return;
    }

    // 🚀 URL'in sonuna sellerId=2 tuzağını silip email=${userEmail} dinamiğini çakıyoruz kanka!
    const response = await axios.get(`http://localhost:8080/seller-workshops-orders?email=${userEmail}`);
    workshopOrders.value = response.data || [];
  } catch (error) {
    console.error("Atölye başvuruları çekilemedi:", error);
  }
};

// Satış Onaylama Fonksiyonu
const confirmSale = async (purchaseId) => {
  if (!confirm('Bu satışı onaylamak istediğinize emin misiniz?')) return;
  
  try {
    await axios.post('http://localhost:8080/confirm-sale', { purchaseId });
    alert('Satış başarıyla onaylandı! ✅');
    await fetchSellerOrders(); // Tabloyu tazele
  } catch (error) {
    alert('Onay işlemi sırasında bir hata oluştu.');
  }
};



// Atölye Başvurusunu Onaylama Fonksiyonu
const confirmWorkshopOrder = async (enrollmentId) => {
  if (!confirm('Bu atölye başvurusunu onaylamak istediğinize emin misiniz?')) return;
  
  try {
    const response = await axios.post('http://localhost:8080/confirm-workshop-enrollment', {
      enrollmentId: parseInt(enrollmentId)
    });
    alert(response.data.message || 'Başvuru başarıyla onaylandı! ✅');
    await fetchWorkshopOrders(); // Tabloyu tazele
  } catch (error) {
    alert('Onay işlemi sırasında bir hata oluştu.');
    console.error(error);
  }
};

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

// 🚀 GÜNCEL: Düzenleme Modalını Açarken eski katılımcı sayısını hafızaya alır
const openEditModal = (enroll) => {
  selectedEnrollment.value = { ...enroll };
  originalParticipantCount.value = enroll.participantCount; // Eski sayıyı kilitledik kanka
  editModalPaymentMethod.value = 'Kredi Kartı'; // Her açılışta yöntemi sıfırla
  showEditModal.value = true;
};

// 🚀 GÜNCEL: Bakiye ve Ekstra Ödeme Denetimli Güncelleme Fonksiyonu
const updateEnrollment = async () => {
  try {
    const userEmail = localStorage.getItem('userEmail');
    
    // Katılımcı sayısı veritabanındaki halinden yüksekse ekstra ödeme kontrolü devreye girer
    if (selectedEnrollment.value.participantCount > originalParticipantCount.value) {
      const diff = selectedEnrollment.value.participantCount - originalParticipantCount.value;
      const extraCost = diff * selectedEnrollment.value.workshopPrice;
      
      // Bakiye kontrolü ön yüz güvenlik duvarı
      if (editModalPaymentMethod.value === 'Uygulama Bakiyesi' && user.value.balance < extraCost) {
        alert("Ekstra katılımcı ödemesi için uygulama bakiyeniz yetersiz! 💸 Lütfen bakiye yükleyin.");
        return;
      }
      
      if (!confirm(`Farktan dolayı ${extraCost} TL ekstra ödeme alınacaktır. Onaylıyor musunuz?`)) return;
    }

    // Go backend'deki o yeni transactionlı PUT endpointine verileri yolluyoruz
    await axios.put(`http://localhost:8080/update-enrollment`, {
      id: selectedEnrollment.value.id,
      participantCount: parseInt(selectedEnrollment.value.participantCount),
      reservedDate: selectedEnrollment.value.reservedDate,
      paymentMethod: editModalPaymentMethod.value, // Boş gitmesi engellendi kanka!
      email: userEmail
    });

    alert("Rezervasyon başarıyla güncellendi ve ödemeniz tescillendi! ✨");
    showEditModal.value = false;
    
    // Bakiye düşebileceği için profil objesini ve listeyi tazeleyelim
    const profileRes = await axios.get(`http://localhost:8080/profile?email=${userEmail}`);
    user.value = profileRes.data;
    fetchEnrollments();
  } catch (error) {
    alert(error.response?.data?.message || "Güncelleme başarısız oldu.");
  }
};

// Satın alınan eserleri getir
const fetchArtworkPurchases = async () => {
  const email = localStorage.getItem('userEmail');
  try {
    const res = await axios.get(`http://localhost:8080/user-purchases?email=${email}`);
    purchases.value = res.data || [];
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
    const res = await axios.get(`http://localhost:8080/artworks`);
    const allArtworks = res.data || [];
    
    const profileRes = await axios.get(`http://localhost:8080/profile?email=${userEmail}`);
    const userName = profileRes.data.firstName + ' ' + profileRes.data.lastName;
    
    myArtworks.value = allArtworks.filter(art => art.artist === userName);
  } catch (error) {
    console.error("Kendi eserler yüklenemedi:", error);
  }
};

// Eser sil
const deleteMyArtwork = async (artworkId) => {
  if (confirm("Bu eseri silmek istediğinize emin misiniz? This işlem geri alınamaz.")) {
    try {
      const userEmail = localStorage.getItem('userEmail');
      await axios.post(`http://localhost:8080/delete-artwork`, {
        artworkId: artworkId,
        email: userEmail
      });
      alert("Eser başarıyla silindi.");
      
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
      await fetchMyWorkshops(); 
    } catch (error) {
      console.error("Silme hatası:", error);
      alert("Silme işlemi başarısız: " + (error.response?.data || "Sunucu hatası"));
    }
  }
};

// Bakiye Yükleme Fonksiyonu
const addBalance = async () => {
  const input = prompt("Yüklemek istediğiniz tutarı girin (TL):", "500");
  if (input === null) return;

  const amount = parseFloat(input);
  if (isNaN(amount) || amount <= 0) {
    alert("⚠️ Geçerli bir sayı giriniz!");
    return;
  }

  try {
    const userEmail = localStorage.getItem('userEmail');
    await axios.put('http://localhost:8080/profile/update-balance', {
      email: userEmail,
      amount: amount
    });

    const profileRes = await axios.get(`http://localhost:8080/profile?email=${userEmail}`);
    user.value = profileRes.data;

    alert(`🎉 ${amount} TL başarıyla yüklendi!`);
  } catch (error) {
    console.error("Bakiye hatası:", error);
    alert("Bakiye güncellenirken bir sorun oluştu.");
  }
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

const formatDate = (d) => new Intl.DateTimeFormat('tr-TR', { day: 'numeric', month: 'long' }).format(new Date(d));

onMounted(() => {
  fetchUser();
  fetchFavorites();
  fetchPurchases();
  fetchComparisons();
});
</script>
