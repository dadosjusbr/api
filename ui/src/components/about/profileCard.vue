<template>
  <div class="teamPhotoCard">
    <img class="profilePhoto" src="https://via.placeholder.com/180" />
    <h4>{{ this.name }}</h4>
    <h5>{{ this.role }}</h5>
  </div>
</template>

<script>
export default {
  name: "profileCard",
  props: {
    image: {
      type: String,
      default: ""
    },
    role: {
      type: String,
      default: ""
    },
    name: {
      type: String,
      default: ""
    }
  }
};
</script>

<style scoped>
.profilePhoto {
  border-radius: 50%;
}

.teamPhotoCard {
  text-align: center;
  margin-right: 70px;
  margin-left: 70px;
}

h4 {
  margin-top: 7px;
  font-weight: 600;
}
</style>
