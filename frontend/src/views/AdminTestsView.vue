<template>
  <v-table>
    <thead>
    <tr>
      <th class="text-left">
        Идентификатор
      </th>
      <th class="text-left">
        Идентификатор создателя
      </th>
      <th class="text-left">
        Тип
      </th>
      <th class="text-left">
        Название
      </th>
      <th class="text-left">
        Краткое описание
      </th>
      <th class="text-left">
        Длинное описание
      </th>
      <th class="text-left">
        Тэги
      </th>
    </tr>
    </thead>
    <tbody>
    <tr
        v-for="item in tests"
        :key="item.id"
    >
      <td>{{ item.id }}</td>
      <td>{{ item.creator_id }}</td>
      <td>{{ item.type }}</td>
      <td>{{ item.title }}</td>
      <td>{{ item.short_text }}</td>
      <td>{{ item.long_text }}</td>
      <td>{{ item.tags }}</td>
      <v-btn @click="deleteTest(item.id)" icon="mdi-trash-can"></v-btn>
    </tr>
    </tbody>
  </v-table>
</template>
<script>
import {useStore} from "@/store";

let store = useStore();

export default {
  computed: {
    tests() {
      return store.getters["tests/tests"]
    }
  },
  mounted() {
    store.dispatch("tests/getTests", true);
  },
  methods: {
    deleteTest(id) {
      store.dispatch("tests/deleteTest", id).then(() => {
        store.dispatch("tests/getTests");
      })
    }
  }
}
</script>