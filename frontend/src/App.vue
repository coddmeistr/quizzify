<template>
  <div class="header">
    <div class="auth_buttons">
      <v-btn class="auth_button">Войти</v-btn>
      <v-btn class="auth_button">Регистрация</v-btn>
    </div>
  </div>
  <v-main>
    <router-view />
  </v-main>
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


<script setup lang="ts">
</script>