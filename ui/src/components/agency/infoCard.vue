<template>
  <div class="container">
    <div class="md-elevation-24">
    <md-card class="score">
      <md-card-content class="header">
         {{ info.value }}
      </md-card-content>
      <md-card-content class="text">
       {{ info.name.replace("_", " ") }}
      </md-card-content>
    </md-card>
    </div>
  </div>
</template>

<script>
export default {
  name: "infoCard",
  props: {
    info: {
      type: Object,
      default: null
    }
  }
};
</script>

<style scoped>
.container {
  height: 7vw;
  width: 14vw;
}

.header {
  font-size: 1.2vw;
  line-height: 0px;
  padding: 1.0vw;
}


.score {
  padding-top: 25%;
  text-align: center;
  background: #87bba2;
  border: 1px solid black;
  margin-top: 10px;
  font-size: 20px;
  height: 7vw;
  width: 120%;
  color: black;
  font-weight: 500;
}

.text {
  padding-top: 2%;
  padding-left: 5px;
  text-align: left;
  position: relative;
  font-weight: 410;
  font-size: 1vw;
  color: black;
}

@media only screen and  (max-width: 358px)  {

  .container {
    height: 26vw;
    width: 46vw;
    padding-right: 0px; 
    padding-left: 0px;
  }

  .score {
    justify-content: space-evenly;
    padding-top: 0%;
    height: 23vw;
    width: 100%;
  }
  .header {
    font-size: 4.2vw;
    line-height: 0px;
    margin-top: 15%;
    margin-bottom: 0%;
    padding: 16px;
  }

  .text {
    padding-top: 0%;
    font-size: 3.7vw;
  }
}

@media only screen and  (min-width: 359px) and (max-width: 600px)  {
  .container {
    height: 26vw;
    width: 44vw;
    padding-right: 0px; 
    padding-left: 0px;
  }

  .score {
    justify-content: space-evenly;
    padding-top: 0%;
    height: 23vw;
    width: 100%;
  }
  .header {
    font-size: 4.2vw;
    line-height: 0px;
    margin-top: 15%;
    margin-bottom: 0%;
    padding: 16px;
  }

  .text {
    padding-top: 4%;
    font-size: 3.7vw;
  }

}

@media only screen and  (min-width: 601px) and (max-width: 870px)  {
  .container {
    width: 16vw;
  }

  .header {
    padding: 0.7vw;
    font-size: 1.5vw;
  }

  .text {
    font-size: 1.2vw;
    padding-top: 0%;
  }

}

@media only screen and  (min-width: 771px) and (max-width: 1025px)  {
  .score {
    padding-top: 15%;
    width: 130%;  
  }

  .text {
    font-size: 1.2vw;
    padding-top: 5%;
  }
  .header {
    font-size: 1.5vw;
  }

}



</style>
