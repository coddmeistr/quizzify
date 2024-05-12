<template>
  <div class="tests-grid">
    <v-row v-for="(tests, index) in getTestsRows()" :key="index">
      <v-col v-for="test in tests" :key="test.id">
        <test-short
            :key="test.id"
            v-bind="test"
            @action="displayFullTest"
            :id = "test.id"
            :title = "test.title"
            :shortText = test.shortText
            :questionsCount = "test.questionsCount">
        ></test-short>
      </v-col>
    </v-row>
  </div>
</template>

<script>
import TestShort from "@/components/TestShort.vue";

export default {
  components: {TestShort},
  data() {
    return {
      tests: [
        {
          id: 1,
          title: "Тест 1",
          shortText: "Тест по сексу",
          questionsCount: 5
        },
        {
          id: 10,
          title: "Тест 2",
          shortText: "Тест по математике",
          questionsCount: 1
        },
        {
          id: 11,
          title: "Тест 3",
          shortText: "Тест по математике",
          questionsCount: 2
        },
        {
          id: 12,
          title: "Тест 4",
          shortText: "Тест по математике",
          questionsCount: 3
        }
      ]
    }
  },

  methods: {
    displayFullTest(testId) {
      alert(testId)
    },
    getTestsRows(){
      const inOneRow = 3
      let rows = []
      for (let i = 0; i < this.tests.length; i += inOneRow) {
        rows.push(this.tests.slice(i, Math.min(this.tests.length, i + inOneRow)));
      }
      return rows
    }
  }
}
</script>

<style scoped>
.tests-grid{
  margin-top: 20px;
  margin-left: auto;
  margin-right: auto;
  max-width: 50%;
}
</style>