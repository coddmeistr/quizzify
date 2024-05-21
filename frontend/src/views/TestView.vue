<template>
  <div>
    <TestFull v-if="test.id && test.id !== ''"
        :key="test.id"
        v-bind="test"
        @action="alert('action')"
        :id = "test.id"
        :title = "test.title"
        :shortText = test.short_text
        :longText = test.long_text
        :creatorName = "test.creator_id"
        :tags = "test.tags"
        :questions = "test.questions"
        :type = "test.type"
      />
    <div v-else>
      Тест не найден
    </div>
  </div>
</template>

<script>
import TestFull from "@/components/TestFull.vue";
import {useStore} from "@/store";

let store = useStore();

export default {
  components: {TestFull},

  mounted() {
    store.dispatch("tests/getTest", this.$route.params.testId)
  },
  computed:{
    test() {
      return store.getters["tests/test"]
    }
  },
}
</script>

<style lang="scss" scoped>

</style>