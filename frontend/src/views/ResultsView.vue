<template>
  <v-table>
    <thead>
    <tr>
      <th class="text-left">
        Идентификатор
      </th>
      <th class="text-left">
        Идентификатор пользователя
      </th>
      <th class="text-left">
        Идентификатор теста
      </th>
      <th class="text-left">
        Процент правильности
      </th>
    </tr>
    </thead>
    <tbody>
    <tr
        v-for="item in accounts"
        :key="item.resultId"
    >
      <td>{{ item.resultId }}</td>
      <td>{{ item.userId }}</td>
      <td>{{ item.testId }}</td>
      <td>{{ item.percentage }}</td>
      <v-btn @click="deleteResult(item.resultId)" icon="mdi-trash-can"></v-btn>
    </tr>
    </tbody>
  </v-table>
</template>
<script>
import {useStore} from "@/store";

let store = useStore();

export default {
  computed: {
    results() {
      return store.getters["tests/results"]
    }
  },
  mounted() {
    store.dispatch("tests/getResults");
  },
  methods: {
    deleteResult(id) {
      store.dispatch("tests/deleteResult", id).then(() => {
        store.dispatch("tests/getResults");
      })
    }
  }
}
</script>