<template>
  <div class="contactContainer">
    <h1>CONTATO</h1>
    <p>
      Alguma dica? Tem algum feedback geral? Se você tiver uma ideia que você
      gostaria de ver no DadosJusBr, envie-nos um e-mail (dadosjusbr@gmail.com)
      ou visite nosso github.
    </p>
    <div class="logoContainer">
      <a href="mailto:dadosjusbr@gmail.com">
        <img src="../../assets/logo_email.png" />
      </a>
      <a href="https://github.com/dadosjusbr">
        <img class="img2" src="../../assets/logo_github.png" />
      </a>
    </div>
  </div>
</template>

<script>
export default {
  name: "contact"
};
</script>

<style scoped>
.contactContainer {
  border: 1px solid #01805A;
  margin-left: 200px;
  margin-right: 200px;
  margin-top: 10px;
  background-color: #01805A;
  height: 700px;
  margin-bottom: 10px;
}

h1 {
  color: white;
  text-align: center;
  padding-top: 50px;
  font-size: 50px;
}

p {
  padding-left: 60px;
  padding-right: 60px;
  padding-top: 25px;
  text-align: justify;
  color: white;
  font-size: 25px;
  margin-top: 50px;
}

img {
  width: 70px;
  height: 70px;
}

.img2 {
  margin-left: 400px;
}

.logoContainer {
  margin-top: 180px;
  justify-content: center;
  display: flex;
  align-self: center;
  align-content: center;
}

@media only screen and (max-width: 1150px) {
  .img2 {
    margin-left: 200px;
  }
  .contactContainer {
    margin-left: 100px;
    margin-right: 100px;
  }
}

@media only screen and (max-width: 895px) {
  .img2 {
    margin-left: 100px;
  }
  .contactContainer {
    margin-left: 20px;
    margin-right: 20px;
  }
}

@media only screen and (max-width: 544px) {
  .img2 {
    margin-left: 100px;
  }
  .contactContainer {
    margin-left: 10px;
    margin-right: 10px;
    height: 600px;
  }
  .logoContainer {
    margin-top: 100px;
  }
  p {
    font-size: 20px;
    padding-left: 30px;
    padding-right: 30px;
  }
  h1 {
    font-size: 40px;
  }
  img {
    width: 50px;
    height: 50px;
  }
}
</style>
