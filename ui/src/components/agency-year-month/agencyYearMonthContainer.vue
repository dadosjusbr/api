<template>
  <div>
    <h1>{{ this.agencyName }}</h1>
    <h1>{{ this.year }}</h1>
    <h1>{{ this.month }}</h1>
  </div>
</template>

<script>
export default {
  data() {
    return {
      agencyName: this.$route.params.agencyName,
      year: this.$route.params.year,
      month: this.$route.params.month
    };
  }
};
</script>

<style></style>
