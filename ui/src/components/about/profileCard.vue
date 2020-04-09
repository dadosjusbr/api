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
  width: 125px;
  height: 125px;
}

.teamPhotoCard {
  color: white;
  text-align: center;
  margin-right: 40px;
  margin-left: 40px;
  height: 250px;
  width: 200px;
}

@media only screen and (max-width: 1276px) {
  .teamPhotoCard {
    text-align: center;
    margin-right: 5px;
    margin-left: 5px;
    height: 200px;
    width: 200px;
  }
  .profilePhoto {
    border-radius: 50%;
    width: 100px;
    height: 100px;
  }
  h4 {
    font-size: 20px;
  }
  h5 {
    font-size: 15px;
  }
}

@media only screen and (max-width: 950px) {
  .teamPhotoCard {
    text-align: center;
    margin-right: 5px;
    margin-left: 5px;
    height: 125px;
    width: 200px;
  }
  .profilePhoto {
    border-radius: 50%;
    width: 75px;
    height: 75px;
  }
  h4 {
    font-size: 15px;
  }
  h5 {
    font-size: 10px;
  }
}

h4 {
  font-weight: 600;
}
</style>
