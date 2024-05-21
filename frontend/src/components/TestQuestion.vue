<template>
  <v-card class="question">
    <v-card-title>{{ longText }}</v-card-title>

    <v-card-text class="question__text">{{ shortText }}</v-card-text>

    <template v-if="type === 'single_choice' || type === 'multiple_choice'">
      <v-list>
        <v-radio-group v-if="type === 'single_choice'" v-model="selectedVariant">
        <v-list-item v-for="(variant) in variantsFields" :key="variant.id">
          <v-radio
              @change="onVariantChange(selectedVariant)"
              :value="variant.id"
              :label="variant.text"
          />
        </v-list-item>
        </v-radio-group>
        <div v-if="type === 'multiple_choice'">
          <v-list-item v-for="(variant) in variantsFields" :key="variant.id">
            <v-checkbox
                @change="onVariantsChange(selectedVariants)"
                v-model="selectedVariants"
                :value="variant.id"
                :label="variant.text"
            />
          </v-list-item>
        </div>
      </v-list>
    </template>

    <template v-else-if="type === 'manual_input'">
      <v-text-field
          @change="onWritedTextChange(writedText)"
          v-model="writedText"
          label="Введите ответ"
      />
    </template>

    <v-card-subtitle v-if="required">
      <v-icon color="error">mdi-star</v-icon> Обязательный
    </v-card-subtitle>

    <v-card-subtitle>Баллы: {{ points }}</v-card-subtitle>
  </v-card>
</template>

<script>
import {defineModel} from 'vue'

const selectedVariant = defineModel({default: null})
const selectedVariants = defineModel({default: [], type: Array})
const writedText = defineModel({default: ''})

export default {
  name: "TestQuestion",

  data(){
    return {
      selectedVariant,
      selectedVariants,
      writedText
    }
  },
  props: {
    id: Number,
    type: String,
    longText: String,
    shortText: String,
    required: Boolean,
    variants: Object,
    points: Number,
  },
  components: {},
  methods:{
    resetAnswers(){
      if (this.type === "single_choice") {
        this.selectedVariant = null;
      } else if (this.type === 'manual_input') {
        this.writedText = "";
      } else if (this.type === 'multiple_choice') {
        this.selectedVariants = [];
      }
    },
    onVariantChange(value) {
      this.$emit('update:selectedVariant', value);
    },
    onVariantsChange(value) {
      this.$emit('update:selectedVariants', value);
    },
    onWritedTextChange(value) {
      this.$emit('update:writedText', value);
    },
  },
  mounted() {
    this.selectedVariants = []
  },
  computed: {
    variantsFields() {
      if (this.variants && this.variants[this.type] && this.variants[this.type].fields) {
        return this.variants[this.type].fields;
      } else {
        return [];
      }
    }
  },
};
</script>

<style scoped>
.question {
  margin-bottom: 20px;
}
</style>