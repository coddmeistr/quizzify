<template>
  <div class="container">
     <v-text-field class="field" v-model="title" label="Название теста"></v-text-field>
     <v-text-field class="field" v-model="short_text" label="Краткое описание(показывается в шапке)"></v-text-field>
     <v-text-field class="field" v-model="long_text" label="Полное описание"></v-text-field>
    <h2><br>Вопросы</h2>
     <TestQuestionCreating class="question" v-for="q in questions" :key="q.id"
                           v-model:id="questions[q.id].id"
                           v-model:shortText="questions[q.id].short_text"
                           v-model:longText="questions[q.id].long_text"
                           v-model:required="questions[q.id].required"
                           v-model:type="questions[q.id].type"
                           v-model:points="questions[q.id].points"
                           v-model:answers="questions[q.id].answers"
                           v-model:variants="questions[q.id].variants"
                           @delete="deleteQuestion"
     />
    <v-btn variant="text" @click="createQuestion" icon="true">
      <v-icon>mdi-plus</v-icon>
    </v-btn>
    <v-btn variant="tonal" @click="createTest">
      СОЗДАТЬ ТЕСТ
    </v-btn>
  </div>
</template>


<script>
import TestQuestionCreating from '@/components/TestQuestionCreating.vue'
import {useStore} from '@/store'

let store = useStore();

export default {
  name: "test-creating-view",
  components: {
        TestQuestionCreating
    },
  data(){
        return {
            currentQuestionCount: 1,
            type: '',
            title: '',
            short_text: '',
            long_text: '',
            creator_id: 0,
            tags: [],
            questions: {}
        }
    },
    methods:{
      createQuestion(){
        this.questions[this.currentQuestionCount] = this.questionSample(this.currentQuestionCount)
        this.currentQuestionCount++
      },
      deleteQuestion(id){
        delete this.questions[id]
      },
      questionSample(id){
        return {
          id: id,
          short_text: '',
          long_text: '',
          required: true,
          points: 10,
          answers: {},
          variants: {single_choice: {fields: [
                {text: 'Вариант 1', id: 0}, {text: 'Вариант 2', id: 1}, {text: 'Вариант 3', id: 2}
              ]}, multiple_choice: {max: 2,fields: [{text: 'Вариант 1', id: 0}, {text: 'Вариант 2', id: 1}]}}
        }
      },
      createTest(){
        let test = {
          type: "strict_test",
          title: this.title,
          short_text: this.short_text,
          long_text: this.long_text,
          tags: this.tags,
          main_image: {
            name: "main_image_stub",
            content: [1, 1]
          }
        }

        let questions = []
        Object.keys(this.questions).forEach(id => {
          let q = this.questions[id]
          if (q.points && q.points !== ""){
            q.points = parseInt(q.points)
          }
          questions.push(q)
        });

        test.questions = questions
        store.dispatch('tests/createTest', {test: test})
      }
    }

}
</script>

<style scoped>

.container{
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 30px;
  flex-direction: column;
  margin-bottom: 100px;
}
.question{
  width: 50%;
  margin-top: 20px;
}
.field{
  min-width: 50%;
}
</style>