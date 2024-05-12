<template>
  <v-app>

  <div class="header">
    <div class="auth_buttons">
      <v-btn class="auth_button">Войти</v-btn>
      <v-btn class="auth_button">Регистрация</v-btn>
    </div>
  </div>

  <v-navigation-drawer
  app
  permanent
  :mini-variant="isMenuMinimize"
  >

  <v-list style="height: 100%" v-if="true" dense>

    <div style="height: 100%" id="drawer-menu-hover">
      <v-list-item :to="{ name: 'Tests' }">
        <v-list-item-title>{{ "Тесты" }}</v-list-item-title>
        </v-list-item>
    </div>
  </v-list>
  </v-navigation-drawer>

  <div>
    <router-view />
  </div>
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
</style>


<script>
import { ref, reactive } from "vue";
import {red, teal} from "vuetify/util/colors";
export default {
  computed: {
    red() {
      return red
    },
    teal() {
      return teal
    }
  },
  setup() {
    return {
      isMenuMinimize: ref(false),
      isMouseOnMenu: ref(false),
      easterEgg: ref(false),
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
}
</script>