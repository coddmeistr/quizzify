<template>
  <div class="card">
    <v-card variant="outlined">
      <v-card-title>{{title}}</v-card-title>
      <v-card-subtitle v-if="questionsCount > 0">{{questionsCount}} вопросов</v-card-subtitle>
      <v-card-text>{{shortText}}</v-card-text>
      <v-card-actions>
        <v-btn color="primary" @click="action">Подробнее</v-btn>
      </v-card-actions>
    </v-card>
  </div>
</template>

<script>
export default {
  props: {
    id: Number,
    title: String,
    shortText: String,
    questionsCount: Number
  },
  methods: {
    action() {
      this.$emit('action', this.id)
    }
  }
}
</script>

<style scoped>

.card{
  max-width: 300px;
}

</style>