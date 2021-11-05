import Vue from 'vue'
import VueRouter, { RouteConfig } from 'vue-router'
import Login from '../components/Login.vue'
import AuthInProcess from '../components/AuthInProcess.vue'
import FeedbackCmp from '../components/Feedback.vue'
import SendFeedback from '../components/SendFeedback.vue'
import * as global from '../common/GlobalConstants'

Vue.use(VueRouter)

const routes: Array<RouteConfig> = [
  {
    path: '/',
    name: 'Login',
    component: Login
  },
  {
    path: global.REDIRECT_PATH,
    name: 'AuthInProcess',
    component: AuthInProcess
  },
  {
    path: "/feedback",
    name: 'Feedback',
    component: FeedbackCmp
  },
  {
    path: "/feedback-form",
    name: 'FeedbackForm',
    component: SendFeedback
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
