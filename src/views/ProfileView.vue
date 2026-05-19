<template>
  <div class="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div v-if="user" class="max-w-4xl mx-auto">
      
      <div class="bg-white rounded-3xl shadow-xl overflow-hidden mb-8 border border-gray-100">
        <div class="h-32 bg-gradient-to-r from-blue-600 to-indigo-700"></div>
        <div class="px-8 pb-8">
          <div class="relative flex justify-between items-end -mt-12 mb-6">
            <div class="w-32 h-32 bg-white p-2 rounded-3xl shadow-lg">
              <div class="w-full h-full bg-blue-100 rounded-2xl flex items-center justify-center text-4xl">
                {{ user.firstName ? user.firstName.charAt(0) : 'U' }}
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
            <p class="text-gray-500 font-medium">{{ user.email }} • <span class="text-blue-600">{{ user.userRole || user.role }}</span></p>
            <p v-if="user.biography" class="mt-4 text-gray-600 leading-relaxed max-w-2xl">{{ user.biography }}</p>
          </div>
        </div>
      </div>

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

      <div class="space-y-6">
        
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

        <div class="mt-8 space-y-12">
          <div class="bg-blue-50 border-2 border-blue-100 rounded-2xl p-6 flex items-center justify-between">
            <div class="flex items-center space-x-4">
              <div class="p-3 bg-blue-600 rounded-xl text-white shadow-lg text-2xl">💰</div>
              <div>
                <p class="text-sm font-bold text-blue-600 uppercase tracking-wider">Hesap Bakiyeniz</p>
                <h3 class="text-3xl font-black text-blue-900">
                  {{ user?.balance ? user.balance.toLocaleString() : '0' }} TL
                </h3>
              </div>
            </div>
            <button @click="addBalance" class="px-6 py-3 bg-white border-2 border-blue-600 text-blue-600 font-bold rounded-xl hover:bg-blue-600 hover:text-white transition-all shadow-sm">
              + Bakiye Yükle
            </button>
          </div>

          <div class="border-t-2 border-gray-100 pt-8">
            <h3 class="text-2xl font-bold text-gray-800 mb-6 flex items-center">🛍️ Satın Aldığım Eserler Detay</h3>
            <div v-if="boughtArtworks.length > 0" class="space-y-4">
              <div v-for="item in boughtArtworks" :key="item.id" class="bg-white rounded-2xl p-5 border border-gray-100 shadow-sm flex items-center justify-between">
                <div class="flex items-center space-x-4">
                  <img :src="item.image" class="w-20 h-20 object-cover rounded-xl shadow-sm">
                  <div>
                    <h4 class="text-lg font-bold text-gray-900">{{ item.title }}</h4>
                    <p class="text-blue-600 font-semibold">{{ item.price }} TL</p>
                    <p class="text-xs text-gray-400">{{ new Date(item.date).toLocaleDateString('tr-TR') }}</p>
                  </div>
                </div>
                <div class="text-right">
                  <span :class="{
                      'bg-amber-100 text-amber-700': item.status === 'Hazırlanıyor',
                      'bg-blue-100 text-blue-700': item.status === 'Kargoda',
                      'bg-green-100 text-green-700': item.status === 'Teslim Edildi'
                    }" class="px-4 py-1.5 rounded-full text-xs font-black uppercase tracking-widest">
                    {{ item.status }}
                  </span>
                </div>
              </div>
            </div>
            <div v-else class="text-center py-10 bg-gray-50 rounded-2xl border-2 border-dashed border-gray-200">
              <p class="text-gray-400 italic">Henüz bir eser satın almadınız.</p>
            </div>
          </div>

          <div class="pt-8 border-t-2 border-gray-100">
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
                  <span :class="enroll.status === 'Onaylandı' ? 'bg-green-100 text-green-700 border-green-200' : 'bg-amber-100 text-amber-700 border-amber-200'" class="px-4 py-1.5 rounded-full text-xs font-black uppercase tracking-widest border">
                    {{ enroll.status }}
                  </span>
                </div>
                <div class="flex gap-2">
                  <button v-if="enroll.status !== 'Onaylandı'" @click="openEditModal(enroll)" class="px-4 py-2 bg-blue-100 text-blue-700 font-bold rounded-lg hover:bg-blue-200 transition-colors">Düzenle</button>
                  <button v-if="enroll.status !== 'Onaylandı'" @click="cancelEnrollment(enroll.id)" class="px-4 py-2 bg-red-100 text-red-700 font-bold rounded-lg hover:bg-red-200 transition-colors">İptal Et</button>
                  <span v-else class="text-green-600 font-bold text-sm flex items-center gap-1 bg-green-50 px-4 py-2 rounded-xl border border-green-100">✨ Rezervasyon Kesinleşti</span>
                </div>
              </div>
            </div>
            <div v-else class="text-center py-10 bg-gray-50 rounded-2xl border-2 border-dashed border-gray-200">
              <p class="text-gray-400 italic">Henüz bir atölye rezervasyonunuz bulunmuyor.</p>
            </div>
          </div>

          <div v-if="userRole === 'Instructor' || user?.userRole === 'Instructor'">
            
            <div class="pt-8 border-t-2 border-gray-100 mt-12">
              <h3 class="text-2xl font-bold text-gray-800 mb-8 flex items-center">🎨 Gelen Siparişler (Satış Onayı)</h3>
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
                        <span :class="order.status === 'Onaylandı' ? 'bg-green-100 text-green-700' : 'bg-orange-100 text-orange-700'" class="px-3 py-1 rounded-full text-xs font-bold">
                          {{ order.status }}
                        </span>
                      </td>
                      <td class="px-6 py-4 text-right">
                        <button v-if="order.status === 'Hazırlanıyor'" @click="confirmSale(order.id)" class="bg-blue-600 text-white px-4 py-2 rounded-xl text-sm font-bold hover:bg-blue-700 shadow-lg shadow-blue-100 transition-all">Onayla</button>
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

            <div class="pt-8 border-t-2 border-gray-100 mt-12">
              <h3 class="text-2xl font-bold text-gray-800 mb-8 flex items-center">🎫 Gelen Atölye Başvuruları (Rezervasyon Onayı)</h3>
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
                        <span :class="order.status === 'Onaylandı' ? 'bg-green-100 text-green-700' : 'bg-amber-100 text-amber-700'" class="px-3 py-1 rounded-full text-xs font-bold">
                          {{ order.status }}
                        </span>
                      </td>
                      <td class="px-6 py-4 text-right">
                        <button v-if="order.status !== 'Onaylandı'" @click="confirmWorkshopOrder(order.id)" class="bg-indigo-600 text-white px-4 py-2 rounded-xl text-sm font-bold hover:bg-indigo-700 shadow-lg shadow-indigo-100 transition-all">Başvuruyu Onayla ✅</button>
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

            <div class="pt-8 border-t-2 border-gray-100 mt-12">
              <h3 class="text-2xl font-bold text-gray-800 mb-8 flex items-center">M12 Eklediğim Eserler</h3>
              <div v-if="myArtworks && myArtworks.length > 0" class="space-y-4">
                <div v-for="artwork in myArtworks" :key="artwork.id" class="bg-white rounded-2xl p-5 border border-gray-100 shadow-sm flex items-center justify-between hover:shadow-md transition-shadow">
                  <div class="flex items-center space-x-4 flex-1">
                    <img :src="artwork.imageUrl" class="w-20 h-20 object-cover rounded-xl shadow-sm" :alt="artwork.title">
                    <div>
                      <h4 class="text-lg font-bold text-gray-900">{{ artwork.title }}</h4>
                      <p class="text-gray-600 text-sm">{{ artwork.description?.substring(0, 50) }}...</p>
                      <p class="text-green-600 font-semibold mt-1">{{ artwork.price }} TL</p>
                    </div>
                  </div>
                  <button @click="deleteMyArtwork(artwork.id)" class="ml-4 px-4 py-2 bg-red-100 text-red-600 font-bold rounded-lg hover:bg-red-200 transition-colors flex items-center gap-2">🗑️ Sil</button>
                </div>
              </div>
              <div v-else class="text-center py-10 bg-gray-50 rounded-2xl border-2 border-dashed border-gray-200">
                <p class="text-gray-400 italic">Henüz eser eklemediniz.</p>
              </div>
            </div>

            <div class="pt-8 border-t-2 border-gray-100 mt-12">
              <h3 class="text-2xl font-bold text-gray-800 mb-8 flex items-center">Oluşturduğum Atölyeler</h3>
              <div v-if="myWorkshops && myWorkshops.length > 0" class="space-y-4">
                <div v-for="ws in myWorkshops" :key="ws.id" class="bg-white rounded-2xl p-5 border border-gray-100 shadow-sm flex items-center justify-between hover:shadow-md transition-shadow">
                  <div class="flex items-center space-x-4 flex-1">
                    <img :src="ws.imageUrl" class="w-20 h-20 object-cover rounded-xl shadow-sm" :alt="ws.title">
                    <div>
                      <h4 class="text-lg font-bold text-gray-900">{{ ws.title }}</h4>
                      <p class="text-gray-500 text-sm">📍 {{ ws.location }} | 📅 {{ ws.availableDates }}</p>
                      <p class="text-indigo-600 font-semibold mt-1">{{ ws.price }} TL <span class="text-gray-400 text-xs">(Kapasite: {{ ws.capacity }})</span></p>
                    </div>
                  </div>
                  <div class="flex gap-2">
                    <button @click="deleteMyWorkshop(ws.id)" class="px-4 py-2 bg-red-50 text-red-600 font-bold rounded-lg hover:bg-red-100 transition-colors">🗑️ Sil</button>
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
              ⚠️ Katılımcı sayısını artırdığınız için <strong>{{ (selectedEnrollment.participantCount - originalParticipantCount) * (selectedEnrollment.workshopPrice || 0) }} TL</strong> ekstra ödeme yapmanız gerekmektedir.
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
import { ref, onMounted, computed, watch, nextTick } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';

const router = useRouter();

// State Tanımlamaları (Çökmeleri Engelleyen Kısım)
const user = ref(null);
const userRole = ref('User');
const favorites = ref([]);
const purchases = ref([]);
const comparisons = ref([]);
const enrollments = ref([]);
const boughtArtworks = ref([]);
const myArtworks = ref([]);
const myWorkshops = ref([]);
const sellerOrders = ref([]);
const workshopOrders = ref([]);

const activeTab = ref('favorites');
const editingId = ref(null);
const editTitle = ref('');
const showEditModal = ref(false);
const selectedEnrollment = ref(null);
const originalParticipantCount = ref(0);
const editModalPaymentMethod = ref('Kredi Kartı');

const editUser = ref({ firstName: '', lastName: '', biography: '' });
const passwords = ref({ oldPassword: '', newPassword: '' });

const tabs = [
  { id: 'favorites', label: 'Favorilerim' },
  { id: 'orders', label: 'Siparişlerim' },
  { id: 'comparisons', label: 'Analizlerim' },
  { id: 'edit', label: 'Düzenle' }
];

// 🚀 USER PROFILE DETAYLARINI ÇEKER
const fetchUser = async () => {
  const email = localStorage.getItem('userEmail');
  if (!email) return;
  try {
    const res = await axios.get(`http://localhost:8080/profile?email=${email}`);
    user.value = res.data;
    userRole.value = res.data.userRole || res.data.role || 'User';
    editUser.value = { ...res.data };
    
    // Eğer kullanıcı eğitmense kendi içeriklerini de çekelim
    if (userRole.value === 'Instructor') {
      await fetchUserArtworks(email, `${res.data.firstName} ${res.data.lastName}`);
      await fetchMyWorkshops();
      await fetchSellerOrders();
      await fetchWorkshopOrders();
    }
  } catch (e) { 
    console.error("Kullanıcı bilgisi çekilemedi:", e); 
  }
};

// 🛍️ KULLANICININ SATIN ALDIĞI ESERLERİ ÇEKER
const fetchPurchases = async () => {
  const email = localStorage.getItem('userEmail');
  try {
    const res = await axios.get(`http://localhost:8080/user-purchases?email=${email}`);
    purchases.value = res.data || [];
    boughtArtworks.value = res.data || []; // Template'deki iki dizi de besleniyor
  } catch (e) { 
    console.error("Satın alımlar çekilemedi:", e); 
  }
};

// ❤️ FAVORİLERİ ÇEKER
const fetchFavorites = async () => {
  const email = localStorage.getItem('userEmail');
  try {
    const res = await axios.get(`http://localhost:8080/favorites/list?email=${email}`);
    favorites.value = res.data || [];
  } catch (e) { console.error(e); }
};

// 📊 KIYASLAMALARI GETİRİR
const fetchComparisons = async () => {
  const token = localStorage.getItem('userToken');
  if (!token) return;
  try {
    const res = await axios.get('http://localhost:8080/comparisons', { headers: { Authorization: `Bearer ${token}` } });
    comparisons.value = res.data || [];
  } catch (e) { 
    console.log("Kıyaslamalar yüklenemedi (veya simüle ediliyor)");
  }
};

// 📅 ATÖLYE REZERVASYONLARINI GETİRİR
const fetchEnrollments = async () => {
  const userEmail = localStorage.getItem('userEmail');
  try {
    const res = await axios.get(`http://localhost:8080/user-enrollments?email=${userEmail}`);
    enrollments.value = res.data || [];
  } catch (error) {
    console.error("Rezervasyonlar yüklenemedi:", error);
  }
};

// 🎨 SANATÇININ KENDİ EKLEDİĞİ ESERLERİ FİLTRELEME
const fetchUserArtworks = async (userEmail, userName) => {
  try {
    const res = await axios.get(`http://localhost:8080/artworks`);
    const allArtworks = res.data || [];
    myArtworks.value = allArtworks.filter(art => art.artist === userName || art.artistID === user.value?.id);
  } catch (e) {
    console.error("Sanatçı eserleri filtreleme hatası:", e);
  }
};

// 🏫 EĞİTMENİN KENDİ ATÖLYELERİNİ ÇEKER
const fetchMyWorkshops = async () => {
  const email = localStorage.getItem('userEmail');
  try {
    const res = await axios.get(`http://localhost:8080/workshops?email=${email}`);
    myWorkshops.value = res.data || [];
  } catch (e) { 
    console.error("Atölyeler çekilemedi:", e); 
  }
};

// 🛍️ GELEN SİPARİŞLER (SATIŞ ONAYI)
const fetchSellerOrders = async () => {
  try {
    const userEmail = localStorage.getItem('userEmail'); 
    if (!userEmail) return;
    const response = await axios.get(`http://localhost:8080/seller-orders?email=${userEmail}`);
    sellerOrders.value = response.data || [];
  } catch (error) {
    console.error("Siparişler çekilemedi:", error);
  }
};

// 🏫 GELEN ATÖLYE BAŞVURULARI
const fetchWorkshopOrders = async () => {
  try {
    const userEmail = localStorage.getItem('userEmail');
    if (!userEmail) return;
    const response = await axios.get(`http://localhost:8080/seller-workshops-orders?email=${userEmail}`);
    workshopOrders.value = response.data || [];
  } catch (error) {
    console.error("Atölye başvuruları çekilemedi:", error);
  }
};

// SATIŞ ONAYLAMA
const confirmSale = async (purchaseId) => {
  if (!confirm('Bu satışı onaylamak istediğinize emin misiniz?')) return;
  try {
    await axios.post('http://localhost:8080/confirm-sale', { purchaseId });
    alert('Satış başarıyla onaylandı! ✅');
    await fetchSellerOrders();
  } catch (error) {
    alert('Onay işlemi sırasında bir hata oluştu.');
  }
};

// ATÖLYE BAŞVURUSU ONAYLAMA
const confirmWorkshopOrder = async (enrollmentId) => {
  if (!confirm('Bu atölye başvurusunu onaylamak istediğinize emin misiniz?')) return;
  try {
    const response = await axios.post('http://localhost:8080/confirm-workshop-enrollment', {
      enrollmentId: parseInt(enrollmentId)
    });
    alert(response.data.message || 'Başvuru başarıyla onaylandı! ✅');
    await fetchWorkshopOrders();
  } catch (error) {
    alert('Onay işlemi sırasında bir hata oluştu.');
  }
};

// PROFİL GÜNCELLEME
const updateProfile = async () => {
  try {
    await axios.post('http://localhost:8080/profile/update', { ...editUser.value, email: user.value.email });
    alert("Profil güncellendi!");
    await fetchUser();
  } catch (e) { alert("Hata oluştu."); }
};

// REZERVASYON İPTAL ETME
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

// MODAL AÇMA VE GÜNCELLEME İŞLEMLERİ
const openEditModal = (enroll) => {
  selectedEnrollment.value = { ...enroll };
  originalParticipantCount.value = enroll.participantCount;
  editModalPaymentMethod.value = 'Kredi Kartı';
  showEditModal.value = true;
};

const updateEnrollment = async () => {
  try {
    const userEmail = localStorage.getItem('userEmail');
    if (selectedEnrollment.value.participantCount > originalParticipantCount.value) {
      const diff = selectedEnrollment.value.participantCount - originalParticipantCount.value;
      const extraCost = diff * (selectedEnrollment.value.workshopPrice || 0);
      
      if (editModalPaymentMethod.value === 'Uygulama Bakiyesi' && user.value.balance < extraCost) {
        alert("Ekstra katılımcı ödemesi için bakiye yetersiz! 💸");
        return;
      }
      if (!confirm(`Farktan dolayı ${extraCost} TL ekstra ödeme alınacaktır. Onaylıyor musunuz?`)) return;
    }

    await axios.put(`http://localhost:8080/update-enrollment`, {
      id: selectedEnrollment.value.id,
      participantCount: parseInt(selectedEnrollment.value.participantCount),
      reservedDate: selectedEnrollment.value.reservedDate,
      paymentMethod: editModalPaymentMethod.value,
      email: userEmail
    });

    alert("Rezervasyon başarıyla güncellendi! ✨");
    showEditModal.value = false;
    await fetchUser();
    fetchEnrollments();
  } catch (error) {
    alert("Güncelleme başarısız oldu.");
  }
};

// ESER SİLME
const deleteMyArtwork = async (artworkId) => {
  if (confirm("Bu eseri silmek istediğinize emin misiniz?")) {
    try {
      const userEmail = localStorage.getItem('userEmail');
      await axios.post(`http://localhost:8080/delete-artwork`, { artworkId, email: userEmail });
      alert("Eser başarıyla silindi.");
      await fetchUser();
    } catch (error) { alert("Eser silinemedi."); }
  }
};

// ATÖLYE SİLME
const deleteMyWorkshop = async (id) => {
  if (confirm("Bu atölyeyi silmek istediğinize emin misiniz?")) {
    try {
      const userEmail = localStorage.getItem('userEmail');
      await axios.delete(`http://localhost:8080/delete-workshop?id=${id}`, {
        data: { workshopId: id, email: userEmail }
      });
      alert("Atölye başarıyla silindi.");
      await fetchUser();
    } catch (error) { alert("Silme işlemi başarısız."); }
  }
};

// BAKİYE YÜKLEME
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
    await axios.put('http://localhost:8080/profile/update-balance', { email: userEmail, amount });
    alert(`🎉 ${amount} TL başarıyla yüklendi!`);
    await fetchUser();
  } catch (error) { alert("Bakiye güncellenemedi."); }
};

// DİĞER YARDIMCI METOTLAR
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

const formatDate = (d) => {
  if (!d) return '-';
  return new Intl.DateTimeFormat('tr-TR', { day: 'numeric', month: 'long' }).format(new Date(d));
};

// SAYFA AÇILDIĞINDA TETİKLENEN ANA AKIŞ
onMounted(async () => {
  await fetchUser();
  await fetchFavorites();
  await fetchPurchases();
  await fetchComparisons();
  await fetchEnrollments();
});
</script>