<template>
  <div class="footerContainer">
    <div class="footerAux">
      <img class="logo" src="../../assets/icon_dadosjusbr_white.svg" alt="" />
      <div class="footerTexto">
        Alguma dica? Tem algum feedback geral? Se você tiver uma ideia que você
        gostaria de ver no DadosJusBr, envie-nos um e-mail
        (dadosjusbr@gmail.com) ou visite nosso
        <a
          style="color: white; font-size: 20px;"
          href="https://github.com/dadosjusbr"
          >github</a
        >.
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "pageFooter",
};
</script>

<style scoped>
.footerContainer {
  height: 280px;
  padding: 0px 100px 0px 100px;
  background-color: #3e5363;
  color: white;
  font-size: 20px;
}

.footerAux {
  border-top: 1px solid white;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  padding: 100px 70px 0px 70px;
}

.footerTexto {
  width: 70%;
}

@media only screen and (max-width: 500px) {
  .footerContainer {
    height: 280px;
    padding: 0px 20px 0px 20px;
  }
  .logo {
    display: none;
  }
  .footerTexto {
    width: 100%;
  }
  .footerAux {
    padding: 60px 0px 0px 0px;
    width: 100%;
  }
}
</style>
