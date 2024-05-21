<template>
  <div v-if="userChoices">
    <Question  v-for="(q) in test.questions" :key="q.id" v-bind="q"
              :id = "q.id"
              :variants = "q.variants"
              :shortText = q.short_text
              :longText = q.long_text
              :type = "q.type"
              :required = "q.required"
              :points = "q.points"
              v-model:selectedVariant="userChoices[`${q.id}`].selectedVariant"
              v-model:selectedVariants="userChoices[`${q.id}`].selectedVariants"
              v-model:writedText="userChoices[`${q.id}`].writedText"
    />
    <v-btn class="mt-5 end-btn" @click="sendResults">ЗАКОНЧИТЬ ТЕСТ</v-btn>
  </div>
</template>

<script>
import {useStore} from "@/store";
import Question from "@/components/TestQuestion.vue";

let store = useStore();

export default {
  components: {Question},

  data() {
    return {
      userChoices: null
    }
  },

  async mounted() {
    await store.dispatch("tests/getTest", this.$route.params.testId)
    this.userChoices = {}
    this.test.questions.forEach((q) => {
      this.userChoices[`${q.id}`] = {
        selectedVariant: null,
        selectedVariants: [],
        writedText: null
      }
    })
  },
  computed:{
    test() {
      return store.getters["tests/test"]
    },
  },
  methods: {
    sendResults() {
      let answers = []
      for (let ind in this.test.questions) {
        let q = this.test.questions[ind]
        if (!this.userChoices[`${q.id}`]){
          continue
        }

        let variants = this.userChoices[`${q.id}`].selectedVariants
        if (variants.length === 0){
          variants = null
        }

        answers.push({
          question_id: q.id,
          chosen_id: this.userChoices[`${q.id}`].selectedVariant,
          chosen_ids: variants,
          writed_text: this.userChoices[`${q.id}`].writedText
        })
      }

      store.dispatch("tests/sendResult", {
        test_id: this.test.id,
        answers: answers
      })
    }
  }
}
</script>

<style scoped>

.container{
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 100px;
}

.end-btn{
  margin-left: 2%;
  margin-bottom: 50px;
}

</style>