<template>
  <v-table>
    <thead>
    <tr>
      <th class="text-left">
        Имя пользователя
      </th>
      <th class="text-left">
        Электронная почта
      </th>
      <th class="text-left">
        Доступы
      </th>
      <th class="text-left">
        Идентификатор
      </th>
    </tr>
    </thead>
    <tbody>
    <tr
        v-for="item in accounts"
        :key="item.login"
    >
      <td>{{ item.login }}</td>
      <td>{{ item.email }}</td>
      <td>{{ item.permissions }}</td>
      <td>{{ item.userId }}</td>
      <v-btn @click="deleteUser(item.userId)" icon="mdi-trash-can"></v-btn>
    </tr>
    </tbody>
  </v-table>
</template>
<script>
import {useStore} from "@/store";

let store = useStore();

export default {
  computed: {
    accounts() {
      return store.getters["auth/accounts"]
    }
  },
  mounted() {
    store.dispatch("auth/accountsList");
  },
  methods: {
    deleteUser(id) {
      if (store.getters["auth/userdata"].userId === id){
        alert("Нельзя удалить себя");
        return
      }
      store.dispatch("auth/deleteAccount", id).then(() => {
        store.dispatch("auth/accountsList");
      })
    }
  }
}
</script>