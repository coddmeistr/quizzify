<template>
  <v-app>

  <div class="header">

    <div v-if="!isLoggedIn" class="auth_buttons">
      <router-link to="/login">
        <v-btn class="auth_button">Войти</v-btn>
      </router-link>
      <router-link to="/register">
        <v-btn class="auth_button">Регистрация</v-btn>
      </router-link>
    </div>
    <div v-else class="authenticated">
      <span>{{userData.login}}</span>
      <v-btn class="auth_button" @click="logout">Выйти</v-btn>
    </div>
  </div>

  <v-navigation-drawer
  app
  permanent
  fixed
  :rail="isMenuMinimize"
  >

    <router-link to="/">
    <div class="d-flex gg-15px align-center justify-center pa-5">
      <template v-if="!config.logoSrc">
        <v-img
            alt=""
            src="/assets/logo.svg"
            max-height="60px"
            max-width="60px"
            contain
        ></v-img>
      </template>
      <template v-else>
        <v-img
            v-if="!isMenuMinimize"
            transition="fade-transition"
            alt="logo"
            src="/assets/logo.svg"
            contain
        ></v-img>
      </template>
    </div>
    </router-link>

  <v-list style="height: 100%" v-if="true" dense>

    <div
        :class="{
            'd-flex': true,
            'align-center': true,
            'justify-space-between': !isMenuMinimize,
            'flex-column-reverse': isMenuMinimize,
          }"
    >
      <v-list-subheader v-if="!isMenuMinimize">Основное</v-list-subheader>
      <v-btn variant="text" @click="isMenuMinimize = !isMenuMinimize" icon="true">
        <v-icon v-if="isMenuMinimize">mdi-arrow-right</v-icon>
        <v-icon v-else>mdi-arrow-left</v-icon>
      </v-btn>
    </div>

    <div style="height: 100%" id="drawer-menu-hover">
      <v-list-item prepend-icon="mdi-view-dashboard-variant" :to="{ name: 'Tests' }" title="Тесты" />
      <v-list-item prepend-icon="mdi-view-dashboard-variant" :to="{ name: 'Results' }" title="Результаты" />
      <v-list-item prepend-icon="mdi-plus" :to="{ name: 'TestCreating' }" title="Создать тест" />

      <v-list-subheader v-if="!isMenuMinimize && userData?.permissions?.includes(2)">Администрирование</v-list-subheader>
      <v-list-item prepend-icon="mdi-view-dashboard-variant" :to="{ name: 'Accounts' }"
                   v-if="userData?.permissions?.includes(2)"
                   title="Аккаунты" />
      <v-list-item prepend-icon="mdi-view-dashboard-variant" :to="{ name: 'TestsAdmin' }"
                   v-if="userData?.permissions?.includes(2)"
                   title="Тесты" />
      <v-list-item prepend-icon="mdi-view-dashboard-variant" :to="{ name: 'Results' }"
                   v-if="userData?.permissions?.includes(2)"
                   title="Результаты" />
    </div>
  </v-list>
  </v-navigation-drawer>

  <v-main>
    <router-view />
  </v-main>
  </v-app>
</template>

<style>
.header {
  background-color: lightsteelblue;
  min-height: 100px;
  display: flex;
  justify-content: flex-end;
  align-items: center;
}

.auth_buttons {
  margin-right: 50px;
}

.auth_button {
  margin-left: 10px;
}

.authenticated{
  margin-right: 50px;
}
</style>


<script>
import { ref, reactive } from "vue";
import {red, teal} from "vuetify/util/colors";
import config from "@/config"
import { useStore } from "@/store";

let store = useStore();
export default {
  computed: {
    red() {
      return red
    },
    teal() {
      return teal
    },
    isLoggedIn() {
      return store.getters["auth/isLoggedIn"]
    },
    userData(){
      return store.getters["auth/userdata"]
    }
  },
  setup() {
    return {
      isMenuMinimize: ref(false),
      isMouseOnMenu: ref(false),
      easterEgg: ref(false),
      config,
      overlay: reactive({
        timeoutId: null,
        isVisible: true,
        buttonTitle: "",
        uuid: "",
        x: 0,
        y: 0,
      }),
    };
  },
  methods: {
    configureHoverOnMenu() {
      document
          .getElementById("drawer-menu-hover")
          ?.addEventListener("mouseenter", () => {
            if (this.isMenuMinimize) {
              this.isMouseOnMenu = true;
              setTimeout(() => {
                if (this.isMouseOnMenu) {
                  this.isMenuMinimize = false;
                }
              }, 5000);
            }
          });
      document
          .getElementById("drawer-menu-hover")
          ?.addEventListener("mouseleave", () => {
            if (this.isMouseOnMenu) {
              this.isMouseOnMenu = false;
              setTimeout(() => {
                this.isMenuMinimize = true;
              }, 1000);
            }
          });
    },
    logout(){
      store.dispatch("auth/logout");
    }
  },
  mounted() {
    this.configureHoverOnMenu();
    store.dispatch("auth/fetchUserData", { token: store.getters["auth/token"]});
  },
}
</script>