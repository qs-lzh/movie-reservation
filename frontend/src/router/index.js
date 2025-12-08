import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/HomeView.vue')
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue')
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/RegisterView.vue')
    },
    {
      path: '/movies/:id',
      name: 'movie-detail',
      component: () => import('../views/MovieDetailView.vue')
    },
    {
      path: '/showtimes/:id/reserve',
      name: 'seat-reservation',
      component: () => import('../views/SeatReservationView.vue')
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/ProfileView.vue')
    },
    {
      path: '/admin/add-movie',
      name: 'admin-add-movie',
      component: () => import('../views/AdminAddMovieView.vue')
    },
    {
      path: '/admin/update-movie/:id',
      name: 'admin-update-movie',
      component: () => import('../views/AdminUpdateMovieView.vue')
    },
  ]
})

export default router
