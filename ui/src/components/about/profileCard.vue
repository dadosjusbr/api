<template>
  <div class="teamPhotoCard">
    <a :href="this.github" target="_blank">
      <img class="profilePhoto" :src="this.image" />
    </a>
    <br />
    <br />
    <b>{{ this.name }}</b>
    <br />
    {{ this.role }}
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
  width: 190px;
  height: 190px;
}

.teamPhotoCard {
  color: white;
  text-align: center;
  height: 250px;
  width: 250px;
  font-size: 22px;
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
</style>
