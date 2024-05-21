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
            :shortText = test.short_text
            :questionsCount = "test.questions.length">
        ></test-short>
      </v-col>
    </v-row>
  </div>
</template>

<script>
import TestShort from "@/components/TestShort.vue";
import {useStore} from "@/store";

let store = useStore();

export default {
  components: {TestShort},
  computed:{
    tests(){
      return store.getters["tests/tests"]
    }
  },

  mounted() {
    store.dispatch("tests/getTests")
  },

  methods: {
    displayFullTest(testId) {
      this.$router.push({name: 'TestFull', params: {testId: testId}})
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