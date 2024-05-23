<template>
  <v-card>
    <v-card-title>
      <v-text-field
          v-model="shortText"
          label='Текст вопроса'
          required
      ></v-text-field>
    </v-card-title>
    <v-card-text>
      <v-row>
        <v-col>
          <v-text-field
              v-model="longText"
              label='Дополнительный текст вопроса'
              required
          ></v-text-field>
        </v-col>

        <v-col>
          <v-select
              :items="types"
              label="Тип вопроса"
              @update:modelValue="updateVariants"
          ></v-select>
        </v-col>
      </v-row>

      <v-row v-if="type === 'single_choice'">
        <v-col>
          <v-radio-group v-model="answers.correct_id">
            <div class="variant" v-for="(variant, index) in variants.single_choice.fields" :key="index">
              <v-radio
                  class="item-btn"
                  :value="variant.id"
              ></v-radio>
              <v-responsive class="ma-0 pa-0" :width="`${variant.text.length}.5rem`">
                <v-text-field class="item" v-model="variant.text" rounded variant="outlined"></v-text-field>
              </v-responsive>
              <v-btn class="delete-btn-variants" @click="this.variants.single_choice.fields.splice(index, 1)">
                <v-icon>mdi-delete</v-icon>
              </v-btn>
            </div>
            <v-btn class="create-variant" @click="this.variants.single_choice.fields.push({id: this.variants.single_choice.fields.length, text: ''})">
              <v-icon>mdi-plus</v-icon>
            </v-btn>
          </v-radio-group>
        </v-col>
      </v-row>

      <v-row v-if="type === 'multiple_choice'">
        <v-col>
          <div class="variant" v-for="(variant, index) in variants.multiple_choice.fields" :key="index">
            <v-checkbox
                class="item-btn"
                v-model="answers.correct_ids"
                :value="variant.id"
            ></v-checkbox>
            <v-responsive class="ma-0 pa-0" :width="`${variant.text.length}.5rem`">
              <v-text-field class="item" v-model="variant.text" rounded variant="outlined"></v-text-field>
            </v-responsive>
            <v-btn class="delete-btn-variants" @click="deleteVariantMultipleChoice(index)">
              <v-icon>mdi-delete</v-icon>
            </v-btn>
          </div>
          <v-btn class="create-variant" @click="createVariantMultipleChoice">
            <v-icon>mdi-plus</v-icon>
          </v-btn>
        </v-col>
      </v-row>

      <v-row v-if="type === 'manual_input'">
        <v-col>
          <v-text-field
              v-model="answers.correct_text"
              label="Правильный ответ"
          ></v-text-field>
        </v-col>
      </v-row>

      <v-row>
        <v-col>
          <v-text-field
              v-model="points"
              label='Количество баллов за правильный ответ'
              type="number"
          ></v-text-field>
        </v-col>
      </v-row>
    </v-card-text>
    <div class="bot-row">
      <v-checkbox
          class="required"
          v-model="required"
          label='Обязательный'
      ></v-checkbox>
      <v-btn class="delete-btn" @click="$emit('delete', id)">
        <v-icon>mdi-delete</v-icon>
      </v-btn>
    </div>
  </v-card>
</template>

<script>

export default {
  data() {
    return {
      types : ["Один выбор", "Множественный выбор", "Ввод вручную"],
    };
  },
  emits: ['delete'],
  methods: {
    updateVariants(e) {
      switch (e) {
        case 'Один выбор':
          this.$emit('update:type', "single_choice");
          this.answers.correct_id = null;
          this.variants.single_choice.fields = [
            {text: 'Вариант 1', id: 0}, {text: 'Вариант 2', id: 1}, {text: 'Вариант 3', id: 2}
          ];
          break;
        case 'Множественный выбор':
          this.$emit('update:type', "multiple_choice");
          this.answers.correct_ids = [];
          this.variants.multiple_choice.fields = [
            {text: 'Вариант 1', id: 0}, {text: 'Вариант 2', id: 1}, {text: 'Вариант 3', id: 2}
          ];
          this.updateMaxChoices()
          break;
        case 'Ввод вручную':
          this.$emit('update:type', "manual_input");
          this.answers.correct_text = '';
          break;
      }
    },
    updateMaxChoices(){
       this.variants.multiple_choice.max = this.variants.multiple_choice.fields.length
    },
    createVariantMultipleChoice(){
      this.variants.multiple_choice.fields.push({id: this.variants.multiple_choice.fields.length, text: ''})
      this.updateMaxChoices()
    },
    deleteVariantMultipleChoice(id){
      this.variants.multiple_choice.fields.splice(id, 1)
      this.updateMaxChoices()
    }
  },
};
</script>

<script setup>
import {defineModel} from "vue";

const id = defineModel('id')
const longText = defineModel('longText')
let shortText = defineModel('shortText')
let type = defineModel('type')
let required = defineModel('required')
let answers = defineModel('answers', {default: {correct_id: null, correct_ids: [], correct_text: ''}})
let variants = defineModel('variants', {default: {single_choice: {fields: []}, multiple_choice: {max_choices: 0, fields: []}}})
let points = defineModel('points')
</script>

<style scoped>
.inline-text-field{
  display: inline-block;
  margin-left: 10px;
}

.variant{
  display: flex;
  align-items: flex-start;
  width: 35%;
}

.item-btn{
  max-width: 15%;
  padding-top: 7px;
}

.required{
  margin-left: 10px;
}

.bot-row{
  display: flex;
  justify-content: space-between;
}

.delete-btn{
  margin-top: 7px;
  margin-right: 15px;
}

.delete-btn-variants{
  margin-left: 10px;
  margin-top: 7px;
}

.create-variant{
  margin-left: 40px;
  margin-bottom: 15px;
  max-width: 5%;
}
</style>
