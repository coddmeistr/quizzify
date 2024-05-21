<template>
  <div class="container">
  <v-card width="1800">
      <v-card-title class="text-white">{{ title }}</v-card-title>

      <v-card-text>
        <div class="mb-3">
          <v-alert-title>{{shortText}}</v-alert-title>
        </div>

        <v-chip-group class="mb-6">
          <v-chip v-for="(tag, index) in tags" :key="index">
            {{ tag }}
          </v-chip>
        </v-chip-group>

        <div class="text-muted">
          Создатель: {{ creatorName }}
        </div>
      </v-card-text>

      <v-card-actions>
        <v-btn :to="{ path: `/tests/${id}` }" color="primary">Начать тест</v-btn>
      </v-card-actions>
  </v-card>
  </div>
</template>

<script>
export default {
  props: {
    id: Number,
    title: String,
    type: String,
    creatorName: String,
    shortText: String,
    longText: String,
    questions: Object,
    tags: Array,
  },
  methods: {
    action() {
      this.$emit('action', this.id)
    }
  }
}
</script>

<style scoped>
.mb-3 {
  margin-bottom: 16px;
}

.test-img{
  max-height: 100px;
  max-width: 100%;
}

.mb-6 {
  margin-bottom: 32px;
}

.container{
  max-width: 70%;
  min-height: 70%;
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>