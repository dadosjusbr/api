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
  color: #3e5363;
  text-align: center;
  height: 250px;
  width: 250px;
  font-size: 22px;
}

@media only screen and (max-width: 500px) {
  .teamPhotoCard {
    text-align: center;
    height: 350px;
    width: 90%;
    font-size: 19px;
  }
  .profilePhoto {
    border-radius: 50%;
    width: 150px;
    height: 150px;
  }
}
@media only screen and (max-width: 320px) {
  .teamPhotoCard {
    text-align: center;
    height: 250px;
    width: 95%;
    font-size: 17px;
  }
  .profilePhoto {
    border-radius: 50%;
    width: 140px;
    height: 140px;
  }
}
</style>
