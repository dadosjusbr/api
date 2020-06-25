<template>
  <div class="teamPhotoCard">
    <a :href="this.github" target="_blank">
      <img class="profilePhoto" :src="this.image" />
    </a>
    <h4 style="padding-top: 15px">{{ this.name }}</h4>
    <h5>{{ this.role }}</h5>
  </div>
</template>

<script>
export default {
  name: "profileCard",
  props: {
    image: {
      type: String,
      default: "",
    },
    role: {
      type: String,
      default: "",
    },
    name: {
      type: String,
      default: "",
    },
    github: {
      type: String,
      default: "",
    },
  },
};
</script>

<style scoped>
.profilePhoto {
  border-radius: 50%;
  width: 135px;
  height: 135px;
}

.teamPhotoCard {
  color: white;
  text-align: center;
  height: 250px;
  width: 250px;
}

@media only screen and (max-width: 500px) {
  .teamPhotoCard {
    color: white;
    text-align: center;
    height: 200px;
    width: 150px;
  }
  .profilePhoto {
    border-radius: 50%;
    width: 70px;
    height: 70px;
  }
}

h4 {
  font-weight: 600;
}
</style>
