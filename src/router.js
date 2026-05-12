import { createRouter, createWebHistory } from 'vue-router'

import ProfileView from './views/ProfileView.vue';
import HomeView from './views/HomeView.vue'
import LoginView from './views/LoginView.vue'
import ArtworkDetail from './views/ArtworkDetail.vue'
import WorkshopView from './views/WorkshopView.vue'
import ArtworkCard from './components/ArtworkCard.vue'
import ArtistDetail from './views/ArtistDetail.vue';
import WorkshopDetail from './views/WorkshopDetail.vue';
import AddArtwork from './views/AddArtwork.vue'
import AddWorkshop from './views/AddWorkshop.vue'




const routes = [
  { 
    path: '/login', 
    name: 'Login', 
    component: LoginView 
  },
  {
    path: '/profile',
    name: 'Profile',
    component: ProfileView
  },
  { 
    path: '/register', 
    name: 'Register',
    component: () => import('./views/RegisterView.vue') 
  },
  { 
    path: '/', 
    name: 'Home',
    component: HomeView
  },
  { 
    path: '/artwork/:id', // ID bazlı dinamik rota
    name: 'ArtworkDetail',
    component: ArtworkDetail,
    props: true // URL'deki id'yi sayfaya veri olarak gönderir
  },
  { 
    path: '/workshops', 
    name: 'Workshops',
    component: WorkshopView 
  },
  {
  path: '/artist/:name',
  name: 'ArtistDetail',
  component: ArtistDetail,
  props: true
  },
  {
  path: '/workshops/:id',
  name: 'WorkshopDetail',
  component: WorkshopDetail,
  props: true
  },
  {
    path: '/add-artwork',
    name: 'AddArtwork',
    component: AddArtwork
  },
  { 
    path: '/add-workshop', 
    name: 'AddWorkshop', 
    component: AddWorkshop 
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('userToken'); // Tarayıcıdaki bilekliği kontrol et
  
  // Eğer kullanıcı giriş yapmamışsa ve 'login' veya 'register' dışında bir yere gitmeye çalışıyorsa
  if (!token && to.path !== '/login' && to.path !== '/register') {
    next('/login'); // Dur yolcu! Önce giriş yap.
  } else if (token && (to.path === '/login' || to.path === '/register')) {
    next('/'); // Zaten giriş yapmışsın, tekrar Login'e girmene gerek yok, ana sayfaya dön.
  } else {
    next(); // Her şey yolunda, devam et.
  }
});

export default router