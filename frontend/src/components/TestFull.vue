<template>
  <v-container class="bg-surface-variant">
    <v-row no-gutters>
      <v-col
          v-for="n in 3"
          :key="n"
          cols="12"
          sm="4"
      >
        <v-sheet class="ma-2 pa-2">
          One of three columns
        </v-sheet>
      </v-col>
    </v-row>
  </v-container>
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

</style>