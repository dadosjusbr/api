<template>
  <b-container class="pl-0 pr-0" style="background-color: #2fbb96;" fluid>
  <b-row class=" align-items-center justify-content-center">
    <a :href="this.github" target="_blank">
      <img class="profilePhoto" :src="this.image" />
    </a>
    <br />
    <br />
  </b-row>
  <b-row class="mt-2 align-items-center justify-content-center">
    <b>{{ this.name }}</b>
    <br />
  </b-row>
  <b-row class=" align-items-center justify-content-center">
    {{ this.role }}
  </b-row>
  </b-container>
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
  height: 6em;
  width: 60%;
  font-size: 1.3em;
}

@media only screen and (max-width: 500px) {
  .teamPhotoCard {
    text-align: center;
    height: 350px;
    width: 90%;
    font-size: 1.3em;
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
    font-size: 1.3em;
  }
  .profilePhoto {
    border-radius: 50%;
    width: 140px;
    height: 140px;
  }
}

</style>
