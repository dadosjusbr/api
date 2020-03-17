<template>
  <div class="teamPhotoCard">
    <a :href="this.github" target="_blank">
      <img class="profilePhoto" :src="this.image" />
    </a>
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
    },
    github: {
      type: String,
      default: ""
    }
  }
};
</script>

<style scoped>
.profilePhoto {
  border-radius: 50%;
  width: 150px;
  height: 150px;
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
