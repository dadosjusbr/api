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
  border: 1px solid #6a757a;
  margin-left: 200px;
  margin-right: 200px;
  margin-top: 10px;
  background-color: #6a757a;
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
</style>
