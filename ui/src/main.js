import {
  extend,
  ValidationObserver,
  ValidationProvider
} from "vee-validate";
import {
  messages
} from "vee-validate/dist/locale/id.json";
import * as rules from "vee-validate/dist/rules";
import Vue from 'vue'
import './plugins/axios'
import App from './App.vue'
import './registerServiceWorker'
import router from './router'
import store from './store'
import vuetify from './plugins/vuetify'

//import '@babel/polyfill'
import "core-js/stable";
import "regenerator-runtime/runtime";

Vue.component("Validate", ValidationProvider);
Vue.component("ValidateAll", ValidationObserver);

extend("fully_required", {
  validate: (value) => {
    return value!="";
  },
  message: (field) => {
    return `${field} harus diisi`;
  },
});

Vue.config.productionTip = false

new Vue({
  router,
  store,
  vuetify,
  render: h => h(App)
}).$mount('#app')
